package main

import (
	"fmt"
	"my-distributed-service/common/broadcast"
	"my-distributed-service/common/dlock"
	"net/http"
	"sync"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	curState   string     // 当前节点的状态数据
	stateMutex sync.Mutex // 状态的互斥锁
)

// 把状态改变的逻辑封装到一个函数中
func changeStateAndBroadcast(newState string, broadcaster *broadcast.Broadcaster) {
	// 试图获取锁
	mutex, err := dlock.ObtainLock("my-lock", 10)
	if err != nil {
		log.Printf("Failed to obtain lock: %v\n", err)
		return
	}
	defer dlock.ReleaseLock(mutex)

	// 加锁以同步状态的变更
	stateMutex.Lock()
	curState = newState
	stateMutex.Unlock()

	// 发布新状态到所有节点
	err = broadcaster.Publish("state-channel", newState)
	if err != nil {
		log.Printf("Failed to publish state: %v\n", err)
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// 初始化分布式锁
	dlock.InitDistruibutedLock(redisClient)
	// 初始化广播器
	broadcaster := broadcast.NewBroadcaster(redisClient)

	// 订阅频道以接收同步状态数据的通知
	pubSub := broadcaster.Subscribe("state-channel")
	defer pubSub.Close()

	go func() {
		ch := pubSub.Channel()
		for msg := range ch {
			fmt.Printf("Received message: %s\n", msg.Payload)
			// 收到消息后，将状态数据同步到本地
			curState = msg.Payload
			fmt.Printf("Current state: %s\n", curState)
		}
	}()

	// 启动Gin路由并在端口8080监听
	router := gin.Default()
	// 添加一个路由，用于改变状态
	router.POST("/change-state", func(c *gin.Context) {
		newState := c.PostForm("state")
		if newState == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "State not provided"})
			return
		}

		go changeStateAndBroadcast(newState, broadcaster) // 同步广播
		c.JSON(http.StatusOK, gin.H{"message": "State change initiated"})
	})
	// 启动Gin服务
	go router.Run(":8860")

	// 阻塞主进程
	select {}
}
