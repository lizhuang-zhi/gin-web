package main

import "fmt"

/*
	求1-8000的素数有多少个？

	解决思路：
	1. 开启一个协程，用于存放1-8000的数，叫intChan
	2. 再开4个协程，用于从intChan管道中取数据
	3. 并判断取出来的数是否为素数，如果为素数，则放入primeChan管道中
	4. 开启的四个协程，执行完一个就往exitChan中输入一个true，目的是为了知道4个异步的协程什么时候执行完毕
	5. 在主线程的最后，执行一个for循环，不断从exitChan中取，直到取不出来，则为完成
*/
// 将数据存放到intChan管道
func putIntChan(intChan chan int) {
	for i := 1; i <= 80; i++ {
		intChan <- i
	}
	close(intChan)
}

// 从intChan中获取数据并判断是否为素数, 是则存入primeChan中
func getDataToJudge(intChan, primeChan chan int, exitChan chan bool) {
	for {
		num, ok := <-intChan
		if !ok {
			fmt.Println("intChan管道中已经取完，没有数了！！")
			break // 跳出循环，intChan管道中已经取完，没有数了
		}
		var isPrime bool = true // 默认为素数
		for i := 2; i < num; i++ {
			if num%i == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {
			primeChan <- num // num为素数，则存入primeChan管道
		}
	}

	// 执行完毕，则标记当前的计算协程执行完毕
	exitChan <- true
}

// 判断开始的4个协程是否执行完毕
func isFinish(exitChan chan bool, primeChan chan int) {
	// 方式一：会出现死锁：因为我其他地方没有去关闭exitChan，导致这里会一直等待exitChan管道，所以没法执行到程序最后的close(primeChan)，
	// 所以导致主进程中的prime, ok := <-primeChan，也会因为primeChan无法关闭而一直等待，导致死锁
	// for {
	// 	_, ok := <-exitChan
	// 	if !ok {
	// 		break
	// 	}
	// }

	// 方式二：不会出现死锁，因为这里通过循环次数来控制，只有在接收到指定次数的值后，主协程才会退出。这种方式不依赖于通道是否关闭，而是通过循环次数的约束来确定何时退出
	for i := 0; i < cap(exitChan); i++ {
		<-exitChan
	}

	fmt.Println("【方式一，无法到这里， 但是方式二可以到这里】")
	close(primeChan)
}

func main() {
	var intChan chan int = make(chan int, 1000)   // 流动管道
	var primeChan chan int = make(chan int, 2000) // 存储所有素数的管道
	var exitChan chan bool = make(chan bool, 4)   // 用于标记多个协程是否都执行完毕

	go putIntChan(intChan)

	for i := 0; i < cap(exitChan); i++ { // 之所以开启四个协程去判断素数，是为了体现多协程的异步计算，提高效率
		go getDataToJudge(intChan, primeChan, exitChan)
	}

	go isFinish(exitChan, primeChan)

	// 遍历primieChan管道，展示存入的素数
	for {
		prime, ok := <-primeChan
		if !ok {
			break
		}
		fmt.Printf("素数 = %d\n", prime)
	}

	fmt.Println("主程序执行完毕....")
}
