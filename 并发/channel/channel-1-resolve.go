package main

import (
	"fmt"
	"sync"
	"time"
)

/*
解决 channel-1.go 文件中的需求问题

解决方式
1. 全局变量加锁(互斥锁): 当单个协程在写入dataMap时,加锁
*/

var (
	dataMap = make(map[int]int, 10) // 初始化map大小
	// sync包: 大部分都是适用于低水平程序线程, 高水平的同步使用channel通信更好一些
	lock sync.Mutex // 全局互斥锁(写锁), Mutex是一个互斥锁
)

func main() {
	for i := 1; i <= 20; i++ {
		go buildMap(i)
	}

	time.Sleep(time.Second * 3) // 延缓主线程结束时间

	/*
		这里读操作也要添加互斥锁: 因为程序设计上可能知道10s执行完所有协程,
		但是主线程并不知道,所以底层仍然可能出现资源争夺

		不加锁的话, 生成的channel-1-resolve.exe报错: Found 2 data race(s)
	*/
	lock.Lock()
	// 遍历map
	for key, val := range dataMap {
		fmt.Printf("map[%d]%d\n", key, val)
	}
	lock.Unlock()
}

// 计算阶乘, 并将数据存入map中
func buildMap(num int) {
	factorialRes := 1
	for i := 1; i <= num; i++ {
		factorialRes *= i
	}

	lock.Lock() // 只允许单个协程进行写入, 避免资源竞争
	dataMap[num] = factorialRes
	lock.Unlock() // 解锁
}
