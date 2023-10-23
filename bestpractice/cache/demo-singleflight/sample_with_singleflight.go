package demo_singleflight

import (
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
)

func Exec2() {
	var wg sync.WaitGroup
	wg.Add(10)

	// 10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			// 拿数据
			data, err := getData2("热点key")

			if err != nil {
				log.Println(err)
				return
			}

			log.Println(data)
		}()
	}

	wg.Wait()
}

var gsf singleflight.Group

//获取数据(加了singleflight)
func getData2(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		//模拟从db中获取数据
		v, err, _ := gsf.Do(key, func() (interface{}, error) {
			return getDataFromDB(key)
			//set cache
		})
		if err != nil {
			log.Println(err)
			return "", err
		}

		//TOOD: set cache
		data = v.(string)
	} else if err != nil {
		return "", err
	}
	return data, nil
}
