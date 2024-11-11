// Backend/core/actor/manager.go

package actor

import (
	"sync"

	"Solarland/Backend/core/utils/option"
)

// Manager 是一个用于管理 Actor 的结构体
type Manager struct {
	actors   *sync.Map       // 存储已注册的 Actor 实例
	metadata option.MetaData // 元数据，用于存储管理器的相关信息
}

// NewManager 创建并返回一个新的 Manager 实例
func NewManager() *Manager {
	m := &Manager{
		actors:   new(sync.Map),                      // 初始化 sync.Map 用于存储 Actor
		metadata: option.NewMetaData("ActorManager"), // 初始化元数据，名称为 "ActorManager"
	}
	return m
}

// GetActor 根据名称获取已注册的 Actor 实例
func (mgr *Manager) GetActor(name string) Actor {
	value, ok := mgr.actors.Load(name) // 从 sync.Map 中加载指定名称的 Actor
	if !ok {
		return nil // 如果不存在，则返回 nil
	}

	return value.(Actor) // 类型断言并返回 Actor 实例
}

// Register 注册一个新的 Actor 实例
func (mgr *Manager) Register(a Actor) bool {
	_, loaded := mgr.actors.LoadOrStore(a.Name(), a) // 尝试将 Actor 存储到 sync.Map 中
	if loaded {
		return false // 如果已存在同名的 Actor，则返回 false
	}

	a.Run() // 运行新注册的 Actor
	// 统计
	mgr.metadata.AddInt(StatActors, 1) // 增加 StatActors 计数器
	// 元数据
	if mgr.metadata.Len() <= 100 { // 如果元数据长度不超过 100
		mgr.metadata.Set(a.Name(), a.MetaData()) // 将 Actor 的元数据存储到管理器的元数据中
	}
	return true
}

// UnRegister 取消注册指定名称的 Actor
func (mgr *Manager) UnRegister(name string) {
	a, ok := mgr.actors.Load(name) // 从 sync.Map 中加载指定名称的 Actor
	if !ok {
		return // 如果不存在，则直接返回
	}

	a.(Actor).Stop() // 停止 Actor

	mgr.actors.Delete(name) // 从 sync.Map 中删除指定名称的 Actor
	// 统计
	mgr.metadata.AddInt(StatActors, -1) // 减少 StatActors 计数器
	// 元数据
	mgr.metadata.Del(name) // 从管理器的元数据中删除指定名称的键值对
}

// Stop 停止所有已注册的 Actor
func (mgr *Manager) Stop() {
	mgr.actors.Range(func(name, m interface{}) bool { // 遍历 sync.Map 中的所有 Actor
		m.(Actor).Stop() // 停止当前 Actor
		return true      // 继续遍历
	})
}

// MetaData 返回管理器的元数据
func (mgr *Manager) MetaData() option.MetaData {
	return mgr.metadata
}
