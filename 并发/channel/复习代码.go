package main

import (
	"fmt"
)

/*
需求:
1. 开启一个writeData协程, 向管道intChan中写入50个整数
2. 开启一个readData协程, 从管道intChan中读取writeData写入的数据[注意这里可能是边写边读]
3. 注意: writeData和readData操作的是同一个管道
4. 主线程需要等待writeData和readData协程都完成工作才能退出管道
*/

var intChan chan int = make(chan int, 5)
var exitChan chan bool = make(chan bool, 1)

func main() {
	go writeData()
	go readData()

	for i := 0; i < cap(exitChan); i++ {
		<-exitChan
	}
	close(exitChan)

	fmt.Println("退出主线程............")
}

func writeData() {
	for i := 1; i <= 50; i++ {
		intChan <- i
	}
	close(intChan)
}

func readData() {
	for {
		val, ok := <-intChan
		if !ok {
			break
		}
		fmt.Printf("读取管道中数据：%d\n", val)
	}
	exitChan <- true
}
