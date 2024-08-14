package main

import (
	"fmt"
	"sync"
	"time"
)

// Event represents the abstract structure of an event.
type Event struct {
	Name string
	Data interface{}
}

// EventHandler is a function that handles an Event.
type EventHandler func(Event)

// EventListener is a struct that holds the ID and EventHandlers for each listener.
type EventListener struct {
	ID            int
	EventHandlers map[string][]EventHandler
}

// EventBus holds the subscribers to the events.
var EventBus = &EventBusStruct{
	subscribers: make(map[int]*EventListener),
}

type EventBusStruct struct {
	sync.RWMutex
	subscribers map[int]*EventListener
}

func (eb *EventBusStruct) Subscribe(listener EventListener) {
	eb.Lock()
	eb.subscribers[listener.ID] = &listener
	eb.Unlock()
}

func (eb *EventBusStruct) Unsubscribe(id int) {
	eb.Lock()
	delete(eb.subscribers, id)
	eb.Unlock()
}

func (eb *EventBusStruct) Publish(event Event) {
	eb.RLock()
	defer eb.RUnlock()

	for _, listener := range eb.subscribers {
		if handlers, ok := listener.EventHandlers[event.Name]; ok {
			for _, handler := range handlers {
				// Invoke the event handler in a separate goroutine.
				go handler(event)
			}
		}
	}
}

// Example event handlers
func NewBugHandler(event Event) {
	fmt.Printf("New Bug: %s\n", event.Data)
}

func NewTaskHandler(event Event) {
	fmt.Printf("New Task: %s\n", event.Data)
}

func main() {
	// Create and register a new listener
	listener1 := EventListener{
		ID: 1,
		EventHandlers: map[string][]EventHandler{
			"NewBug":  []EventHandler{NewBugHandler},
			"NewTask": []EventHandler{NewTaskHandler},
		},
	}

	EventBus.Subscribe(listener1)

	// Trigger some events
	EventBus.Publish(Event{Name: "NewBug", Data: "Bug1"})
	EventBus.Publish(Event{Name: "NewTask", Data: "Task1"})

	// Unsubscribe the listener
	EventBus.Unsubscribe(listener1.ID)

	time.Sleep(1 * time.Second)
}
