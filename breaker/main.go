package main

import (
	"booking-app/breaker/hystrix"
	"context"
	"errors"
	"fmt"
)

func main() {
	config := &hystrix.Config{
		Name: "test",
		CommandConfig: hystrix.CommandConfig{
			Timeout:               1000,
			MaxConcurrentRequests: 100,
			ErrorPercentThreshold: 50,
		},
	}

	myBreaker := hystrix.NewBreaker(config)

	myTask := func(ctx context.Context) error {
		return errors.New("task failed")
	}

	myFallback := func(ctx context.Context, err error) error {
		fmt.Println("Doing fallback")
		return err
	}

	err := myBreaker.Do(context.Background(), myTask, myFallback)
	fmt.Println(err)
}
