package main

import (
	"fmt"
	"runtime"
)

/*
 go设置cpu数目
*/

func main() {
	// 查看电脑CPU个数
	cpuNum := runtime.NumCPU()
	fmt.Println(cpuNum) // 8

	// 设置最大CPU数量
	runtime.GOMAXPROCS(cpuNum - 1)
}
