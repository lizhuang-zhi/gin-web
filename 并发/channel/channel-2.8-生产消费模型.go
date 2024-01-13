package main

import (
	"fmt"
	"sync"
	"time"
)

/*
生产消费模型: 简单的理解就是从一个资源池(管道)中获取数据, 并进行消费
*/

// 生产
func producer(sourceChan chan<- *AliYunLog, wg *sync.WaitGroup) {
	defer wg.Done()

	// 产生日志
	produceLog := &AliYunLog{
		EndPoint:   1,
		Message:    "日志内容123",
		Level:      2,
		CreateTime: time.Now().Nanosecond(),
	}
	produceLog2 := &AliYunLog{
		EndPoint:   13,
		Message:    "Erorr Panic",
		Level:      4,
		CreateTime: time.Now().Nanosecond(),
	}

	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			sourceChan <- produceLog2
			continue
		}
		sourceChan <- produceLog
	}

	close(sourceChan)
}

// 消费
func consumer(sourceChan <-chan *AliYunLog, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		log, ok := <-sourceChan
		if !ok {
			break
		}
		// 打印日志
		fmt.Printf("日志端点: %d, 日志信息: %s, 日志级别: %d, 创建时间: %d\n", log.EndPoint, log.Message, log.Level, log.CreateTime)
	}
}

type AliYunLog struct {
	EndPoint   int
	Message    string
	Level      int
	CreateTime int
}

func main() {
	var sourceChan chan *AliYunLog = make(chan *AliYunLog, 5)

	var wg sync.WaitGroup

	wg.Add(1)
	go producer(sourceChan, &wg)

	// 多消费者(开启多个协程)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go consumer(sourceChan, &wg)
	}

	wg.Wait()

	fmt.Println("消费完所有生产内容.....")

}
