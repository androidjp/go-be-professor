package workerid

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/assert"
)

func TestInitClusterWorkerID(t *testing.T) {

	// 启动一个模拟的 Redis 服务器
	mrs, err := miniredis.Run()
	if err != nil {
		t.Fatalf("无法启动模拟 Redis 服务器: %v", err)
	}
	defer mrs.Close()

	redisCli := redis.NewClient(&redis.Options{
		Addr: mrs.Addr(),
	})

	redisCli.Del(context.Background(), "worker_id:test")
	id, err := InitClusterWorkerID(redisCli, "test")
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, int64(1), id)
	id, err = InitClusterWorkerID(redisCli, "test")
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, int64(1), id)
}
