package idgen

import (
	"demo-redis-snowflake-distribute-id/workerid"
	"fmt"
	"github.com/zheng-ji/goSnowFlake"
	"sync"
	"time"
)

type IDGenerator struct {
	GoSnowFlakeIDGenerator *goSnowFlake.IdWorker
}

var generator *IDGenerator
var initOnce sync.Once

func Get() *IDGenerator {
	initOnce.Do(func() {
		// 从redis获取唯一的workerID
		workerID, err := workerid.GetWorkerID()
		if err != nil {
			panic(any(err))
		}

		// 初始化ID生成器
		idGenerator, err := goSnowFlake.NewIdWorker(workerID)
		if err != nil {
			panic(any(err))
		}
		generator = &IDGenerator{
			GoSnowFlakeIDGenerator: idGenerator,
		}
	})
	return generator
}

func (g *IDGenerator) GenSnowID() int64 {
	snowID, err := g.GoSnowFlakeIDGenerator.NextId()
	if err != nil {
		panic(any(err))
	}
	return snowID
}

const epoch = int64(1474802888000) // 自定义起始时间戳（毫秒） 对应UTC时间：2016年9月25日 19:28:08

// SnowIdToTime 雪花算法id转时间戳，无法精确定位，定位到日即可
func (g *IDGenerator) SnowIdToTime(metaId int64) int64 {
	snowflakeID := metaId                    // Snowflake算法生成的ID
	timestamp := (snowflakeID >> 22) + epoch // 右移22位，获取时间戳
	// 将时间戳转换为时间格式
	t := time.Unix(0, timestamp*int64(time.Millisecond))
	fmt.Println("Generated Timestamp:", t)
	return t.Unix()
}
