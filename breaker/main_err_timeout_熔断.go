package main

import (
	"booking-app/breaker/breaker"
	"booking-app/breaker/gin_server"
	"booking-app/breaker/hystrix"
	"booking-app/shttp/shttp"
	"context"
	"fmt"
	"time"
)

var GlobalData string

func main() {
	// 异步启动gin服务
	go gin_server.InitGinServer()

	time.Sleep(1 * time.Second) // 等待1秒，保证gin服务启动

	// 初始化熔断器配置
	config := &hystrix.Config{
		Name: "test",
		CommandConfig: hystrix.CommandConfig{
			Timeout:               1000,
			MaxConcurrentRequests: 100,
			ErrorPercentThreshold: 50,
		},
	}

	// 创建熔断器
	myBreaker := hystrix.NewBreaker(config)

	// 创建任务函数
	myTask := func(ctx context.Context) error {
		// 设置请求超时时间为3秒
		resp := shttp.Get(context.Background(),
			"http://localhost:7878/timeout",
			nil,
			shttp.WithHTTPTimeout(3*time.Second))
		if resp.Err != nil {
			// fmt.Printf("req.Do error: %v", resp.Err)
			return breaker.ErrTimeout // 返回超时错误(直接熔断, 防止系统雪崩, 不执行降级操作)
		}

		// 将获取的输入存入全局内存中
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		GlobalData = string(buf[:n])

		return nil
	}

	// 创建 fallback 函数
	myFallback := func(ctx context.Context, err error) error {
		fmt.Println("主请求超时，降级操作, 返回用户默认的数据")
		GlobalData = "default data" // 降级操作，返回默认数据
		return err
	}

	err := myBreaker.Do(context.Background(), myTask, myFallback)
	fmt.Println("报错:", err)
	fmt.Println("全局数据GlobalData:", GlobalData)
}
