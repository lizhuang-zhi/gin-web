package main

import (
	"fmt"
)

/*
	计算1-8000的素数有哪些？（使用gorountine）
	1. 使用intChan存储1-8000的素数
	2. 开启4个协程, 从intChan中读取数据，然后判断是否为素数，然后将素数存入primeChan中
	3. 每个协程执行完毕，则往exitChan中存入一个true
	4. 主线程中循环等待exitChan全部取出后退出
	5. 遍历primeChan中的素数
*/

func main() {
	var intChan chan int = make(chan int, 1000)
	var primeChan chan int = make(chan int, 3000)
	var exitChan chan bool = make(chan bool, 4)

	go inputNum(intChan)

	for i := 0; i < cap(exitChan); i++ {
		go judgePrime(intChan, primeChan, exitChan)
	}

	for i := 0; i < cap(exitChan); i++ {
		<-exitChan
	}
	close(primeChan)
	for {
		num, ok := <-primeChan
		if !ok {
			break
		}
		fmt.Printf("素数: %d\n", num)
	}

	fmt.Println("退出主线程............")
}

// 只往intChan中写入数据
func inputNum(intChan chan<- int) {
	defer close(intChan)
	for i := 1; i <= 8000; i++ {
		intChan <- i
	}
}

// 判断从intChan中取出的是否为素数
func judgePrime(intChan <-chan int, primeChan chan<- int, exitChan chan<- bool) {
	for {
		num, ok := <-intChan
		if !ok {
			break
		}
		var isPrime = true // 默认为素数
		for i := 2; i < num; i++ {
			if num%i == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {
			primeChan <- num
		}
	}

	exitChan <- true
}
