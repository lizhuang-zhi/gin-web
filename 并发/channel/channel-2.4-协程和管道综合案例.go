package main

import (
	"fmt"
	"sync"
)

/*
需求:
1. 开启一个writeData协程, 向管道intChan中写入50个整数
2. 开启一个readData协程, 从管道intChan中读取writeData写入的数据[注意这里可能是边写边读]
3. 注意: writeData和readData操作的是同一个管道
4. 主线程需要等待writeData和readData协程都完成工作才能退出管道
*/

var (
	intChan = make(chan int, 50)
	wg      sync.WaitGroup
)

func main() {
	wg.Add(1)
	go writeData()

	wg.Add(1)
	go readData()

	wg.Wait() // 等待协程执行完毕
}

// 写入数据到管道
func writeData() {
	defer wg.Done()
	defer close(intChan) // 数据全部录入后关闭管道

	for i := 0; i < 50; i++ {
		intChan <- i
		fmt.Printf("写入数据: %d\n", i)
	}
}

// 读取管道数据
func readData() {
	defer wg.Done()

	// 方式一:
	for v := range intChan {
		fmt.Println(v)
	}

	// // 方式二:(需要保证writeData中关闭管道)
	// for {
	// 	x, ok := <-intChan
	// 	if !ok {
	// 		break
	// 	}
	// 	fmt.Println(x)
	// }
}
