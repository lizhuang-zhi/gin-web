package main

import (
	"fmt"
	"sync"
	"time"
)

// 文章地址: https://xupin.im/2023/07/21/game-behavior-tree/

/*
	行为树是由一系列节点构成的，每个节点对应一个动作，自上而下执行且每个节点都会返回执行状态。

	行为树一般由：根节点（Root）、序列节点（Sequence）、选择节点（Selector）、并行节点（Parallel）等等节点组成，可能还会有条件节点（Condition）。这些节点也可以被抽象的归为 4 类:

	控制节点
	- 序列节点（Sequence）
	- 选择节点（Selector）
	- 并行节点（Parallel）

	条件节点
	- 条件节点（Condition）

	行为节点
	- 自定义节点。比如：攻击、防御、吃药等，是真正执行游戏逻辑的节点。

	装饰节点
	- 逆变节点（Inverter）
	- 重复执行节点（Repeater）
	- 装饰节点也算是自定义节点，主要作用是辅助子节点执行。

	* 需要注意的是，无论是什么节点都必须实现Exec()执行方法
*/

// 定义节点接口
type Status int8

const (
	Success Status = iota
	Failure
	Running // 新增 Running 状态
)

type INode interface {
	Exec(db IBlackboard) Status
}

// IBlackboard 用于传递数据
type IBlackboard interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}

type Blackboard struct {
	data map[string]interface{}
	lock sync.RWMutex
}

func (b *Blackboard) Get(key string) interface{} {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.data[key]
}

func (b *Blackboard) Set(key string, value interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.data[key] = value
}

// 序列节点 (Sequence): 有序的执行子节点，任意子节点失败终止执行返回失败，全部执行完毕返回成功
type Sequence struct {
	children []INode
}

func NewSequence(children ...INode) *Sequence {
	return &Sequence{children: children}
}

func (r *Sequence) Exec(db IBlackboard) Status {
	for _, v := range r.children {
		if v.Exec(db) == Failure {
			return Failure
		}
	}
	return Success
}

// 选择节点 (Selector): 同样有序的执行子节点、和序列节点不同的是任意子节点执行成功返回成功
type Selector struct {
	children []INode
}

func NewSelector(children ...INode) *Selector {
	return &Selector{children: children}
}

func (r *Selector) Exec(db IBlackboard) Status {
	for _, v := range r.children {
		if v.Exec(db) == Success {
			return Success
		}
	}
	return Failure
}

// 并行节点 (Parallel): 同时执行所有子节点
type Parallel struct {
	children []INode
}

func NewParallel(children ...INode) *Parallel {
	return &Parallel{children: children}
}

func (r *Parallel) Exec(db IBlackboard) Status {
	var wg sync.WaitGroup
	status := Success
	for _, v := range r.children {
		wg.Add(1)
		go func(v INode) {
			defer wg.Done()
			if v.Exec(db) == Failure {
				status = Failure
			}
		}(v)
	}
	wg.Wait()
	return status
}

// 条件节点 / 行为节点: 顾名思义，执行子节点前需要检查是不是符合条件
type Condition struct {
	f func(IBlackboard) bool
}

func NewCondition(f func(IBlackboard) bool) *Condition {
	return &Condition{f: f}
}

func (r *Condition) Exec(db IBlackboard) Status {
	if r.f(db) {
		return Success
	}
	return Failure
}

// 行为节点: 休息补充能量50, 蓄力需要能量大于等于100, 绝招需要蓄力, 释放绝招消耗能量100
// 绝招
type JueZhao struct{}

func NewJueZhao() *JueZhao {
	return &JueZhao{}
}

func (r *JueZhao) Exec(db IBlackboard) Status {
	db.Set("can_use_skill", false)
	db.Set("energy", db.Get("energy").(int)-100)
	fmt.Println("[绝招]done, 当前能量:", db.Get("energy"))
	return Success
}

// 蓄力
type XuLi struct{}

func NewXuLi() *XuLi {
	return &XuLi{}
}

func (r *XuLi) Exec(db IBlackboard) Status {
	db.Set("can_use_skill", true)
	fmt.Println("[蓄力]done")
	return Success
}

// 休息
type XiuXi struct{}

func NewXiuXi() *XiuXi {
	return &XiuXi{}
}

func (r *XiuXi) Exec(db IBlackboard) Status {
	energy, ok := db.Get("energy").(int)
	if !ok {
		energy = 0
	}
	fmt.Println("[休息]done, 补充能量50, 当前能量:", energy+50)
	db.Set("energy", energy+50)
	return Success
}

// 平A
type Attack struct{}

func NewAttack() *Attack {
	return &Attack{}
}

func (r *Attack) Exec(db IBlackboard) Status {
	db.Set("energy", db.Get("energy").(int)-20)
	fmt.Println("[平A]done, 当前能量:", db.Get("energy"))
	return Success
}

// 移动
type Move struct{}

func NewMove() *Move {
	return &Move{}
}

func (r *Move) Exec(db IBlackboard) Status {
	db.Set("energy", db.Get("energy").(int)-10)
	fmt.Println("[移动]done, 当前能量:", db.Get("energy"))
	return Success
}

/*
行为树的运行通常是以帧为单位进行的，每一帧都会跑一次行为树。为了支持这种运行方式，我们需要添加一个新的状态，叫做Running。
当一个节点返回Running状态时，它表示该节点正在执行一个耗时的操作，比如移动到一个地点或者播放一个动画。在这种情况下，
行为树不会从头开始执行，而是从正在运行的节点开始。这就需要我们在行为树中存储一些状态，记录当前正在运行的节点。
在Go语言中，我们可以用goroutine和channel来模拟这种行为。我们可以启动一个新的goroutine来执行耗时的操作，并在完成时向一个channel发送一个信号
*/
type WalkTo struct {
	Target string
	Done   chan bool
}

func NewWalkTo(target string) *WalkTo {
	return &WalkTo{
		Target: target,
		Done:   make(chan bool),
	}
}

func (r *WalkTo) Exec(db IBlackboard) Status {
	select {
	case <-r.Done:
		fmt.Printf("[到达%s]done\n", r.Target)
		return Success
	default:
		go func() {
			fmt.Printf("[正在往%s移动]...\n", r.Target)
			time.Sleep(2 * time.Second)
			r.Done <- true
		}()
		return Running
	}
}

func main() {
	// 定义根节点
	bt := NewSelector(
		NewSequence(
			NewCondition(func(db IBlackboard) bool { // 附近是否有敌人
				result := db.Get("is_enemy_near")
				if result == nil || !result.(bool) {
					return false
				}
				return true
			}),
			NewWalkTo("x12,y12"), // 远离敌人
		),
		NewXiuXi(), // 休息
	)

	db := &Blackboard{
		data: make(map[string]interface{}, 0),
	}

	go func() {
		time.Sleep(500 * time.Millisecond) // 模拟敌人出现
		db.Set("is_enemy_near", true)
		fmt.Println("****敌人出现*****")
	}()

	count := 5 // 行为树执行次数
	for i := 0; i < count; i++ {
		fmt.Printf("执行行为树%d次-----------------------\n", i+1)
		result := bt.Exec(db)
		fmt.Print("结果: ")
		switch result {
		case Success:
			fmt.Println("成功")
		case Failure:
			fmt.Println("失败")
		}

		time.Sleep(1 * time.Second)
	}
}
