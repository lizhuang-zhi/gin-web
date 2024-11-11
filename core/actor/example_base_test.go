package actor

import (
	"fmt"

	"Solarland/Backend/core/slog"
	"Solarland/Backend/core/timer"
)

// Calculator 是一个实现简单计算器功能的结构体
type Calculator struct {
	*Base     // 嵌入 Base 结构体，以使用 Actor 模型
	num   int // num 表示计算器的当前值
}

// handler 是 Calculator 的消息处理函数
// 它根据接收到的消息执行相应的操作
func (c *Calculator) handler(message *Message) {
	switch message.Command {
	case "add": // 如果命令是 "add"，则将消息体（一个整数）加到 num 上
		c.num += message.Body.(int)
	case "dec": // 如果命令是 "dec"，则将消息体（一个整数）从 num 中减去
		c.num -= message.Body.(int)
	}
}

// ExampleBase 是一个使用 Base 结构体实现简单计算器的示例函数
// 它展示了如何创建一个 Actor，发送消息，执行操作，并停止 Actor
func ExampleBase() {
	calculator := &Calculator{} // 创建一个 Calculator 实例

	// 初始化 Base，并将 calculator.handler 设置为消息处理函数
	calculator.Base = NewBase(timer.NewHeapTimerManager(), slog.NewLogger(), "calculator", 10, calculator.handler, false, nil)

	calculator.Run() // 启动 Actor

	_ = calculator.Send(&Message{Command: "add", Body: 10}) // 发送 "add 10" 消息
	_ = calculator.Send(&Message{Command: "dec", Body: 5})  // 发送 "dec 5" 消息

	// 执行一个操作，打印 num 的值
	_ = calculator.Do(func() {
		fmt.Println(calculator.num)
	})

	calculator.Stop() // 停止 Actor
	// Output: 5
}
