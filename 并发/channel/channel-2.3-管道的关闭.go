package main

/*
close() 关闭管道: 关闭后不能再写入数据，但是可以读取数据
*/

func main() {
	intChan := make(chan int, 3)
	intChan <- 1
	intChan <- 2

	// 关闭管道: 关闭后不能再写入数据，但是可以读取数据
	close(intChan)

	// 读取管道数据
	n1, ok := <-intChan
	println(n1, ok) // 1 true

	// 继续写入报错
	// intChan <- 3 // panic: send on closed channel
}
