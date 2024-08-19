package main

import (
	"fmt"
	"time"
)

type Event struct {
	Data string
}

type EventChannel chan Event

type EventHandler func(Event)

type EventDispatcher struct {
	Handlers map[string][]EventHandler
}

func (d *EventDispatcher) Register(eventName string, handler EventHandler) {
	if _, found := d.Handlers[eventName]; !found {
		d.Handlers[eventName] = []EventHandler{}
	}
	d.Handlers[eventName] = append(d.Handlers[eventName], handler)
}

func (d *EventDispatcher) Dispatch(eventName string) EventChannel {
	eventChannel := make(EventChannel)
	if handlers, found := d.Handlers[eventName]; found {
		for _, handler := range handlers {
			go func(handler EventHandler) {
				for e := range eventChannel {
					handler(e)
				}
			}(handler)
		}
	}
	return eventChannel
}

func main() {
	d := EventDispatcher{
		Handlers: map[string][]EventHandler{},
	}

	d.Register("event1", func(e Event) {
		fmt.Println("event1 handler 1:", e.Data)
	})

	d.Register("event1", func(e Event) {
		time.Sleep(1 * time.Second)
		fmt.Println("event1 handler 2:", e.Data)
	})

	d.Register("event2", func(e Event) {
		fmt.Println("event2 handler 1:", e.Data)
	})

	eventChannel1 := d.Dispatch("event1")
	eventChannel1 <- Event{Data: "event1 data"}
	eventChannel1 <- Event{Data: "event1 data222"}

	eventChannel2 := d.Dispatch("event2")
	eventChannel2 <- Event{Data: "event2 data"}

	time.Sleep(2 * time.Second) // Wait for goroutines to finish
}
