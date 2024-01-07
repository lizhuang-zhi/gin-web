package main

import "fmt"

/*
	// 管道默认是双向的
	var chan1 chan int // 可读可写
	chan1 = make(chan int, 3)
	chan1 <- 20
	num := <-chan1
	fmt.Println(num) // 20

	// 声明只写的channel
	var chan2 chan<- int
	chan2 = make(chan int, 3)
	chan2 <- 20
	// num2 := <-chan2 // 报错：invalid operation: <-chan2 (receive from send-only type chan<- int)

	// 声明只读的channel
	// var chan3 <-chan int
	// num3 := <-chan3  // 报错：invalid operation: chan3 <- (send to receive-only type <-chan int)
	// fmt.Println(num3) // 先关闭, 打开会报错死锁
*/

func main() {
	// 案例
	var ch chan int = make(chan int, 10)                // 可读可写的channel
	var exitChan chan struct{} = make(chan struct{}, 2) // 用于退出主线程的channel

	go send(ch, exitChan) // 只往ch里写
	go rece(ch, exitChan) // 只从ch里读

	for i := 0; i < cap(exitChan); i++ {
		<-exitChan
	}
	close(exitChan)
	fmt.Println("主线程退出.....")
}

// ch chan<- int 设置只往ch中写
func send(ch chan<- int, exitChan chan struct{}) {
	for i := 0; i < cap(ch); i++ {
		ch <- i
	}
	close(ch)
	var a struct{}
	exitChan <- a
}

// 设置只从ch中读
func rece(ch <-chan int, exitChan chan struct{}) {
	for {
		val, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("从ch中读: %d\n", val)
	}
	var b struct{}
	exitChan <- b
}
