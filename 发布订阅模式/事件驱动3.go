package main

import (
	"fmt"
	"sync"
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

func (d *EventDispatcher) Dispatch(eventName string, wg *sync.WaitGroup) EventChannel {
	eventChannel := make(EventChannel)
	if handlers, found := d.Handlers[eventName]; found {
		wg.Add(len(handlers)) // Add the number of handlers to the WaitGroup
		for _, handler := range handlers {
			go func(handler EventHandler) {
				for e := range eventChannel {
					handler(e)
				}
				wg.Done() // Decrement the counter when the goroutine completes
			}(handler)
		}
	}
	return eventChannel
}

func main() {
	var wg sync.WaitGroup // Create a new WaitGroup

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

	eventChannel1 := d.Dispatch("event1", &wg)
	eventChannel1 <- Event{Data: "event1 data"}
	eventChannel1 <- Event{Data: "event1 data222"}
	close(eventChannel1)

	eventChannel2 := d.Dispatch("event2", &wg)
	eventChannel2 <- Event{Data: "event2 data"}
	close(eventChannel2)

	wg.Wait() // Wait for all goroutines to finish
}
