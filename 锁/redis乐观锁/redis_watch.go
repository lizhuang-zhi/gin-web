package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 密码
		DB:       0,                // 数据库
		PoolSize: 100,              // 连接池大小
	})

	ctx := context.Background()

	// 初始化库存
	err := rdb.Set(ctx, "product_stock", 100, 0).Err()
	if err != nil {
		log.Fatal("初始化库存失败:", err)
	}

	// 模拟并发扣减库存
	for i := 0; i < 5; i++ {
		go func(id int) {
			if success := deductStock(ctx, rdb, id); success {
				fmt.Printf("协程%d: 扣减成功\n", id)
			} else {
				fmt.Printf("协程%d: 扣减失败\n", id)
			}
		}(i)
	}

	// 保持主线程运行
	time.Sleep(3 * time.Second)
}

// 使用 WATCH 实现库存扣减
/*
	这里一共包含6次IO操作
	1.WATCH命令
	2.GET命令
	3.事务执行（包含MULTI、DECR、EXEC命令），这里包含3次IO操作
	4.UNWATCH命令
*/
func deductStock(ctx context.Context, rdb *redis.Client, routineID int) bool {
	const maxRetries = 3 // 最大重试次数
	key := "product_stock"

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// 创建事务 Pipeline
		txf := func(tx *redis.Tx) error {
			// 读取当前库存
			currentStock, err := tx.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				return err
			}

			if currentStock <= 0 {
				return fmt.Errorf("库存不足")
			}

			// 模拟业务处理时间
			time.Sleep(100 * time.Millisecond)

			// 开启事务
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Decr(ctx, key)
				return nil
			})
			return err
		}

		// 使用 Watch 自动重试
		err := rdb.Watch(ctx, txf, key)
		switch {
		case err == nil:
			return true // 扣减成功
		case err == redis.TxFailedErr:
			fmt.Printf("协程%d-尝试%d: 版本冲突，重试...\n", routineID, attempt)
			continue // 冲突重试
		default:
			fmt.Printf("协程%d-错误: %v\n", routineID, err)
			return false
		}
	}
	return false
}
