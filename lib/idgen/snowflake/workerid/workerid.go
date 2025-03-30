package workerid

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/zheng-ji/goSnowFlake"
)

var (
	workerID int64
	once     sync.Once
)

var (
	MaxWorkerID = int64(-1 ^ -1<<goSnowFlake.CWorkerIdBits)
)

func InitClusterWorkerID(redisCli redis.UniversalClient, appName string) (id int64, err error) {
	once.Do(func() {
		key := "worker_id:" + appName
		const script = `local id=redis.call('incr',KEYS[1]);if(id<=tonumber(ARGV[1])) then return id;end;redis.call('set',KEYS[1],1);return 1`
		id, err = redisCli.Eval(context.Background(), script, []string{key}, MaxWorkerID).Int64()
		workerID = id
	})
	if err != nil {
		return 0, err
	}
	return workerID, nil
}

func SetWorkerID(id int64) {
	workerID = id
}

func GetWorkerID() int64 {
	return workerID
}
