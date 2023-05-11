# 仿造开源的限流器（gin作为示例）

https://github.com/juju/ratelimit/blob/master/reader.go

默认的令牌桶，fillInterval 指每过多长时间向桶里放一个令牌，capacity 是桶的容量，超过桶容量的部分会被直接丢弃。桶初始是满的

```go
func NewBucket(fillInterval time.Duration, capacity int64) *Bucket
```


和普通的 NewBucket() 的区别是，每次向桶中放令牌时，是放 quantum 个令牌，而不是一个令牌。

```go
func NewBucketWithQuantum(fillInterval time.Duration, capacity, quantum int64) *Bucket
```


按照提供的比例，每秒钟填充令牌数。例如 capacity 是100，而 rate 是 0.1，那么每秒会填充10个令牌。
```go
func NewBucketWithRate(rate float64, capacity int64) *Bucket
```
