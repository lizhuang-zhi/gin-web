package sync

import (
	"sync"
	"sync/atomic"
)

/*
Leader选举（Leader Election）是分布式系统中的一个基本过程，用于在多个进程（或服务器、服务实例等）中选出一个领导者（Leader），这个领导者通常负责协调工作、管理资源或者进行决策，确保系统的一致性和高可用性。以下是Leader选举过程的一些关键点，以及它的重要性：

协调和决策：在分布式系统中，可能有很多个副本或者节点同时工作。为了避免冲突和确保数据的一致性，通常需要一个领导者来协调这些节点的工作和做出决策。

容错性：如果一个领导者因为某些原因失效了，系统需要能够迅速选出一个新的领导者，从而继续正常的工作流程，这就需要领导者选举机制。

负载分配：在某些情况下，领导者可以负责分配任务给其他节点，以负载均衡的方式提高系统的效率。

避免脑裂（Split-Brain）：在没有领导者的情况下，系统的不同部分可能会试图独立操作，产生冲突的决策，这种情况被称为脑裂。领导者选举机制可以帮助系统避免这种问题。

在分布式系统中，leader选举可以通过多种算法来实现，比如Raft或者Paxos。这些算法通常保证即使在网络分区、节点失效等异常情况下，系统也能够选举出一个领导者，并且整个系统能够继续正常运作。

在实际使用中，leader选举的实现通常会依赖于像Consul、Etcd、Zookeeper等分布式一致性存储服务。这些服务提供了选举的机制，确保在分布式环境中，选举过程是可靠和一致的。

总结来说，Leader选举是在一组节点中选出一个主导者，这个主导者对系统的正常运行至关重要，它确保了系统在面对故障时能够持续运行，并且对外提供一致性的服务。
*/

// Leader 接口代表领导者选举的功能。
type Leader interface {
	Elect(key string, val interface{}) (Elected, error)
}

// Elected 接口代表一个被选举的实例。
type Elected interface {
	Resign() error            // 放弃领导者的地位
	Revoked() <-chan struct{} // 返回一个通道，当领导者被撤销时，通道会被关闭
}

// LeaderAction 封装了领导者选举的逻辑。
type LeaderAction struct {
	key       string
	val       interface{}
	electFn   func()        // 当选举成功时执行的函数
	revokeFn  func()        // 当领导者被撤销时执行的函数
	elected   Elected       // 被选举的实例
	isElected int32         // 是否被选举为领导者的标志
	stopChan  chan struct{} // 停止信号的通道
	wg        sync.WaitGroup
}

func NewLeaderAction(key string, val interface{}) *LeaderAction {
	// 初始化LeaderAction结构体
	return &LeaderAction{
		key:      key,
		val:      val,
		stopChan: make(chan struct{}),
	}
}

func (la *LeaderAction) Elect(leader Leader) error {
	// 试图成为领导者
	elected, err := leader.Elect(la.key, la.val)
	if err != nil {
		return err
	}
	la.elected = elected
	atomic.StoreInt32(&la.isElected, 1) // 标记为已被选举
	if la.electFn != nil {
		la.electFn()
	}
	return nil
}

func (la *LeaderAction) Resign() error {
	// 放弃领导者的地位
	if la.elected != nil {
		return la.elected.Resign()
	}
	return nil
}

func (la *LeaderAction) Watch() {
	// 监控领导者的地位是否被撤销
	la.wg.Add(1)
	go func() {
		defer la.wg.Done()
		for {
			select {
			case <-la.elected.Revoked():
				atomic.StoreInt32(&la.isElected, 0) // 标记为未被选举
				if la.revokeFn != nil {
					la.revokeFn()
				}
				return
			case <-la.stopChan:
				return
			}
		}
	}()
}

func (la *LeaderAction) Stop() {
	// 停止监控
	close(la.stopChan)
	la.wg.Wait()
}

func (la *LeaderAction) IsElected() bool {
	// 返回是否被选举为领导者
	return atomic.LoadInt32(&la.isElected) == 1
}

func (la *LeaderAction) SetElectFn(fn func()) {
	// 设置当选举成功时执行的函数
	la.electFn = fn
}

func (la *LeaderAction) SetRevokeFn(fn func()) {
	// 设置当领导者被撤销时执行的函数
	la.revokeFn = fn
}
