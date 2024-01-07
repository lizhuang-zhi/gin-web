package main

import (
	"fmt"
	"time"
)

/*
	求1-500000的素数有多少个？
*/

func main() {
	start := time.Now().Unix()

	for num := 1; num <= 500000; num++ {
		var isPrime bool = true // 默认为素数
		for i := 2; i < num; i++ {
			if num%i == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {

		}
	}

	end := time.Now().Unix()
	fmt.Printf("传统计算耗时：%d秒\n", end-start) // 传统计算耗时：9秒
}
