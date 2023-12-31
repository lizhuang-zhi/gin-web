package main

import (
	"fmt"
	"sync"
)

/*
解决 channel-1.go 文件中的需求问题

解决方式
2. channel(管道)

前面的channel-1-resolve.go文件中的第一种通过加锁的方式解决的不完美:
1. 主线程在等待所有的goroutine执行完毕, 但是goroutine执行的时间不确定, 设置10s, 但只是作为估算时间,
如果goroutine执行的时间超过10s, 主线程就会结束, 但是goroutine还没有执行完毕, 这样就会导致数据不完整
如果设置的时间过长, 又会导致主线程等待时间过长, 降低程序的执行效率
所以使用channel来解决这个问题!!

channel介绍:
1. channel的本质就是一个队列,
2. 数据是先进先出的
3. 线程安全, 多goroutine访问时, 不需要加锁, 就是线程安全的
4. channel是有类型的, 一个string的channel只能存放string类型数据

视频地址:  https://www.bilibili.com/video/BV1ME411Y71o?p=272&spm_id_from=pageDriver&vd_source=e339d0ca63ebabe15afcd98b996033f7
*/

var (
	dataMap = make(map[int]int, 10) // 初始化map大小
)

func main() {
	// 创建一个channel, 用于存放goroutine计算的结果
	var channelS = make(chan struct {
		num int
		res int
	}, 20)

	var wg sync.WaitGroup // 创建一个同步等待的组

	// 开启20个goroutine, 计算1-20的阶乘, 并将结果存入map中
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go buildMap(i, wg, channelS)
	}

	go func() {
		wg.Wait()

		// 关闭channel
		close(channelS)
	}()

	// 遍历通道中的数据, 将数据存入map中
	for v := range channelS {
		dataMap[v.num] = v.res
	}

	// 打印map中的数据
	for k, v := range dataMap {
		fmt.Printf("map[%d] = %d\n", k, v)
	}
}

// 计算阶乘, 并将数据存入map中
func buildMap(num int, wg sync.WaitGroup, channelS chan struct {
	num int
	res int
}) {
	defer wg.Done() // goroutine执行完毕, 就将同步等待组中的计数器减1

	// 计算阶乘
	factorialRes := 1
	for i := 1; i <= num; i++ {
		factorialRes *= i
	}

	// 将计算的结果存入channel中
	channelS <- struct {
		num int
		res int
	}{num, factorialRes}
}
