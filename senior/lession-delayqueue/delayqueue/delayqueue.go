package delayqueue

import (
	"context"
	"lession-delayqueue/generics/container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type entry[T any] struct {
	value      T
	expiration time.Time // 到期时间
}

/**
主要的结构可以看到就是一个heap，entry是每个元素在堆中的表示，value是具体的元素值，expiration是为了堆中元素根据到期时间排序。
mutex是一个互斥锁，主要是保证操作并发安全。
sleeping则是表示Take()，也就是阻塞获取元素操作，是否在等待队列不为空或者有更小到期时间元素加入。这样Push()，也就是添加元素操作，才知道是否去唤醒Take()。（重点）
wakeup是一个通道，通过它实现添加元素的时候唤醒等待的Take()。（重点）
简单来说就是sleeping是表示是否要唤醒Take()，唤醒操作则是通过wakeup通道实现。
*/

/**
原理：
Take()的时候如果队列已经没有元素，或者没有元素到期，那么协程就需要挂起等待。而被唤醒的条件是元素到期、队列不为空或者有更小到期时间元素加入。
其中元素到期协程在Take()时发现堆顶元素还没到期，因此这个条件可以自己构造并等待。
但是条件队列不为空和有更小到期时间元素加入则需要另外一个协程在Push()时才能满足，
因此必须通过一个中间结构来进行协程间通信，一般Golang里面会使用Channel来实现。
而Take()是否在等待则是通过sleeping来表示。
*/

// 延迟队列
// 参考https://github.com/RussellLuo/timingwheel/blob/master/delayqueue/delayqueue.go
type DelayQueue[T any] struct {
	h *heap.Heap[*entry[T]]
	// // 保证并发安全
	mutex sync.Mutex
	// 表示Take()是否正在等待队列不为空或更早到期的元素
	// 0表示Take()没在等待，1表示Take()在等待
	sleeping int32
	// 唤醒通道
	wakeup chan struct{}
}

// 创建延迟队列
func New[T any]() *DelayQueue[T] {
	return &DelayQueue[T]{
		h: heap.New(nil, func(e1, e2 *entry[T]) bool {
			return e1.expiration.Before(e2.expiration)
		}),
		wakeup: make(chan struct{}),
	}
}

// 添加延迟元素到队列
/**
一开始加了一个互斥锁，避免并发冲突，然后把元素加到堆里。
因为我们Take()操作在不满足条件时会去设置sleeping为1表示正在等待Push()来唤醒它，并在wakeup通道阻塞读取。
因为sleeping可能被Take()和Push()同时操作，因此使用CAS()来进行设置，也就是如果sleeping本来是1，我们就设置为0，然后往wakeup写入一个元素，表示如果有Take()在等待，则唤醒它。
*/
func (q *DelayQueue[T]) Push(value T, delay time.Duration) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	entry := &entry[T]{
		value:      value,
		expiration: time.Now().Add(delay),
	}
	q.h.Push(entry)
	// 唤醒等待的Take()
	// 这里表示新添加的元素到期时间是最早的，或者原来队列为空
	// 因此必须唤醒等待的Take()，因为可以拿到更早到期的元素
	if q.h.Peek() == entry {
		// 把sleeping从1修改成0，也就是唤醒等待的Take()
		if atomic.CompareAndSwapInt32(&q.sleeping, 1, 0) {
			q.wakeup <- struct{}{}
		}
	}
}

// 等待直到有元素到期
// 或者ctx被关闭
/**
这里先判断堆是否有元素，如果有获取堆顶元素，然后判断是否已经到期，如果到期则直接出堆并返回。

否则等待直到超时或者元素到期或者有新的元素到达。

在这里需要设置sleeping为1表示Take()正在等待。
*/
func (q *DelayQueue[T]) Take(ctx context.Context) (T, bool) {
	for {
		var timer *time.Timer
		q.mutex.Lock()
		// 有元素
		if !q.h.Empty() {
			// 获取元素
			entry := q.h.Peek()
			now := time.Now()
			if now.After(entry.expiration) {
				q.h.Pop()
				q.mutex.Unlock()
				return entry.value, true
			}
			// 到期时间，使用time.NewTimer()才能够调用Stop()，从而释放定时器
			timer = time.NewTimer(entry.expiration.Sub(now))
		}
		// 走到这里表示需要等待了，设置为1告诉Push()在有新元素时要通知
		atomic.StoreInt32(&q.sleeping, 1)
		q.mutex.Unlock()

		// 不为空，需要同时等待元素到期，并且除非timer到期，否则都需要关闭timer避免泄露
		if timer != nil {
			select {
			case <-q.wakeup: // 新的更快到期元素
				timer.Stop()
			case <-timer.C: // 首元素到期
				// 设置为0，如果原来也为0表示有Push()正在q.wakeup被阻塞
				if atomic.SwapInt32(&q.sleeping, 0) == 0 {
					// 避免Push()的协程被阻塞
					<-q.wakeup
				}
			case <-ctx.Done(): // 被关闭
				timer.Stop()
				var t T
				return t, false
			}
		} else {
			select {
			case <-q.wakeup: // 新的更快到期元素
			case <-ctx.Done(): // 被关闭
				var t T
				return t, false
			}
		}
	}
}

// 返回一个通道，输出到期元素
// size是通道缓存大小
// channel方式阻塞读取
/**
Golang里面可以使用Channel进行流式消费，因此简单包装一个Channel形式的阻塞读取接口，给通道一点缓冲区大小可以带来更好的性能。
使用方式：
for entry := range q.Channel(context.Background(), 10) {
    // do something
}
*/
func (q *DelayQueue[T]) Channel(ctx context.Context, size int) <-chan T {
	out := make(chan T, size)
	go func() {
		for {
			entry, ok := q.Take(ctx)
			if !ok {
				close(out)
				return
			}
			out <- entry
		}
	}()
	return out
}

// 获取队头元素
func (q *DelayQueue[T]) Peek() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.h.Empty() {
		var t T
		return t, false
	}
	return q.h.Peek().value, true
}

// 获取到期元素
func (q *DelayQueue[T]) Pop() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	// 没元素
	if q.h.Empty() {
		var t T
		return t, false
	}
	entry := q.h.Peek()
	// 还没元素到期
	if time.Now().Before(entry.expiration) {
		var t T
		return t, false
	}
	// 移除元素
	q.h.Pop()
	return entry.value, true
}

// 是否队列为空
func (q *DelayQueue[T]) Empty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.h.Empty()
}
