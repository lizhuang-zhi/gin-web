package main

import (
	"fmt"
	"sync"
	"time"
)

/*
1. [当管道的容量(10)小于写入的数量(50)时], 如果只向管道中写入数据, 而没有读, 则会出现阻塞而dead lock; 写入数据大于等于管道容量,则不会dead lock
2. 从管道中读取的速度小于写入的速度, 不会出现死锁
*/

var (
	intChan = make(chan int, 10) // 降低intChan的容量
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

	for v := range intChan {
		fmt.Println(v)
		time.Sleep(time.Second * 1) // 从管道中读的速度变慢
	}
}
