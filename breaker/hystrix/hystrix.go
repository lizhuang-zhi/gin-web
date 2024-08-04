// hystrix.go
package hystrix

import (
	"booking-app/breaker/breaker"
	"context"
	"errors"
)

type Config struct {
	Name          string
	CommandConfig CommandConfig
}

type CommandConfig struct {
	Timeout               int
	MaxConcurrentRequests int
	ErrorPercentThreshold int
}

type hystrixBreaker struct {
	config *Config
}

func NewBreaker(config *Config) breaker.Breaker {
	return &hystrixBreaker{config: config}
}

func (b *hystrixBreaker) IsOpen() bool {
	// 这是一个模拟的实现，通常你需要访问真实的 hystrix 状态
	return false
}

func (b *hystrixBreaker) Do(ctx context.Context, fn func(ctx context.Context) error, fallback func(context.Context, error) error) error {
	err := fn(ctx)
	switch {
	case errors.Is(err, breaker.ErrMaxConcurrency):
		return breaker.ErrMaxConcurrency
	case errors.Is(err, breaker.ErrCircuitOpen):
		return breaker.ErrCircuitOpen
	case errors.Is(err, breaker.ErrTimeout):
		return breaker.ErrTimeout
	}

	if err != nil && fallback != nil {
		err = fallback(ctx, err)
	}
	return err
}
