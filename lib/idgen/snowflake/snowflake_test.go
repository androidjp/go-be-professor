package snowflake

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestIdGenerator_NextId(t *testing.T) {
	sf, err := NewIdGenerator(1)
	if err != nil {
		t.Fatal(err)
		return
	}
	lastId := int64(0)
	for i := 0; i < 100; i++ {
		newId, err := sf.NextId(context.Background())
		if err != nil {
			t.Fatal(err)
			return
		}
		//t.Log(newId)
		if newId < lastId {
			t.Fatalf("newId %d should be greater than %d", newId, lastId)
			return
		}
		lastId = newId
	}
}

func TestNewClusterIdGenerator(t *testing.T) {
	var err error
	var ID int64
	// 启动一个模拟的 Redis 服务器
	mrs, err := miniredis.Run()
	if err != nil {
		t.Fatalf("无法启动模拟 Redis 服务器: %v", err)
	}
	defer mrs.Close()

	redisCli := redis.NewClient(&redis.Options{
		Addr: mrs.Addr(),
	})

	sf, err := NewClusterIdGenerator(redisCli, "test")
	assert.Nil(t, err)
	assert.NotNil(t, sf)

	m := make(map[int64]struct{})
	for i := 0; i < 100; i++ {
		ID, err = sf.NextId(context.Background())
		assert.Nil(t, err)
		m[ID] = struct{}{}
	}

	assert.Len(t, m, 100)
}
