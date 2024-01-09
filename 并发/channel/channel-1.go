package main

import (
	"fmt"
	"time"
)

/*
需求: 现在要计算 1-200 的各个数的阶乘，并且把各个数的阶乘放入到 map 中。最后显示出来。要求使用 goroutine 完成
*/

var (
	dataMap = make(map[int]int, 10) // 初始化map大小
)

func main() {
	for i := 1; i <= 200; i++ {
		go buildMap(i)
	}

	time.Sleep(time.Second * 10) // 延缓主线程结束时间

	// 遍历map
	for key, val := range dataMap {
		fmt.Printf("map[%d]%d\n", key, val)
	}
}

// 计算阶乘, 并将数据存入map中
func buildMap(num int) {
	// 阶乘计算
	factorialRes := 1
	for i := 1; i <= num; i++ {
		factorialRes *= i
	}

	// 写入map
	dataMap[num] = factorialRes
}

/*
	代码解读:
	1. 程序执行无任何输出: 如果不添加“time.Sleep(time.Second * 10) // 延缓主线程结束时间”这段代码,
	则程序执行无任何输出, 因为执行的第一个for循环的协程需要花费大量时间, 而主线程已经执行完毕
	(第二个for循环也执行了,只是没有值), 所以没有任何输出

	2. 基于第一个问题, 所以添加“// time.Sleep(time.Second * 10) // 延缓主线程结束时间”这段代码
	但此时运行代码, 出现了第二个问题:
	fatal error: concurrent map writes
	fatal error: concurrent map writes
	这是典型的资源写入的竞争问题(多个协程同时往dataMap中写入)

	tips: channel-1.exe文件是通过运行“go build -race channel-1.go”生成的
*/

// 后续复习练习代码：

// var showMap map[int]int = make(map[int]int, 20)
// var wg sync.WaitGroup
// var lock sync.Mutex

// func main() {
// 	wg.Add(20)

// 	for i := 1; i <= 20; i++ {
// 		go storeNumToMap(i)
// 	}

// 	wg.Wait()

// 	for key, val := range showMap {
// 		fmt.Printf("key为%d, 对应的val为%d\n", key, val)
// 	}
// }

// func factorial(num int) int {
// 	if num == 0 || num == 1 {
// 		return 1
// 	}
// 	return num * factorial(num-1)
// }

// func storeNumToMap(num int) {
// 	defer wg.Done()
// 	lock.Lock()
// 	showMap[num] = factorial(num)
// 	lock.Unlock()
// }
