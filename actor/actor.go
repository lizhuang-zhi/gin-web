package main

import (
	"fmt"
	"time"
)

type Message struct {
	Data string
}

type Actor struct {
	Inbox chan Message
}

func NewActor() *Actor {
	return &Actor{
		Inbox: make(chan Message, 10),
	}
}

func (a *Actor) Receive(message Message) {
	fmt.Println("Received: ", message.Data)
}

func main() {
	actor := NewActor()

	go func() {
		for {
			message := <-actor.Inbox
			actor.Receive(message)
		}
	}()

	for i := 0; i < 10; i++ {
		actor.Inbox <- Message{Data: fmt.Sprintf("message %d", i)}
	}

	time.Sleep(time.Second)
}
