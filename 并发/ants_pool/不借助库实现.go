package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

var goRoutinePoolChan = make(chan bool, 2) // 设置最大并发数3

func main() {
	fmt.Println("不借助库实现控制并发个数【通过管道】")

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		goRoutinePoolChan <- true // 管道塞入数据，记录当前执行情况

		go longTimeTask(i)
	}

	wg.Wait()
	close(goRoutinePoolChan)
}

func longTimeTask(num int) {
	defer wg.Done()

	fmt.Println("开始执行任务", num)
	time.Sleep(3 * time.Second)
	fmt.Println(".......结束执行任务", num)

	<-goRoutinePoolChan // 任务执行完释放一个管道空间
}
