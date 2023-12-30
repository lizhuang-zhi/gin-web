package main

import (
	"fmt"
	"time"
)

/*
	Go协程的特点:
	1. 有独立的栈空间
	2. 共享程序堆空间
	3. 调度由用户控制
	4. 协程是轻量级的线程
*/

/*
task1:
在主线程中创建一个goroutine, 该协程每隔1秒输出"hello world"
在主线程中也每隔1秒输出"hello golang", 输出10次后, 退出程序
要求主线程和goroutine同时执行
*/
func main() {
	for i := 1; i <= 10; i++ {
		go func() {
			fmt.Printf("hello,world !! %v\n", i)
		}()
		fmt.Printf("hello,golang %v\n", i)
		time.Sleep(time.Second * 1)
	}
}
