package main

import (
	"fmt"
	"time"
)

/*
	在之前的代码中, 如果不关闭管道, 就从管道中取数据, 会出现deadlock的错误
	在实际开发中, 使用select来解决
*/

func main() {
	intChan := make(chan int, 10)
	for i := 0; i < cap(intChan); i++ {
		intChan <- i
	}

	stringChan := make(chan string, 5)
	for i := 0; i < cap(stringChan); i++ {
		stringChan <- "Hello World" + fmt.Sprintf("%d", i)
	}
	// 在之前的代码中, 如果不关闭管道, 就从管道中读数据, 会由于管道阻塞, 而导致死锁
	// label:
	for {
		select {
		case v := <-intChan:
			fmt.Printf("从intChan管道中读数据%d\n", v)
			time.Sleep(time.Second)
		case v := <-stringChan:
			fmt.Printf("从stringChan管道中读数据%s\n", v)
			time.Sleep(time.Second)
		default:
			fmt.Println("读取不到数据了, 不玩了")
			time.Sleep(time.Second)
			// 退出方式1:
			return
			// 退出方式2:(不推荐)
			// break label
		}
	}
}
