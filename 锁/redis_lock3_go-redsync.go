package main

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

func main() {
	pools := []redsync.Pool{
		&redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", "localhost:6379")
			},
		},
	}

	rs := redsync.New(pools)

	mutex := rs.NewMutex("test-redsync")

	err := mutex.Lock()
	if err != nil {
		panic(err)
	}

	fmt.Println("We got the lock!")

	time.Sleep(3 * time.Second)

	fmt.Println("Releasing the lock...")

	_, err = mutex.Unlock()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully released the lock")
}
