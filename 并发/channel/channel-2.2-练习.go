package main

import "fmt"

/*
 */

type Cat struct {
	Name string
	Age  int
}

func main() {
	// 定义一个可存放任意数据类型的管道
	var interfaceChan chan interface{}
	interfaceChan = make(chan interface{}, 10)

	cat1 := Cat{"tom", 8}

	interfaceChan <- 10
	interfaceChan <- "test data"
	interfaceChan <- cat1

	// 推出前两个数据
	<-interfaceChan
	<-interfaceChan

	newCat := <-interfaceChan
	fmt.Printf("newCat=%T, newCat=%v\n", newCat, newCat) // newCat=main.Cat, newCat={tom 8}
	// 这里需要使用类型断言来转换 newCat.(Cat) 表示将newCat转成Cat类型
	fmt.Printf("newCat.Name=%v\n", newCat.(Cat).Name) // newCat.Name=tom
}
