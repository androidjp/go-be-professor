package demo_singleflight

import (
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
	"time"
)

// 带有singleflight 和 ttl key过期时间，过期后，这个key就不会再缓存了，需要重新往DB查。
// 比如：我设置了ttl为2秒，那么，2秒内，只会往DB中查一次数据。

func Exec3() {
	var wg sync.WaitGroup
	wg.Add(10000)

	// 10个并发
	for i := 0; i < 10000; i++ {
		go func() {
			defer wg.Done()

			// 拿数据
			data, err := getDataWithTTL("热点key", 2*time.Second)

			if err != nil {
				log.Println(err)
				return
			}

			log.Println(data)
		}()
	}

	wg.Wait()
}

var gsf2 singleflight.Group

//获取数据(加了singleflight)
func getDataWithTTL(key string, ttl time.Duration) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		//模拟从db中获取数据
		v, err, _ := gsf2.Do(key, func() (interface{}, error) {
			return getDataFromDB(key)
			//set cache
		})
		if err != nil {
			log.Println(err)
			return "", err
		}

		// 设置一个异步任务，2秒后忘记这个key
		go func(key string) {
			time.Sleep(ttl)
			gsf2.Forget(key)
		}(key)

		//TOOD: set cache
		data = v.(string)
	} else if err != nil {
		return "", err
	}
	return data, nil
}
