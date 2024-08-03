package main

import (
	"booking-app/shttp/shttp"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

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

	// 访问Gin服务
	resp := shttp.NewHttpRequest(context.Background(), `http://localhost:7748/send-mail?to=`+to, "GET", nil).Do()
	if resp.Err != nil {
		fmt.Printf("req.Do error: %v", resp.Err)
		return
	}

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
