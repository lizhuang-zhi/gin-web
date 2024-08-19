package main

import (
	"fmt"
	"time"
)

/*
我的系统中有两种类型的定时器，分别是"每三分钟执行任务1和任务2“、"每10分钟执行任务2“，
而任务1是从飞书读取数据，然后用数据去请求第三方API，然后将返回的数据存储到数据库；
而任务2是从P4V的API获取数据，然后用数据去请求第三方API，然后将返回的数据计算下，然后调用飞书API发送给用户；
而任务3是从飞书的多维表格获取数据，然后用数据去请求第三方API，不同步等待三方API的返回，因为此时访问的第三方API在执行完对应的任务后，会请求本服务外放出去的API，本服务外放出去的API会将请求的数据存入数据库；
*/
type Event struct {
	Data string
}

type EventChannel chan Event

type EventHandler func(Event)

type EventDispatcher struct {
	Handlers map[string]EventChannel
}

func (d *EventDispatcher) Register(eventName string, handler EventHandler) {
	if _, found := d.Handlers[eventName]; !found {
		d.Handlers[eventName] = make(EventChannel)
	}
	go func(c EventChannel) {
		for e := range c {
			handler(e)
		}
	}(d.Handlers[eventName])
}

func (d *EventDispatcher) Dispatch(eventName string, data string) {
	if c, found := d.Handlers[eventName]; found {
		c <- Event{Data: data}
	}
}

func main() {
	d := EventDispatcher{
		Handlers: map[string]EventChannel{},
	}

	d.Register("task1", func(e Event) {
		fmt.Println("task1:", e.Data)
		time.Sleep(1 * time.Second)
	})

	d.Register("task2", func(e Event) {
		fmt.Println("task2:", e.Data)
		time.Sleep(1 * time.Second)
	})

	d.Register("task3", func(e Event) {
		fmt.Println("task3:", e.Data)
		time.Sleep(1 * time.Second)
	})

	go func() {
		ticker := time.NewTicker(3 * time.Minute)
		for range ticker.C {
			d.Dispatch("task1", "Feishu data for Task1")
			d.Dispatch("task2", "P4V data for Task2")
		}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			d.Dispatch("task2", "P4V data for Task2")
		}
	}()

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		for now := range ticker.C {
			if now.Hour() == 10 {
				d.Dispatch("task3", "Feishu data for Task3")
			}
		}
	}()

	select {}
}
