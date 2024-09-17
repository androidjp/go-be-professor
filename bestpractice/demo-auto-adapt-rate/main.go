package main

import (
	"demo-auto-adapt-rate/autorate"
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func login1(mgr *autorate.AutoRateMgr, dur time.Duration) {
	defer mgr.Statistic("login_1")()

	// 打印
	fmt.Printf("login1 执行，耗时：%v\n", dur)
	// 模拟目标函数执行一些耗时操作
	time.Sleep(dur)
}

func main() {
	var (
		err error
	)

	// 设置配置文件名称（不包含文件类型后缀）
	viper.SetConfigName("app")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 设置配置文件所在路径
	viper.AddConfigPath("./conf")
	// 启动一个监控项,当配置文件发生改动后自动加载新的配置项目，需要注意的是，在使用时需提前设置配置文件路径信息
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgMap := make(map[autorate.AutoRateKey]autorate.AutoRateConfig)
	err = viper.UnmarshalKey("autorate", &cfgMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr := autorate.NewAutoRateMgr(cfgMap)
	defer mgr.Stop()
	mgr.Start()

	// 模拟并发请求高耗时情况
	for i := 0; i < 10; i++ {
		wg := sync.WaitGroup{}
		cnt := int(mgr.GetCurNumber("login_1"))
		wg.Add(cnt)

		fmt.Printf("底%d轮的测试验证，请求量为：%d\n", i+1, cnt)

		for j := 0; j < cnt; j++ {
			go func() {
				defer wg.Done()
				login1(mgr, 1*time.Second)
			}()
		}
		time.Sleep(2 * time.Second)
		wg.Wait()
	}

	// 模拟并发请求低耗时情况
	for i := 0; i < 10; i++ {
		wg := sync.WaitGroup{}
		cnt := int(mgr.GetCurNumber("login_1"))
		wg.Add(cnt)

		fmt.Printf("底%d轮的测试验证，请求量为：%d\n", i+1, cnt)

		for j := 0; j < cnt; j++ {
			go func() {
				defer wg.Done()
				login1(mgr, 100*time.Millisecond)
			}()
		}
		time.Sleep(2 * time.Second)
		wg.Wait()
	}

	// 模拟并发请求 正常情况
	for i := 0; i < 10; i++ {
		wg := sync.WaitGroup{}
		cnt := int(mgr.GetCurNumber("login_1"))
		wg.Add(cnt)

		fmt.Printf("底%d轮的测试验证，请求量为：%d\n", i+1, cnt)

		for j := 0; j < cnt; j++ {
			go func() {
				defer wg.Done()
				login1(mgr, 500*time.Millisecond)
			}()
		}
		time.Sleep(2 * time.Second)
		wg.Wait()
	}

}
