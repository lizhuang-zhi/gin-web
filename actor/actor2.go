package main

import (
	"fmt"
	"time"
)

type Message struct {
	Data string
}

type Actor struct {
	Inbox      chan Message
	Collleague *Actor
}

func NewActor() *Actor {
	return &Actor{
		Inbox: make(chan Message, 10),
	}
}

func (a *Actor) SetCollleague(c *Actor) {
	a.Collleague = c
}

func (a *Actor) Receive(message Message) {
	fmt.Println("Received: ", message.Data)

	// 在接收到消息后，该 Actor 将消息转发给它的 Collleague
	if a.Collleague != nil {
		a.Collleague.Inbox <- Message{Data: "Hello from actor! from collleague actor"}
	}
}

func main() {
	actor1 := NewActor()
	actor2 := NewActor()

	actor1.SetCollleague(actor2)
	actor2.SetCollleague(actor1)

	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			message := <-actor1.Inbox
			actor1.Receive(message)
		}
	}()

	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			message := <-actor2.Inbox
			actor2.Receive(message)
		}
	}()

	// 主线程向 actor1 发送消息
	actor1.Inbox <- Message{Data: "Hello, actor! from main"}

	// 等待一秒钟，以便所有的 actor 都有足够的处理和转发消息的时间
	time.Sleep(3 * time.Millisecond)
}
