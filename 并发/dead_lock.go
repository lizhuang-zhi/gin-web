package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1)

	go func() {
		defer wg.Done()

		fmt.Println("goroutine 1: acquiring lock")
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("goroutine 1: lock acquired")

		// 模拟一些处理时间
		time.Sleep(1 * time.Second)

		fmt.Println("goroutine 1: waiting for lock 2")
		mu.Lock() // 这里会导致死锁
		defer mu.Unlock()
		fmt.Println("goroutine 1: lock 2 acquired")
	}()

	// 让第一个goroutine先获取锁
	time.Sleep(500 * time.Millisecond)

	wg.Add(1)

	go func() {
		defer wg.Done()

		fmt.Println("goroutine 2: acquiring lock")
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("goroutine 2: lock acquired")

		// 模拟一些处理时间
		time.Sleep(1 * time.Second)

		fmt.Println("goroutine 2: waiting for lock 1")
		mu.Lock() // 这里会导致死锁
		defer mu.Unlock()
		fmt.Println("goroutine 2: lock 1 acquired")
	}()

	// 等待goroutines完成
	wg.Wait()

	fmt.Println("程序执行完毕")
}
