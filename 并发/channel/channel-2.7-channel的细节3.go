package main

import (
	"fmt"
	"time"
)

/*
goroutine中使用recover, 解决协程中出现panic, 导致程序奔溃的问题
*/

// 协程1
func sayHello() {
	for i := 0; i < 10; i++ {
		fmt.Println("Hello, world!!")
		time.Sleep(time.Second)
	}
}

// 协程2(会出错)
func errFun() {
	// 通过recover()捕获到panic, 但是不中断其他协程和主线程的执行
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("捕获到错误: %s\n", err)
		}
	}()

	var testMap map[int]string
	testMap[0] = "error msg" // panic: assignment to entry in nil map
}

func main() {
	go sayHello()
	go errFun()

	for i := 0; i < 10; i++ {
		fmt.Printf("主线程任务: %d\n", i)
		time.Sleep(time.Second)
	}
}
