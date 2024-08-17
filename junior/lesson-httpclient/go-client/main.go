package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	DEFAULT = iota
	MAX_CONN_5
)

func main() {

	// 1. 初始化原生httpClient
	httpCli, err := initHTTPClient(DEFAULT)
	if err != nil {
		fmt.Println("initHTTPClient err:", err)
		return
	}

	// 3. 100个协程并发发起http请求
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			startTS := time.Now()
			defer func() {
				wg.Done()
				endTS := time.Now()
				fmt.Println(i, "cost:", endTS.Sub(startTS))
			}()
			// 40秒超时，表示：整个时间窗口 [发起http请求（即使此时没有可用长连接） , 响应返回] 大于40秒就会掐断
			// ctx, cancelFunc := context.WithTimeout(context.Background(), 40*time.Second)
			// defer cancelFunc()
			ctx := context.Background()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8080/api/hello/20", nil)
			if err != nil {
				fmt.Println(i, "http.NewRequestWithContext err:", err)
				return
			}
			resp, respErr := httpCli.Do(req)
			if respErr != nil {
				fmt.Println(i, "http.Client.Do err:", respErr)
				return
			}
			defer resp.Body.Close()
			fmt.Println(i, "resp.StatusCode:", resp.StatusCode)
			// jsonBS, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	fmt.Println("io.ReadAll err:", err)
			// 	return
			// }
			// fmt.Println(i, "resp.Body:", string(jsonBS))
		}(i)
	}

	wg.Wait()
	fmt.Println("Done...")
}
func initHTTPClient(cliType int) (*http.Client, error) {
	switch cliType {
	case DEFAULT:
		cli := http.DefaultClient
		return cli, nil
	case MAX_CONN_5:
		cli := &http.Client{}
		cli.Transport = &http.Transport{
			MaxIdleConns:          5,
			MaxIdleConnsPerHost:   5,
			MaxConnsPerHost:       5,
			TLSHandshakeTimeout:   10 * time.Second, // tls握手超时时间
			IdleConnTimeout:       5 * time.Second,  // 如果5秒内没有请求，就会关闭连接
			ResponseHeaderTimeout: 25 * time.Second, // 时间窗口 [发起http请求（连接已发起开始算） , 开始响应返回] 大于30秒就会掐断，默认就是无限大
			DialContext: (&net.Dialer{
				Timeout: 2 * time.Second, // 连接超时时间为3秒
			}).DialContext,
		}
		return cli, nil
	default:
		return nil, fmt.Errorf("unknown http client type")
	}
}
