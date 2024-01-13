package main

import (
	"fmt"
	"time"
)

/*
time.Ticker: 定时器, 用于定时执行某个任务

需求: 每隔一秒打印一次“完成任务”, 并于3秒后结束该打印
*/

func main() {
	// 创建一个1秒的定时器
	ticker := time.NewTicker(time.Second)

	// 退出管道
	exitChan := make(chan struct{})

	go tickerTime(ticker, exitChan)

	time.Sleep(time.Second * 3)
	close(exitChan)

	/*
		这里并不是想象中的执行三次就退出, 而是由于主线程和协程的执行间隙时间的不确定而导致输出存在不确定性
	*/

	fmt.Println("退出主线程.....")
}

// 定时执行任务
func tickerTime(ticker *time.Ticker, exitChan <-chan struct{}) {
	// 循环监听定时器触发的事件
	for {
		select {
		case <-ticker.C:
			// 定时器触发的处理逻辑
			fmt.Println("执行一次任务")
		case <-exitChan:
			ticker.Stop()
			fmt.Println("退出定时器.....")
			return
		}
	}
}
