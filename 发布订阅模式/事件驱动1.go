package main

import (
	"fmt"
)

type Event struct {
	name string
}

type EventHandler func(e *Event)

// 事件调度器: 包含一个事件列表
type EventDispatcher struct {
	handlers map[string][]EventHandler
}

func (d *EventDispatcher) Register(name string, handler EventHandler) {
	if _, ok := d.handlers[name]; !ok {
		d.handlers[name] = []EventHandler{}
	}
	d.handlers[name] = append(d.handlers[name], handler)
}

func (d *EventDispatcher) Trigger(name string, e *Event) {
	if _, ok := d.handlers[name]; ok {
		for _, h := range d.handlers[name] {
			h(e)
		}
	}
}

func main() {
	// 事件调度器
	dispatcher := &EventDispatcher{make(map[string][]EventHandler)}

	// 注册事件
	dispatcher.Register("say_hello", func(e *Event) {
		fmt.Println("Hello, world!")
	})

	// 触发事件
	dispatcher.Trigger("say_hello", &Event{name: "say_hello"})
}
