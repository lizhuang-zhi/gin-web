package actor

import (
	"testing"
	"time"

	"Solarland/Backend/core/slog"
	"Solarland/Backend/core/timer"

	. "github.com/smartystreets/goconvey/convey"
)

// TestSend 测试向 actor 发送消息的功能
func TestSend(t *testing.T) {
	Convey("Send", t, func() {
		// 创建一个通道用于接收 actor 处理后的消息体
		count := make(chan int)

		// 创建一个新的 actor，设置消息处理函数为将消息体写入 count 通道
		actor := NewBase(timer.NewHeapTimerManager(), slog.NewLogger(), "test", 10, func(message *Message) {
			count <- message.Body.(int)
		}, false, nil)
		// 启动 actor
		actor.Run()

		// 向 actor 发送消息
		err := actor.Send(&Message{Command: "test", Body: 1})

		// 验证消息体是否被正确处理
		So(<-count, ShouldEqual, 1)
		// 验证发送消息是否成功
		So(err, ShouldBeNil)
		// 停止 actor
		actor.Stop()
	})
}

// TestSyncSend 测试同步向 actor 发送消息的功能
func TestSyncSend(t *testing.T) {
	Convey("SyncSend", t, func() {
		// 定义一个计数器用于记录消息处理次数
		count := 0
		// 创建一个新的 actor，设置消息处理函数为增加计数器
		actor := NewBase(timer.NewHeapTimerManager(), slog.NewLogger(), "test", 10, func(message *Message) {
			count++
		}, false, nil)
		// 启动 actor
		actor.Run()

		// 同步向 actor 发送消息，设置超时时间为 1 秒
		err := actor.SyncSend(&Message{Command: "test", Body: 1}, time.Second)

		// 验证消息是否被正确处理
		So(count, ShouldEqual, 1)
		// 验证发送消息是否成功
		So(err, ShouldBeNil)
		// 停止 actor
		actor.Stop()
	})
}

// TestDo 测试在 actor 内部执行函数的功能
func TestDo(t *testing.T) {
	Convey("Do", t, func() {
		// 创建一个新的 actor
		actor := NewBase(timer.NewHeapTimerManager(), slog.NewLogger(), "test", 10, nil, false, nil)
		// 启动 actor
		actor.Run()

		Convey("Normal", func() {
			// 定义一个标志变量用于检查函数是否被执行
			var runFlag bool

			// 在 actor 内部执行函数，将 runFlag 设置为 true
			err := actor.Do(func() {
				runFlag = true
			})

			// 验证函数是否被正确执行
			So(runFlag, ShouldBeTrue)
			// 验证执行函数是否成功
			So(err, ShouldBeNil)
		})

		Convey("Panic", func() {
			// 定义一个标志变量用于检查函数是否被执行
			var runFlag bool

			// 在 actor 内部执行一个会 panic 的函数
			err := actor.Do(func() {
				panic("error")
			})

			// 验证执行函数是否成功（即使发生 panic 也应该返回 nil）
			So(err, ShouldBeNil)

			// 在 actor 内部执行另一个函数，将 runFlag 设置为 true
			err = actor.Do(func() {
				runFlag = true
			})

			// 验证函数是否被正确执行
			So(runFlag, ShouldBeTrue)
			// 验证执行函数是否成功
			So(err, ShouldBeNil)
		})
		Convey("stop", func() {
			// 定义一个标志变量用于检查函数是否被执行
			var runFlag bool

			// 在 actor 内部执行一个会阻塞 100 毫秒的函数
			err := actor.Do(func() {
				time.Sleep(100 * time.Millisecond)
			})
			// 验证执行函数是否成功
			So(err, ShouldBeNil)
			// 在 actor 内部执行另一个函数，将 runFlag 设置为 true
			err = actor.Do(func() {
				runFlag = true
			})
			// 停止 actor
			actor.Stop()
			// 验证函数是否被正确执行
			So(runFlag, ShouldBeTrue)
			// 验证执行函数是否成功
			So(err, ShouldBeNil)
		})
	})
}

// TestCall 测试在 actor 内部调用函数并获取返回值的功能
func TestCall(t *testing.T) {
	Convey("Call", t, func() {
		// 创建一个新的 actor
		actor := NewBase(timer.NewHeapTimerManager(), slog.NewLogger(), "test", 10, nil, false, nil)
		// 启动 actor
		actor.Run()
		// 在 actor 内部调用一个返回字符串 "test" 的函数
		v, err := actor.Call(func() (interface{}, error) {
			return "test", nil
		})

		// 验证返回值是否正确
		So(v, ShouldEqual, "test")
		// 验证调用函数是否成功
		So(err, ShouldBeNil)
		// 停止 actor
		actor.Stop()
	})
}
