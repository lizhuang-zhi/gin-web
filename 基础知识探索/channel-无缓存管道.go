package main

import (
	"fmt"
	"time"
)

func senderWithBuffered(ch chan int) {
	// 往有缓存的管道里发送数据
	for i := 1; i <= 3; i++ {
		ch <- i
		fmt.Printf("发送方往有缓存管道发送了数据: %d\n", i)
	}
}

func receiverWithBuffered(ch chan int) {
	time.Sleep(2 * time.Second) // 模拟接收方延迟启动，体现缓存优势
	for i := 0; i < 3; i++ {
		data := <-ch
		fmt.Printf("接收方从有缓存管道接收到数据: %d\n", data)
	}
}

func main() {
	chBuffered := make(chan int)

	go senderWithBuffered(chBuffered)
	go receiverWithBuffered(chBuffered)

	time.Sleep(5 * time.Second) // 让主协程等待一段时间，保证子协程有足够时间执行完
	fmt.Println("有缓存管道相关操作执行完毕")
}
