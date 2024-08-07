package main

import (
	"booking-app/shttp/gin_server"
	"booking-app/shttp/shttp"
	"context"
	"fmt"
	"io"
	"time"
)

func main() {
	// 异步启动gin服务
	go gin_server.InitGinServer()

	time.Sleep(1 * time.Second) // 等待1秒，保证gin服务启动

	// TestRequestBaidu()

	// TestRequestGinServerGet()

	TestRequestGinServerPost()
	// 保持程序运行
	select {}
}

func TestRequestBaidu() {
	resp := shttp.Get(context.Background(), "http://www.baidu.com", nil)
	if resp.Err != nil {
		fmt.Printf("req.Do error: %v", resp.Err)
		return
	}

	fmt.Println("resp.StatusCode:", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("io.ReadAll error: %v", err)
		return
	}
	fmt.Println("resp.Body:", string(body))
	fmt.Println("resp.Header:", resp.Header)
}

func TestRequestGinServerGet() {
	resp := shttp.Get(context.Background(), "http://localhost:7878/index", nil)
	if resp.Err != nil {
		fmt.Printf("req.Do error: %v", resp.Err)
		return
	}

	fmt.Println("resp.StatusCode:", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("io.ReadAll error: %v", err)
		return
	}
	fmt.Println("resp.Body:", string(body))
	fmt.Println("resp.Header:", resp.Header)
}

func TestRequestGinServerPost() {
	data := `{"name": "张三", "age": 18}`

	resp := shttp.Post(
		context.Background(), "http://localhost:7878/insert", []byte(data),
		shttp.WithHTTPRetry(3),               // 配置重试options
		shttp.WithHTTPTimeout(4*time.Second), // 配置超时options
	)

	respData, err := resp.ReadAll()
	if err != nil {
		fmt.Printf("req.Do error: %v", err)
		return
	}

	fmt.Println("resp.Body:", string(respData))
	fmt.Println("resp.Status:", resp.Status())
}
