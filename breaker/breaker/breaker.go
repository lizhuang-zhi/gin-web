package breaker

import (
	"context"
	"errors"
)

var (
	ErrMaxConcurrency = errors.New("max concurrency") // 达到最大并发
	ErrCircuitOpen    = errors.New("circuit open")    // 熔断开启
	ErrTimeout        = errors.New("timeout")         // 超时
)

type Breaker interface {
	IsOpen() bool // 是否开启了熔断
	Do(ctx context.Context, fn func(ctx context.Context) error, fallback func(context.Context, error) error) error
}

type SwitchWrapper struct {
	Breaker
	isOpen func() bool
}

/*
这种设计通常用于在某些特定场景下临时关闭断路器。比如，在进行某些敏感的操作时（如数据库迁移），可能需要暂时关闭断路器，
以防止由于偶发的错误导致断路器打开，进而影响到正常的操作。

我们并没有使用WrapSwitch，主要是因为我们的示例较简单，没有涉及到需要临时关闭断路器的场景。
*/
func WrapSwitch(breaker Breaker, isOpen func() bool) *SwitchWrapper {
	return &SwitchWrapper{
		Breaker: breaker,
		isOpen:  isOpen,
	}
}

func (sw *SwitchWrapper) Do(ctx context.Context, fn func(ctx context.Context) error, fallback func(context.Context, error) error) error {
	if sw.isOpen() {
		return sw.Breaker.Do(ctx, fn, fallback)
	}
	err := fn(ctx)
	if err != nil && fallback != nil {
		err = fallback(ctx, err)
	}
	return err
}
