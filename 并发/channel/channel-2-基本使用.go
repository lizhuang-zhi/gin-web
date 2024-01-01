package main

import "fmt"

/*
channel说明:
1. channel是引用类型
2. channel必须初始化才能使用, 也就是make之后才能使用
3. channel是有类型的, 一个string的channel只能存放string类型数据

总结:
1. 管道写入时,不能超过其容量cap
2. 不过可以通过不断的读取数据,来释放容量,腾出空间
3. 当不使用协程时, 如果channel数据已经全部取出,再取就会报错: fatal error: all goroutines are asleep - deadlock!
*/

func main() {
	var intChan chan int
	intChan = make(chan int, 3)

	fmt.Printf("intChan的值: %v, intChan本身的地址: %p\n", intChan, &intChan)

	// 向管道写入数据
	intChan <- 10
	num := 211
	intChan <- num
	intChan <- 50

	// intChan <- 60 // 写入数据超过容量, 会报错: fatal error: all goroutines are asleep - deadlock!

	// 查看管道的长度和容量
	fmt.Printf("channel len=%v cap=%v\n", len(intChan), cap(intChan)) // channel len=3 cap=3

	// 在没有使用协程的情况下, 如果管道的数据已经全部取出,再取就会报错: fatal error: all goroutines are asleep - deadlock!
	num2 := <-intChan
	num3 := <-intChan
	num4 := <-intChan
	fmt.Println("num2=", num2, "num3=", num3, "num4=", num4)

	fmt.Printf("channel len=%v cap=%v\n", len(intChan), cap(intChan)) // channel len=0 cap=3

	// num5 := <-intChan  // 此时再取就会报错: fatal error: all goroutines are asleep - deadlock!
}
