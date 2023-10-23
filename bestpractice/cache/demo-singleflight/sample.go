package demo_singleflight

import (
	"errors"
	"log"
	"sync"
)

var errorNotExist = errors.New("not exist")

func Exec1() {
	var wg sync.WaitGroup
	wg.Add(10)

	// 10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			// 拿数据
			data, err := getData("热点key")

			if err != nil {
				log.Println(err)
				return
			}

			log.Println(data)
		}()
	}

	wg.Wait()
}

// 拿数据(无singleflight)
func getData(key string) (string, error) {
	// 1. 优先缓存
	data, err := getDataFromCache(key)
	if err == errorNotExist {
		//2. 模拟从db中获取数据
		data, err = getDataFromDB(key)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	return data, nil
}

//模拟从cache中获取值，cache中无该值
func getDataFromCache(key string) (string, error) {
	return "", errorNotExist
}

//模拟从数据库中获取值
func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database", key)
	return "data", nil
}
