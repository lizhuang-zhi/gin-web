package sync

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

// import (
// 	"log"
// 	"testing"
// )

// // MockLeader 模拟 Leader 实现。
// type MockLeader struct {
// 	elected map[string]Elected
// }

// func NewMockLeader() *MockLeader {
// 	return &MockLeader{elected: make(map[string]Elected)}
// }

// func (m *MockLeader) Elect(key string, val interface{}) (Elected, error) {
// 	elected := &MockElected{
// 		revoked: make(chan struct{}),
// 	}
// 	m.elected[key] = elected
// 	log.Printf("模拟选举: %s 成为领导者", key)
// 	return elected, nil
// }

// type MockElected struct {
// 	revoked chan struct{}
// }

// func (m *MockElected) Resign() error {
// 	close(m.revoked)
// 	log.Printf("模拟领导者放弃地位")
// 	return nil
// }

// func (m *MockElected) Revoked() <-chan struct{} {
// 	return m.revoked
// }

// func TestLeaderElection(t *testing.T) {
// 	mockLeader := NewMockLeader()
// 	leaderAction := NewLeaderAction("test_key", "test_value")

// 	electCh := make(chan struct{})
// 	revokeCh := make(chan struct{})

// 	leaderAction.SetElectFn(func() {
// 		log.Println("获得领导者地位的回调被调用")
// 		close(electCh)
// 	})

// 	leaderAction.SetRevokeFn(func() {
// 		log.Println("放弃领导者地位的回调被调用")
// 		close(revokeCh)
// 	})

// 	err := leaderAction.Elect(mockLeader)
// 	if err != nil {
// 		t.Fatalf("无法选举领导者: %v", err)
// 	}

// 	<-electCh
// 	if !leaderAction.IsElected() {
// 		t.Fatalf("预期成为领导者")
// 	}

// 	leaderAction.Watch()

// 	err = leaderAction.Resign()
// 	if err != nil {
// 		t.Fatalf("无法放弃领导者地位: %v", err)
// 	}

// 	<-revokeCh
// 	if leaderAction.IsElected() {
// 		t.Fatalf("放弃后预期不再是领导者")
// 	}

// 	leaderAction.Stop()
// }

/* *****************************************上下分割********************************************** */
var ErrAlreadyElected = errors.New("a leader has already been elected for this key")

// MockLeader is a simple implementation of a Leader interface that simulates leader election.
type MockLeader struct {
	elected map[string]*MockElected
	mutex   sync.Mutex // Protects access to the elected map
}

func NewMockLeader() *MockLeader {
	return &MockLeader{
		elected: make(map[string]*MockElected),
	}
}

// Elect tries to elect a leader for a given key. If a leader already exists, it fails.
func (m *MockLeader) Elect(key string, val interface{}) (Elected, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.elected[key]; exists {
		// Simulate an election fail because a leader is already elected for this key
		log.Printf("选举失败: %s 已经有领导者", key)
		return nil, ErrAlreadyElected
	}

	// Simulate a successful election
	log.Printf("模拟选举: %s 成为领导者", val)
	elected := &MockElected{
		key:     key,
		val:     val,
		revoked: make(chan struct{}),
	}
	m.elected[key] = elected
	return elected, nil
}

func (m *MockLeader) Resign(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if elected, exists := m.elected[key]; exists {
		// Simulate the leader resigning
		log.Printf("模拟领导者放弃: %s 不再是领导者", elected.val)
		close(elected.revoked)
		delete(m.elected, key)
	}
}

type MockElected struct {
	key     string
	val     interface{}
	revoked chan struct{}
}

func (m *MockElected) Resign() error {
	// Simulated resignation of elected leader
	log.Printf("模拟领导者：%v 主动放弃领导者地位", m.val)
	close(m.revoked)
	return nil
}

func (m *MockElected) Revoked() <-chan struct{} {
	return m.revoked
}

func TestLeaderElection(t *testing.T) {
	mockLeader := NewMockLeader()
	// 启动多个节点的选举
	for i := 0; i < 3; i++ {
		nodeID := i
		go func(id int) {
			// 每个节点尝试选举自己为领导者
			key := "test_key" // 在真实场景中，这可能代表一个资源或服务的标识
			val := "node_" + strconv.Itoa(id)
			elected, err := mockLeader.Elect(key, val)
			if err != nil {
				if err == ErrAlreadyElected {
					log.Printf("节点 %d 选举失败，已存在领导者", id)
				} else {
					log.Printf("节点 %d 选举过程中发生错误: %v", id, err)
				}
				return
			}

			// 假设节点进行了一些工作作为领导者
			time.Sleep(1 * time.Second)

			// 然后自动放弃领导者地位
			err = elected.Resign()
			if err != nil {
				log.Printf("节点 %d 放弃领导者地位时发生错误: %v", id, err)
			}

			log.Printf("节点 %d 放弃了领导者地位", id)
		}(nodeID)
	}

	// 等待足够的时间，以确保选举过程完成
	time.Sleep(5 * time.Second)
}
