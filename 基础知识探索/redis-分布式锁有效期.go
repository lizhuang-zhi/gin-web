package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// 模拟业务处理函数
func processBusinessLogic() {
	time.Sleep(5 * time.Second) // 模拟耗时操作
}

// 获取分布式锁
func acquireLock(ctx context.Context, rdb *redis.Client, lockKey string, value string, expiration time.Duration) bool {
	success, err := rdb.SetNX(ctx, lockKey, value, expiration).Result()
	if err != nil {
		log.Printf("获取锁失败: %v\n", err)
		return false
	}
	return success
}

// 检查锁是否存在
func checkLock(ctx context.Context, rdb *redis.Client, lockKey string) bool {
	exists, err := rdb.Exists(ctx, lockKey).Result()
	if err != nil {
		log.Printf("检查锁失败: %v\n", err)
		return false
	}
	return exists > 0
}

// 获取锁的剩余过期时间
func getLockTTL(ctx context.Context, rdb *redis.Client, lockKey string) time.Duration {
	ttl, err := rdb.TTL(ctx, lockKey).Result()
	if err != nil {
		log.Printf("获取锁TTL失败: %v\n", err)
		return 0
	}
	return ttl
}

// 模拟业务处理的函数
func handleBusiness(ctx context.Context, rdb *redis.Client, workerId int) {
	lockKey := "business_lock"
	value := fmt.Sprintf("worker_%d", workerId)
	expiration := 3 * time.Second // 锁的过期时间为3秒

	// 尝试获取锁
	if !acquireLock(ctx, rdb, lockKey, value, expiration) {
		log.Printf("Worker %d 无法获取锁\n", workerId)
		// 检查锁是否仍然存在，并显示剩余时间
		if checkLock(ctx, rdb, lockKey) {
			ttl := getLockTTL(ctx, rdb, lockKey)
			log.Printf("锁仍然存在，剩余时间: %v，可能存在死锁情况\n", ttl)
		}
		return
	}

	log.Printf("Worker %d 成功获取锁\n", workerId)

	processBusinessLogic()

	// 以下代码由于panic的发生将永远不会执行
	log.Printf("Worker %d 正常完成业务处理\n", workerId)
	_, err := rdb.Del(ctx, lockKey).Result()
	if err != nil {
		log.Printf("释放锁失败: %v\n", err)
	}
}

func main() {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	// 测试Redis连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("无法连接到Redis: %v", err)
	}

	// 确保开始时锁是清除的
	rdb.Del(ctx, "business_lock")

	// 启动多个worker来模拟并发情况
	var wg sync.WaitGroup

	// 启动多个worker来模拟并发情况
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		workerID := i
		go func() {
			defer wg.Done()
			// 使用recover来防止panic导致程序完全退出
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Worker %d 发生异常: %v\n", workerID, r)
				}
			}()
			handleBusiness(ctx, rdb, workerID)
		}()
		// 间隔2秒启动下一个worker，这样可以更清楚地观察锁的状态
		time.Sleep(2 * time.Second)
	}

	// 等待所有worker完成
	wg.Wait()

	// 最后检查锁的状态
	if checkLock(ctx, rdb, "business_lock") {
		ttl := getLockTTL(ctx, rdb, "business_lock")
		log.Printf("程序结束时锁仍然存在，剩余时间: %v\n", ttl)
	} else {
		log.Printf("程序结束时锁已经被清除\n")
	}
}
