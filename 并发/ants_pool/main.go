package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func main() {
	// 创建一个 goroutine 池，设置最大并发数为 5
	pool, err := ants.NewPool(5, ants.WithNonblocking(false))
	if err != nil {
		fmt.Println("创建 goroutine 池失败:", err)
		return
	}
	defer pool.Release()

	// 创建等待任务完成的 WaitGroup
	var wg sync.WaitGroup

	// 提交任务到 goroutine 池
	for i := 0; i < 10; i++ {
		wg.Add(1)
		taskID := i
		if err := pool.Submit(func() {
			defer wg.Done()
			fmt.Printf("任务 %d 开始执行...\n", taskID)
			time.Sleep(3 * time.Second) // 模拟任务执行
			fmt.Printf("任务 %d 执行完毕\n", taskID)
		}); err != nil {
			// 处理提交任务时发生的错误
			fmt.Println("提交任务失败:", err)
		}
	}

	// 等待所有任务完成
	wg.Wait()
	fmt.Println("所有任务已完成")
}
