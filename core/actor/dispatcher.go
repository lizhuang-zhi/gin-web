package actor

import (
	"Solarland/Backend/core/slog"
	"Solarland/Backend/core/timer"
)

// Dispatcher 是一个分发器，用于将消息分发给相应的处理程序
type Dispatcher struct {
	*Base
	m map[interface{}]Handler
}

// NewDispatcher 创建一个新的分发器实例
// tm 是一个计时器管理器
// logger 是一个日志记录器
// name 是分发器的名称
// size 是分发器的队列大小
func NewDispatcher(tm timer.Manager, logger slog.Logger, name string, size int) *Dispatcher {
	d := &Dispatcher{}
	d.Base = NewBase(tm, logger, name, size, d.handle, false, nil)
	d.m = make(map[interface{}]Handler)

	return d
}

// Attach 将一个处理程序附加到分发器上，关联一个键
// key 是一个唯一的键，用于标识处理程序
// handle 是要附加的处理程序
func (d *Dispatcher) Attach(key interface{}, handle Handler) {
	d.m[keyToCommand(key)] = handle
}

// Detach 从分发器上分离一个处理程序
// key 是要分离的处理程序的键
func (d *Dispatcher) Detach(key interface{}) {
	delete(d.m, keyToCommand(key))
}

// handle 是分发器的内部处理函数，用于将消息分发给相应的处理程序
func (d *Dispatcher) handle(message *Message) {
	handler := d.m[message.Command]
	if handler != nil {
		handler(message)
	}
}
