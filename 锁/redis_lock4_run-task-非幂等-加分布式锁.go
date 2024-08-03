package main

import (
	"booking-app/shttp/shttp"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

var pools = []redsync.Pool{
	&redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	},
}

var rs = redsync.New(pools)
var redisConn = pools[0].Get()

func main() {
	go ginServer() // 启动gin服务

	time.Sleep(1 * time.Second) // 等待1秒,保证gin服务启动

	DistributedTasks() // 模拟分布式任务

	time.Sleep(2 * time.Second)
}

func DistributedTasks() {
	// 启动3个goroutine执行任务, 模拟分布式环境下的任务执行(可以简单理解为三个pod,同时执行相同的任务)
	for i := 0; i < 3; i++ {
		go task(i + 1)
	}
}

// 任务
func task(id int) {
	fmt.Printf("任务[%v]执行开始>>\n", id)

	to := "leo"

	mutex := rs.NewMutex("mail-lock-to-" + to)
	if err := mutex.Lock(); err != nil {
		fmt.Printf("任务[%v]未能获取锁: %v\n", id, err)
		return
	}

	// 检查邮件是否发送
	sent, err := redis.Bool(redisConn.Do("GET", "mail-sent-to-"+to))
	if err != nil && err != redis.ErrNil {
		fmt.Printf("任务[%v]检查邮件发送状态失败: %v\n", id, err)
		return
	}
	if sent {
		fmt.Printf("任务[%v]邮件已发送，跳过此任务\n", id)
		mutex.Unlock()
		return
	}

	// 访问Gin服务
	resp := shttp.NewHttpRequest(context.Background(), `http://localhost:7748/send-mail?to=`+to, "GET", nil).Do()
	if resp.Err != nil {
		fmt.Printf("req.Do error: %v", resp.Err)
		return
	}

	// 标记已经发送过邮件
	_, err = redisConn.Do("SET", "mail-sent-to-"+to, true)
	if err != nil {
		fmt.Printf("任务[%v]标记邮件发送状态失败: %v\n", id, err)
	}

	mutex.Unlock()

	fmt.Printf("任务[%v]执行结束!!\n", id)
}

func ginServer() {
	router := gin.Default()
	router.GET("/send-mail", func(c *gin.Context) {
		// 获取参数
		to := c.Query("to")

		fmt.Println("send mail to:", to)

		c.JSON(200, gin.H{
			"message": "send mail" + to,
		})
	})

	router.Run(":7748")
}
