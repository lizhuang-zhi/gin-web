package actor

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"Solarland/Backend/core/exception"
	"Solarland/Backend/core/slog"
	"Solarland/Backend/core/timer"
	"Solarland/Backend/core/utils/backoff"
	"Solarland/Backend/core/utils/option"
)

var (
	ErrOverflow    = errors.New("overflow")
	ErrActorClosed = errors.New("actor closed")
)

// Base 基础Actor结构体
type Base struct {
	slog.Logger

	FlushDu time.Duration // 缓存消息刷新时间间隔

	closed   uint32          // 是否已关闭
	name     string          // Actor名称
	msgChan  chan *Message   // 消息通道
	handler  Handler         // 消息处理函数
	tickers  *sync.Map       // 定时器集合
	metadata option.MetaData // 元数据
	tm       timer.Group     // 定时器管理器组
	timerid  int32           // 定时器ID
	wg       sync.WaitGroup  // 等待组

	mtx       sync.RWMutex         // 读写锁
	overflows *list.List           // 消息缓存队列
	recounter *exception.Recounter // 异常计数器
}

// NewBase 创建一个新的Base实例
func NewBase(tm timer.Manager, logger slog.Logger, name string, size int, handler Handler, noOverflow bool, recounterCfg func() bool) *Base {
	b := &Base{}
	b.Logger = logger
	b.name = name
	b.tm = tm.Group(name)
	if noOverflow {
		b.overflows = list.New()
	}
	b.FlushDu = time.Millisecond * 100
	b.msgChan = make(chan *Message, size)
	b.handler = handler
	b.tickers = new(sync.Map)
	b.recounter = exception.NewRecounter(exception.NewRecounterConf().SetOpen(recounterCfg))
	b.metadata = option.NewMetaData("Actor")
	b.metadata.Refresh(func() {
		var l int
		b.mtx.RLock()
		if b.overflows != nil {
			l = b.overflows.Len()
		}
		b.mtx.RUnlock()
		b.metadata.Set(MetaPending, len(b.msgChan))
		b.metadata.Set(MetaPendingCache, l)
		b.metadata.Set(MetaPendingSource, b.recounter.TopString(10))
	})
	return b
}

// Name 返回Actor名称
func (b *Base) Name() string {
	return b.name
}

// 发送消息
func (b *Base) Send(message *Message) error {
	if atomic.LoadUint32(&b.closed) == 1 {
		putMessagePool(message)
		return ErrActorClosed
	}

	select {
	case b.msgChan <- message:
	default:
		if !b.doCache(message) {
			putMessagePool(message)
			return ErrOverflow
		}
	}
	// 更新元数据
	b.metadata.AddInt(MetaSend, 1)
	return nil
}

// 同步发送消息
func (b *Base) SyncSend(message *Message, timeout time.Duration) error {
	message.done = make(chan struct{}, 1)

	err := b.Send(message)
	if err != nil {
		return err
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-message.done:
	case <-timer.C:
		return errors.New("timeout")
	}

	return nil
}

// 运行
func (b *Base) Run() {
	b.wg.Add(1)
	exception.GO(func() {
		defer b.wg.Done()
		tk := time.NewTicker(b.FlushDu)
		defer tk.Stop()
		for {
			select {
			case msg, ok := <-b.msgChan:
				if !ok {
					return
				}
				b.doMsg(msg)
			case <-tk.C:
				b.mtx.Lock()
				of := b.overflows
				if of != nil {
					b.overflows = list.New()
				}
				b.mtx.Unlock()

				if of != nil && of.Len() > 0 {
					b.Infof("[Actor] flush overflow %s %d", b.Name(), of.Len())
					for e := of.Front(); e != nil; e = e.Next() {
						b.doMsg(e.Value.(*Message))
					}
				}
			}
		}
	}, exception.WithOpen(false))
}

func (b *Base) doMsg(msg *Message) {
	b.invoke(msg)
	putMessagePool(msg)
}

func (b *Base) doCache(msg *Message) bool {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.overflows == nil {
		return false
	}
	b.overflows.PushBack(msg)
	return true
}

func (b *Base) invoke(msg *Message) {
	defer exception.Recovery()

	var st time.Time
	var cv any
	if msg.Context != nil {
		cv = msg.Context.Value(CostWarnContext{})
		if cv != nil {
			st = time.Now()
		}
	}
	defer func() {
		if cv != nil {
			cwc := cv.(*CostWarnContext)
			if cost := time.Since(st); cost > cwc.CostLimit {
				b.Warnf("[ACTOR.BASE] invoke %s cost: %.3fs", cwc.Key, cost.Seconds())
			}
		}
	}()

	switch msg.Command {
	case commandRunner:
		runner := msg.Body.(func())
		runner()
		return
	case commandStop:
		close(b.msgChan)
		return
	}

	if b.handler != nil {
		b.handler(msg)
	}

	if msg.done != nil {
		msg.done <- struct{}{}
	}
}

// 停止
func (b *Base) Stop() {
	if !atomic.CompareAndSwapUint32(&b.closed, 0, 1) {
		return
	}

	// 停止定时器
	b.tickers.Range(func(key, _ interface{}) bool {
		b.tm.StopTimer(b.tickerName(key.(string)))
		return true
	})
	b.tickers = new(sync.Map)
	b.msgChan <- NewMessage(commandStop, nil)
	b.wg.Wait()
}

func (b *Base) recount(opt ...exception.RecounterOption) func() {
	if len(opt) == 0 {
		opt = []exception.RecounterOption{exception.WithCallerSkipOffset(1)}
	}
	return exception.Recount(b.recounter, opt...)
}

