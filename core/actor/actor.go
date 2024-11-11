// Actor模式
package actor

import (
	"time"

	"Solarland/Backend/core/utils/option"
)

// Actor 接口定义了一个Actor的基本行为
type Actor interface {
	// Run 启动Actor
	Run()
	// Stop 停止Actor
	Stop()
	// Name 返回Actor的名称
	Name() string
	// MetaData 返回Actor的元数据
	MetaData() option.MetaData

	// Send 异步发送消息到Actor的消息队列
	Send(*Message) error
	// SyncSend 同步发送消息到Actor的消息队列，并等待指定时间直到消息被处理或超时
	SyncSend(*Message, time.Duration) error
	// Do 在Actor的消息处理循环中执行指定的函数
	Do(runner func()) error
}

const (
	// StatActors 用于统计Actor的指标名称
	StatActors = "Actors"
)

// 元数据
const (
	// MetaSend 用于统计Actor发送的消息数量
	MetaSend = "Send"
	// MetaPending 用于统计Actor等待处理的消息数量
	MetaPending = "Pending"
	// MetaPendingCache 用于统计Actor等待处理的消息缓存数量
	MetaPendingCache = "PendingCache"
	// MetaPendingSource 用于统计Actor等待处理的消息来源
	MetaPendingSource = "PendingSource"
)
