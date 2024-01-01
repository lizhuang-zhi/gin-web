package main

import "fmt"

/*
管道的遍历:
1. 管道支持for-range方式遍历, 不要采用普通的for循环, 因为普通for循环, 管道遍历时len在变化, 会出现问题
2. 在遍历时, 如果channel没有关闭, 会出现deadlock的错误
3. 在遍历时, 如果channel已经关闭, 会正常遍历数据, 遍历完后, 会退出遍历
*/

func main() {
	intChan := make(chan int, 100)

	// 往管道中塞入100个数据
	for i := 0; i < 100; i++ {
		intChan <- i * 2
	}

	// 关闭管道
	close(intChan)

	// 遍历管道
	for v := range intChan {
		fmt.Println(v) // 如果不关闭管道, 就进行遍历, 会出现deadlock的错误
	}
}
