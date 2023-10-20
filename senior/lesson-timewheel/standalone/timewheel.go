package standalone

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type TimeWheel struct {
	// 单例工具，保证时间轮停止操作只能执行一次
	sync.Once
	// 时间轮运行时间间隔
	interval time.Duration
	// 类比：驱动指针运转的齿轮
	// 时间轮定时器
	ticker *time.Ticker
	// 停止时间轮的 channel
	stopc chan struct{}
	// 新增定时任务的入口 channel
	addTaskCh chan *taskElement
	// 删除定时任务的入口 channel
	removeTaskCh chan string
	// 类似时钟的表盘
	// 通过 list 组成的环状数组. 通过遍历环状数组的方式实现时间轮
	// 定时任务数量较大，每个 slot 槽内可能存在多个定时任务，因此通过 list 进行组装
	// 所谓环状数组指的是逻辑意义上的. 在实际的实现过程中，会通过一个定长数组结合循环遍历的方式，来实现这个逻辑意义上的“环状”性质
	slots []*list.List
	// 类似时钟的指针
	// 当前遍历到的环状数组的索引
	curSlot int
	// 定时任务 key 到任务节点的映射，便于在 list 中删除任务节点
	keyToETask map[string]*list.Element
}

func (t *TimeWheel) run() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("catch error: %v\n", err)
			buf := make([]byte, 1<<10)
			size := runtime.Stack(buf, true)
			fmt.Printf("catch stacktrace: %s\n", buf[:size])
		}
	}()

	// 通过 for 循环结合 select 多路复用的方式运行，属于 golang 中非常常见的异步编程风格
	for {
		select {
		// 停止时间轮
		case <-t.stopc:
			return
		// 接收到定时信号
		case <-t.ticker.C:
			// 批量执行定时任务
			t.tick()
		// 接收创建定时任务的信号
		case task := <-t.addTaskCh:
			t.addTask(task)
		// 接收到删除定时任务的信号
		case removeKey := <-t.removeTaskCh:
			t.removeTask(removeKey)
		}
	}
}

// Stop 手动停止时间轮，回收对应的携程和ticker资源
func (t *TimeWheel) Stop() {
	t.Do(func() {
		// 定制定时器 ticker 关闭
		t.ticker.Stop()
		// 关闭定时器运行的 stopc
		close(t.stopc)
	})
}

// AddTask 添加定时任务到时间轮
func (t *TimeWheel) AddTask(key string, task func(), executeAt time.Time) {
	// 根据执行时间推算得到定时任务从属的 slot 位置，以及需要延迟的轮次
	pos, cycle := t.getPosAndCycle(executeAt)
	// 将定时任务通过 channel 进行投递
	t.addTaskCh <- &taskElement{
		pos:   pos,
		cycle: cycle,
		task:  task,
		key:   key,
	}
}

func (t *TimeWheel) getPosAndCycle(executeAt time.Time) (int, int) {
	delay := int(time.Until(executeAt))
	// 定时任务的延迟轮次 = 总延迟 dur / (轮盘格子数*每个格子的时间间隔)
	cycle := delay / (len(t.slots) * int(t.interval))
	// 定时任务从属的环状数组 index，定位到对应的索引位置，准备入链表
	pos := (t.curSlot + delay/int(t.interval)) % len(t.slots)
	return pos, cycle
}

// 常驻 协程 接收到创建定时任务后的处理逻辑：通过pos找到链表，尝试插入链表的适当位置
func (t *TimeWheel) addTask(task *taskElement) {
	// 1. 找到轮盘位置的链表
	list := t.slots[task.pos]
	// 2. 如果这个任务key已经存在任务了，就需要先删除key
	if _, ok := t.keyToETask[task.key]; ok {
		t.removeTask(task.key)
	}
	// 3. 将定时任务追加到list尾部
	eTask := list.PushBack(task)
	// 4. 放入任务映射map中
	t.keyToETask[task.key] = eTask
}

func (t *TimeWheel) removeTask(key string) {
	eTask, ok := t.keyToETask[key]
	if !ok {
		return
	}
	// fixme：注意并发读写异常
	delete(t.keyToETask, key)
	task, _ := eTask.Value.(*taskElement)
	_ = t.slots[task.pos].Remove(eTask)
}

// RemoveTask 删除定时任务
func (t *TimeWheel) RemoveTask(key string) {
	t.removeTaskCh <- key
}

// 内部的定时任务执行逻辑
func (t *TimeWheel) tick() {
	// 1. 当前链表
	list := t.slots[t.curSlot]
	// 2. 方法返回前，让curSlot指针移动到下一个位置
	defer t.circularIncr()
	// 3. 批量处理满足执行条件的定时任务
	t.execute(list)
}

func (t *TimeWheel) circularIncr() {
	t.curSlot = (t.curSlot + 1) % len(t.slots)
}

func (t *TimeWheel) execute(l *list.List) {
	// 遍历 list
	for e := l.Front(); e != nil; {
		// 拿到具体任务信息
		taskInfo, _ := e.Value.(*taskElement)
		// 如果任务cycle还没到，就对他cycle计数器扣减，本轮不做任务执行
		if taskInfo.cycle > 0 {
			taskInfo.cycle--
			e = e.Next()
			continue
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("catch error: %v\n", err)
					buf := make([]byte, 1<<10)
					size := runtime.Stack(buf, true)
					fmt.Printf("catch stacktrace: %s", buf[:size])
				}
			}()
			// 执行
			taskInfo.task()
		}()

		// 任务已经交给 协程处理了，需要把对应任务节点从list中删除
		next := e.Next()
		l.Remove(e)
		// 把任务key从映射map中删除
		delete(t.keyToETask, taskInfo.key)
		e = next
	}
}

// 封装了一笔定时任务的明细信息
type taskElement struct {
	// 内聚了定时任务执行逻辑的闭包函数
	task func()
	// 定时任务挂载在环状数组中的索引位置
	pos int
	// 时间轮中一个 slot 可能需要挂载多笔定时任务，因此针对每个 slot，需要采用 golang 标准库 container/list 中实现的双向链表进行定时任务数据的存储.
	// 定时任务的延迟轮次. 指的是 curSlot 指针还要扫描过环状数组多少轮，才满足执行该任务的条件
	cycle int
	// 定时任务的唯一标识键
	key string
}

// NewTimeWheel 创建单机版时间轮 slotNum——时间轮环状数组长度  interval——扫描时间间隔
func NewTimeWheel(slotNum int, interval time.Duration) *TimeWheel {
	// 环状数组长度默认为 10
	if slotNum <= 0 {
		slotNum = 10
	}
	// 扫描时间间隔默认为 1 秒
	if interval <= 0 {
		interval = time.Second
	}
	// 初始化时间轮实例
	t := TimeWheel{
		interval:     interval,
		ticker:       time.NewTicker(interval),
		stopc:        make(chan struct{}),
		keyToETask:   make(map[string]*list.Element),
		slots:        make([]*list.List, 0, slotNum),
		addTaskCh:    make(chan *taskElement),
		removeTaskCh: make(chan string),
	}
	for i := 0; i < slotNum; i++ {
		t.slots = append(t.slots, list.New())
	}

	// 异步启动时间轮常驻 goroutine
	go t.run()
	return &t
}
