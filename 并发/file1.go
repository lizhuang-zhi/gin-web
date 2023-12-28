package main

import (
	"fmt"
)

func main() {
	src := make(chan int)
	dest := make(chan int, 3)
	fmt.Println("程序开始执行之前..........")
	// A子协程
	go func() {
		fmt.Println("A子协程进入.............")
		defer close(src)
		for i := 0; i < 10; i++ {
			fmt.Printf("A子协程正在计算:%d\n", i)
			// 将数字发送到channel
			src <- i
		}
	}()
	// B子协程
	go func() {
		fmt.Println("B子协程进入.............")
		defer close(dest)
		for i := range src {
			fmt.Printf("B子协程正在计算:%d\n", i)
			// 将计算好的数字发送到channel
			dest <- i * i
		}
	}()
	for i := range dest {
		println(i)
	}
	fmt.Println("程序结束之后.............")
}

// 程序开始执行之前..........
// B子协程进入.............
// A子协程进入.............
// A子协程正在计算:0
// A子协程正在计算:1
// B子协程正在计算:0
// B子协程正在计算:1
// 0
// 1
// A子协程正在计算:2
// A子协程正在计算:3
// B子协程正在计算:2
// B子协程正在计算:3
// 4
// A子协程正在计算:4
// 9
// A子协程正在计算:5
// B子协程正在计算:4
// B子协程正在计算:5
// 16
// 25
// A子协程正在计算:6
// A子协程正在计算:7
// B子协程正在计算:6
// B子协程正在计算:7
// 36
// 49
// A子协程正在计算:8
// A子协程正在计算:9
// B子协程正在计算:8
// B子协程正在计算:9
// 64
// 81
// 程序结束之后.............