// Do 等待执行函数返回
func (b *Base) Do(runner func()) error {
	return b.DoCtx(context.Background(), runner)
}

// DoCtx 等待执行函数返回
func (b *Base) DoCtx(ctx context.Context, runner func()) error {
	wait := make(chan struct{})
	cancel := b.recount()
	err := b.Send(NewMessageWithContext(ctx, commandRunner, func() {
		defer func() {
			wait <- struct{}{}
			cancel()
		}()

		runner()
	}))
	if err != nil {
		cancel()
		return err
	}

	<-wait

	return nil
}

// DoAsync 异步执行函数
func (b *Base) DoAsync(runner func(), opt ...exception.RecounterOption) error {
	cancel := b.recount(opt...)
	if err := b.Send(NewMessage(commandRunner, func() {
		defer cancel()

		runner()
	})); err != nil {
		cancel()
		return err
	}
	return nil
}

// Call 调用函数并返回结果
func (b *Base) Call(runner func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error

	doErr := b.Do(func() {
		result, err = runner()
	})
	if doErr != nil {
		return nil, doErr
	}

	return result, err
}

// TickerAction 定时动作结构体
type TickerAction struct {
	Key    string          // 键
	D      time.Duration   // 时间间隔
	Runner func()          // 执行函数
	B      backoff.Backoff // 退避策略

	mtx              sync.Mutex
	c                int
	du               time.Duration
	lstBegin, lstEnd time.Time
}

func (t *TickerAction) try(logger slog.Logger, f func()) func() {
	return func() {
		if t.begin() {
			defer func() { t.end(logger) }()
			f()
		}
	}
}

func (t *TickerAction) begin() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if !t.lstBegin.IsZero() {
		return false
	}

	if time.Since(t.lstEnd) >= t.du {
		t.lstBegin = time.Now()
		t.lstEnd = time.Time{}
		return true
	}
	return false
}

func (t *TickerAction) end(logger slog.Logger) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	cost := time.Since(t.lstBegin)
	t.lstBegin = time.Time{}
	t.lstEnd = time.Now()

	if t.B != nil {
		if cost <= t.du {
			if t.c <= 0 {
				return
			}
			t.c--
		} else {
			logger.Warnf("[Actor] TickerObj %s(c: %v) cost: %.3f/%.3f", t.Key, t.c, cost.Seconds(), t.du.Seconds())
			t.c += 3
		}
		t.du = t.B.Next(t.c)
	}
}

// TickerObj 添加一个定时动作
func (b *Base) TickerObj(obj *TickerAction) {
	nkey := b.tickerName(obj.Key)
	obj.du = obj.D
	if err := b.tm.Scheduler(nkey,
		obj.try(b.Logger, func() {
			b.runTicker(nkey, obj.Key, obj.Runner)
		}), time.Now(), obj.D, -1,
		timer.WithPayloadFileLine("actor.tickerobj.%s"+obj.Key, 0),
		timer.WithNoOverlap(),
	); err != nil {
		b.Errorf("Ticker %s error: %s", nkey, err)
	} else {
		b.tickers.Store(obj.Key, struct{}{})
	}
}

func (b *Base) tickerName(key string) string {
	return b.Name() + "." + key
}

// Ticker 添加一个定时任务
func (b *Base) Ticker(key string, d time.Duration, runner func()) {
	nkey := b.tickerName(key)
	if err := b.tm.Scheduler(nkey,
		func() {
			b.runTicker(nkey, key, runner)
		},
		time.Now(), d, -1,
		timer.WithPayloadFileLine("actor.tickerobj.%s"+key, 0),
		timer.WithNoOverlap(),
	); err != nil {
		b.Errorf("Ticker %s error: %s", nkey, err)
	} else {
		b.tickers.Store(key, struct{}{})
	}
}

// CostWarnContext 耗时警告上下文
type CostWarnContext struct {
	CostLimit time.Duration
	Key       string
}

func (b *Base) runTicker(nkey, key string, runner func()) {
	ctx := context.WithValue(context.Background(), CostWarnContext{}, &CostWarnContext{CostLimit: time.Millisecond * 50, Key: key})
	err := b.DoCtx(ctx, runner)
	if err != nil {
		if err == ErrActorClosed {
			b.tickers.Delete(key)
			b.tm.StopTimer(nkey)
			return
		}
		b.Errorf("ActorBase Ticker with err: %v", err)
	}
}

// StopTicker 停止一个定时任务
func (b *Base) StopTicker(key string) {
	b.tickers.Delete(key)
	b.tm.StopTimer(b.tickerName(key))
}

// Timer 添加一个延时任务
func (b *Base) Timer(d time.Duration, runner func(), opt ...timer.Option) string {
	if len(opt) == 0 {
		opt = []timer.Option{timer.WithCurrentPayloadFileLine(1)}
	}

	key := fmt.Sprintf("once.%d", atomic.AddInt32(&b.timerid, 1))
	nkey := b.tickerName(key)
	if err := b.tm.Scheduler(nkey, func() {
		b.tickers.Delete(key)
		err := b.Do(runner)
		if err != nil {
			b.Errorf("ActorBase Timer with err: %v", err)
		}
	}, time.Now(), d, 1, opt...); err != nil {
		b.Errorf("Ticker %s error: %s", nkey, err)
		return ""
	} else {
		b.tickers.Store(key, struct{}{})
	}
	return key
}

// MetaData 返回元数据
func (b *Base) MetaData() option.MetaData {
	return b.metadata
}
