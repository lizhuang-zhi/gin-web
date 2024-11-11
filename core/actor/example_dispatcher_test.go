package actor

import (
	"fmt"

	"Solarland/Backend/core/slog"
	"Solarland/Backend/core/timer"
)

// Player 结构体表示一个玩家
type Player struct {
	*Dispatcher        // 嵌入 Dispatcher 结构体
	account     string // 玩家账号
	level       int    // 玩家等级
}

// handleLogin 处理玩家登录消息
func (p *Player) handleLogin(message *Message) {
	p.account = message.Body.(string) // 从消息体中获取账号信息
}

// handleLevelup 处理玩家升级消息
func (p *Player) handleLevelup(message *Message) {
	p.level++ // 玩家等级加1
}

// ExampleDispatcher 演示如何使用 Dispatcher 实现一个简单的玩家消息分发机制
func ExampleDispatcher() {
	player := &Player{} // 创建一个新的玩家实例
	// 初始化 Dispatcher
	player.Dispatcher = NewDispatcher(timer.NewHeapTimerManager(), slog.NewLogger(), "player", 10)
	// 注册消息处理函数
	player.Attach("login", player.handleLogin)
	player.Attach("levelup", player.handleLevelup)
	player.Run() // 启动 Dispatcher

	// 发送登录消息
	_ = player.Send(&Message{Command: "login", Body: "test"})
	// 发送升级消息
	_ = player.Send(&Message{Command: "levelup"})

	// 通过调用 Call 方法获取玩家等级
	level, _ := player.Call(func() (interface{}, error) {
		fmt.Println(player.account) // 输出玩家账号
		return player.level, nil    // 返回玩家等级
	})

	fmt.Println(level) // 输出玩家等级

	player.Stop() // 停止 Dispatcher
	// Output:
	// test
	// 1
}
