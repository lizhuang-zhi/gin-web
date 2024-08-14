package main

import (
	"fmt"
	"sync"
)

type Event struct {
	Data int64
}

type EventListener struct {
	Id         int
	Ch         chan Event
	WaitHandle *sync.WaitGroup
}

var eventBus = struct {
	sync.RWMutex
	subscribers map[int]*EventListener
}{
	subscribers: make(map[int]*EventListener),
}

func Subscribe(listener *EventListener) {
	eventBus.Lock()
	eventBus.subscribers[listener.Id] = listener
	eventBus.Unlock()
}

func UnSubscribe(listener *EventListener) {
	eventBus.Lock()
	delete(eventBus.subscribers, listener.Id)
	eventBus.Unlock()
}

func Publish(event Event, wg *sync.WaitGroup) {
	eventBus.RLock()

	for _, listener := range eventBus.subscribers {
		go func(listener *EventListener) {
			listener.Ch <- event
		}(listener)
	}
	eventBus.RUnlock()
	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	// 创建并注册监听者
	listener1 := EventListener{Id: 1, Ch: make(chan Event, 1), WaitHandle: &wg}
	Subscribe(&listener1)

	// 创建并注册另一个监听者
	listener2 := EventListener{Id: 2, Ch: make(chan Event, 1), WaitHandle: &wg}
	Subscribe(&listener2)

	// 发布事件
	wg.Add(3)
	go Publish(Event{Data: 10}, &wg)

	// 处理事件
	go func() {
		for e := range listener1.Ch {
			fmt.Printf("Listener1 received: %v\n", e.Data)
			listener1.WaitHandle.Done()
		}
	}()

	go func() {
		for e := range listener2.Ch {
			fmt.Printf("Listener2 received: %v\n", e.Data)
			listener2.WaitHandle.Done()
		}
	}()

	// 等待所有事件处理完毕
	wg.Wait()
	fmt.Println("All events processed")
}
