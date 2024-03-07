package workerid

import (
	"context"
	"demo-redis-snowflake-distribute-id/errhandle"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	MinWorkerID = 0
	MaxWorkerID = 1024

	// 缓存过期时间=2分钟
	WorkerCacheTTL = 2 * time.Minute
	// 续租间隔时间=30秒
	WorkerKeepInterval = 30 * time.Second
)

type IDGen struct {
	workerName    string
	workerID      int64
	genTime       string
	hostname      string
	stopChan      chan struct{}
	redisOperator redis.UniversalClient
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewWorkerIDGen() *IDGen {
	// 链接redis 单节点
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("连接redis出错，错误信息：%v", err)
	}
	fmt.Println("成功连接redis")

	return &IDGen{
		workerName:    getWorkerName(),
		hostname:      getHostname(),
		stopChan:      make(chan struct{}),
		redisOperator: rdb,
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("fail to get hostname, error: %s\n", hostname)
		hostname = getRandomStr(32)
	}
	fmt.Printf("success to get hostname: %s\n", hostname)
	return hostname
}

var (
	// 数字+大小写
	bytes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	bLen  = len(bytes)
)

// 生产随机字符串
func getRandomStr(size int) string {
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 基于当前纳秒时间，做随机
	// 根据要求的长度，获取size次随机字符
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(bLen)])
	}
	// TODO: 零拷贝一下
	return string(result)
}

func getWorkerName() string {
	// 可以读取apollo配置，这里写死
	return "test-srv"
}

func (g *IDGen) GenWorkerID() (int64, error) {
	if g.workerID > MinWorkerID && g.workerID < MaxWorkerID {
		fmt.Printf("GenWorkerId: success to gen worker id: %d\n", g.workerID)
		return g.workerID, nil
	}
	var workerID int64 = MinWorkerID
	ctx := context.Background()

	for {
		// 循环询问获取
		if workerID >= MaxWorkerID {
			workerID = MinWorkerID
		}

		g.genTime = time.Now().Format("2006-01-02 15:04:05")

		/**
		redis 存储的东西：
		<k, v> = <test-srv:workerid:1023, os.Hostname或者长度32字符串:2024-01-02 14:12:12>
		*/
		redisKey := g.getWorkerRedisKey(workerID)
		redisVal := g.getWorkerRedisValue(g.hostname, g.genTime)
		ok, err := g.redisOperator.SetNX(ctx, redisKey, redisVal, WorkerCacheTTL).Result()
		if err != nil {
			fmt.Printf("GenWorkerId: fail to setnx key info to redis, error: %s\n", err.Error())
			return 0, err
		}

		if ok {
			// 抢到锁了，直接返回此workerID
			g.workerID = workerID
			break
		}

		// 抢不到锁就新生一个
		// 1）判断是否与当前host匹配
		workerVal, err := g.redisOperator.Get(ctx, redisKey).Result()
		if err != nil {
			fmt.Printf("GenWorkerId: fail to get worker key info from redis, error: %s\n", err.Error())
			return 0, err
		}
		hostname, err := parseHostname(workerVal)
		if err != nil {
			fmt.Printf("GenWorkerId: fail to parse worker value, error: %s", err.Error())
			workerID++
			continue
		}
		// 获取值 与 当前host匹配，直接续租使用（续租会更新genTime）
		if hostname == g.hostname {
			g.workerID = workerID
			break
		}
		workerID++
	}

	go g.keepalive()
	fmt.Printf("GenWorkerId: success to gen worker id: %d\n", workerID)
	return workerID, nil
}

func parseHostname(workerVal string) (string, error) {
	// 简单校验
	values := strings.Split(workerVal, ":")
	if len(values) < 2 {
		return "", errors.New("invalid worker value")
	}
	return values[0], nil
}

func (g *IDGen) getWorkerRedisKey(workerID int64) string {
	return fmt.Sprintf("%s:workerid:%d", g.workerName, workerID)
}

func (g *IDGen) getWorkerRedisValue(hostname string, timeStr string) string {
	return fmt.Sprintf("%s:%s", hostname, timeStr)
}

func (g *IDGen) keepalive() {
	defer errhandle.CatchError()

	redisKey := g.getWorkerRedisKey(g.workerID)
	redisVal := g.getWorkerRedisValue(g.hostname, g.genTime)
	ctx := context.Background()
	ticker := time.NewTicker(WorkerKeepInterval)
	for {
		select {
		case <-ticker.C:
			// 考虑到：保活一直失败，直接重启，否则可能冲突
			_, err := g.redisOperator.Set(ctx, redisKey, redisVal, WorkerCacheTTL).Result()
			if err != nil {
				fmt.Printf("fail to set worker keepalive, worker key: %s, worker value: %s\n", redisKey, redisVal)
			}
			fmt.Printf("WorkerIDTask: %v\n", redisKey)
		case <-g.stopChan:
			ticker.Stop()
			return
		}
	}
}
