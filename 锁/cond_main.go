package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mu     sync.Mutex
	cond   = sync.NewCond(&mu)
	count  int
	maxCap = 10
)

// Producer 生产者
func producer() {
	for {
		mu.Lock()
		for count == maxCap {
			fmt.Println("Full, waiting...")
			cond.Wait() // 等待消费者消耗
		}
		count++
		fmt.Printf("Produced, count = %d\n", count)
		cond.Signal() // 通知消费者
		mu.Unlock()
		// 可以通过调整生产者生产速度和消费者消费速度来测试
		time.Sleep(time.Millisecond * 100) // 模拟更快的生产时间
		// time.Sleep(time.Millisecond * 500) // 模拟生产时间
	}
}

// Consumer 消费者
func consumer() {
	for {
		mu.Lock()
		for count == 0 {
			fmt.Println("Empty, waiting...")
			cond.Wait() // 等待生产者生产
		}
		count--
		fmt.Printf("Consumed, count = %d\n", count)
		cond.Signal() // 通知生产者
		mu.Unlock()
		time.Sleep(time.Millisecond * 500) // 模拟消费时间
	}
}

func main() {
	go producer()
	go consumer()

	// 运行一段时间后退出
	time.Sleep(time.Second * 10)
}
