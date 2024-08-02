package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	lockKey := "counter_lock"
	counterKey := "counter"

	// 尝试获取一个锁
	resp := client.SetNX(lockKey, 1, time.Second*5)
	if resp.Val() { // 获取锁成功

		// 获取锁成功，进行业务操作
		val := client.Get(counterKey)
		counter, _ := val.Int64()
		if counter < 5 {
			counter++
			client.Set(counterKey, counter, 0)
			fmt.Println("Number is now:", counter)
		}

		// 业务操作完毕，释放锁
		client.Del(lockKey)
	} else {
		fmt.Println("Got lock failed")
	}
}
