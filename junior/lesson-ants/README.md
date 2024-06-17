# 协程池管理-ants

# 调整参数
ants 提供了多种参数来优化池的行为。例如，可以调整池的大小、设置非阻塞模式、以及设置空闲的超时时间等。

## 设置池大小
我们可以创建一个不同大小的池，通过 NewPoolWithFunc 来指定池大小：

```go
p, _ := ants.NewPoolWithFunc(20, func(i interface{}) {
runTask(i)
})
```

## 设置非阻塞模式
默认情况下，如果池已满，新任务会阻塞直到有空闲 goroutine 可用。我们可以通过设置非阻塞模式来改变这一行为：

```go
p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
runTask(i)
}, ants.WithNonblocking(true))

```

在非阻塞模式下，如果池已满，新任务会立即返回错误，而不是等待。

## 设置任务超时
我们还可以设置池中任务的超时时间，以避免某些任务耗时过长：
```go
p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
runTask(i)
}, ants.WithExpireDuration(5 * time.Second))

```
在这个示例中，每个任务最多只能运行 5 秒钟，超过这个时间，任务将被自动终止。

## 示例：网络爬虫
假设我们正在编写一个网络爬虫，需要并发访问多个网页。可以使用 ants 来管理并发连接：

```go
package main

import (
"fmt"
"net/http"
"io/ioutil"
"github.com/panjf2000/ants/v2"
)

func main() {
urls := []string{
"http://example.com",
"http://example.org",
"http://example.net",
}

    fetchURL := func(url string) {
        resp, err := http.Get(url)
        if err != nil {
            fmt.Printf("Failed to fetch %s: %s\n", url, err)
            return
        }
        defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Printf("Fetched from %s: %d bytes\n", url, len(body))
    }

    p, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
        fetchURL(i.(string))
    })
    defer p.Release()

    for _, url := range urls {
        _ = p.Invoke(url)
    }

    p.Wait()
    fmt.Println("All URLs fetched")
}
```

在这个示例中，我们创建了一个网络爬虫，使用 ants 来并发地访问多个网页。爬虫只使用 5 个 goroutines，并发访问多个 URL。

# 扩展思路
除了基本的任务提交和参数调整，还有许多高级用法。下面列举几个：

## 动态调整池大小
有时候，我们需要根据运行时环境动态调整池的大小。ants 支持在运行时调整池大小：

```go
pool.Tune(20)
```

## 性能监控
ants 内置了一些性能监控功能，可以用来观察池的运行状态：
```go
fmt.Printf("Running goroutines: %d\n", pool.Running())
fmt.Printf("Free goroutines: %d\n", pool.Free())

```

## 异常处理
可以为每个任务设置异常处理机制，确保 goroutine 不会因为 panic 而崩溃：

```go
p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic", r)
        }
    }()
    runTask(i)
})

```

# 总结
使用 ants 可以显著提高 Go 应用程序的并发能力，使得 goroutines 的管理更加高效和可靠。通过调整参数，我们可以更灵活地控制池的行为，以满足不同的应用场景需求。无论是简单的并发任务，还是复杂的并发控制，ants 都提供了强大的支持。

希望通过本文，你能够对 ants 有一个全面的理解，并能在实际项目中灵活