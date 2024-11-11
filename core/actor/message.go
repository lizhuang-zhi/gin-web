// 包 actor 提供了一种异步消息传递机制，用于在不同的 actor 之间传递消息并执行相应的操作。
package actor

import (
	"context"
	"fmt"
	"sync"
)

// messagePool 是一个消息对象池，用于重用已分配的 Message 对象，提高性能并减少内存分配。
var messagePool = sync.Pool{
	New: func() any {
		return &Message{
			usePool: true,
		}
	},
}

// putMessagePool 将 Message 对象放回对象池中进行重用。
func putMessagePool(msg *Message) {
	if msg.usePool {
		messagePool.Put(msg)
	}
}

// Message 是一个消息对象，包含了消息的命令、数据、上下文以及一些元数据。
type Message struct {
	Command string          // 消息命令
	Body    interface{}     // 消息数据
	Context context.Context // 消息上下文
	usePool bool            // 是否使用对象池
	done    chan struct{}   // 用于通知消息处理完成
}

// Handler 是一个消息处理函数类型，用于处理接收到的消息。
type Handler func(*Message)

// NewMessage 构造一个新的消息对象，从对象池中获取一个 Message 对象，并设置相应的命令和数据。
func NewMessage(key, body interface{}) *Message {
	message := messagePool.Get().(*Message)
	message.Command = keyToCommand(key)
	message.Body = body
	message.Context = context.Background()
	return message
}

// NewMessageWithContext 构造一个新的消息对象，并设置指定的上下文。
func NewMessageWithContext(ctx context.Context, key, body interface{}) *Message {
	message := NewMessage(key, body)
	message.Context = ctx
	return message
}

// 以下是一些内部使用的命令常量。
const (
	commandRunner = "_Runner" // 运行命令
	commandStop   = "_Stop"   // 停止命令
)

// keyToCommand 将给定的键转换为命令字符串。
func keyToCommand(key interface{}) string {
	command, ok := key.(string)
	if !ok {
		command = fmt.Sprintf("%s", key)
	}

	return command
}
