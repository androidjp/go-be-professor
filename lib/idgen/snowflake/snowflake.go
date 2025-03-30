package snowflake

import (
	"context"
	"mylib/idgen/snowflake/workerid"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/zheng-ji/goSnowFlake"
)

type IdGenerator struct {
	worker *goSnowFlake.IdWorker
}

func NewIdGenerator(workerId int64) (*IdGenerator, error) {
	w, err := goSnowFlake.NewIdWorker(workerId)
	if err != nil {
		return nil, err
	}
	return &IdGenerator{w}, nil
}

func NewClusterIdGenerator(redisCli redis.UniversalClient, appName string) (*IdGenerator, error) {
	workerId, err := workerid.InitClusterWorkerID(redisCli, appName)
	if err != nil {
		return nil, err
	}
	w, err := goSnowFlake.NewIdWorker(workerId)
	if err != nil {
		return nil, err
	}
	return &IdGenerator{w}, nil
}

func (s *IdGenerator) NextId(ctx context.Context) (int64, error) {
	return s.worker.NextId()
}

func ParseId(id int64) (t time.Time, ts int64, workerId int64, seq int64) {
	return goSnowFlake.ParseId(id)
}
