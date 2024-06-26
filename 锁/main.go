package main

import "sync"

func main() {

}

// Sync.Mutex：可以在多个goroutine间同步访问共享资源，确保同一时刻只有一个goroutine访问共享资源
var mu sync.Mutex

func writeToMap(m map[string]string, key string, value string) {
	mu.Lock()         // 加锁
	defer mu.Unlock() // 解锁
	m[key] = value    // 写入 map
}

/*
	Sync.RWMutex: 相比于Sync.Mutex更加灵活，适用于读多写少的场景，两种锁定方式：
	1. 读锁(RLock)：多个goroutine可以同时持有读锁
	2. 写锁(Lock)：同一时间，只能有一个goroutine持有写锁，并且持有写锁时，其他goroutine无法持有读锁或写锁
*/
var rwMu sync.RWMutex

func readFromMapRW(m map[string]string, key, value string) string {
	rwMu.RLock()         // 读锁
	defer rwMu.RUnlock() // 释放读锁
	return m[key]        // 读取 map
}

func writeToMapRW(m map[string]string, key string, value string) {
	rwMu.Lock()         // 写锁
	defer rwMu.Unlock() // 释放写锁
	m[key] = value      // 写入 map
}

// Sync.Once: 保证只执行一次操作，适用于初始化操作。例如某个函数只会被调用一次
var once sync.Once

func init() {
	once.Do(func() {
		// 初始化操作内容
	})
}

/*
	Sync.Cond: 提供一种条件变量，可以让goroutine等待或者通知其他goroutine
	1. Wait()：等待通知
	2. Signal()：通知单个goroutine
	3. Broadcast()：通知所有goroutine
	更多完整使用：在/锁/cond_main.go中查看
*/
var condMu sync.Mutex
var cond = sync.NewCond(&condMu)

func waitCondtion() {
	condMu.Lock()
	cond.Wait() // 等待通知
	condMu.Unlock()
}

func notifyCondtion() {
	condMu.Lock()
	cond.Signal() // 通知单个goroutine
	condMu.Unlock()
}

/*
	Sync.Map: 是Go1.9版本新增的一种并发安全的map
	1. 无需初始化，直接声明即可
	2. 提供了高效的读写操作，避免使用`Sync.Mutex`或`Sync.RWMutex`手动管理锁的麻烦
*/
var smap sync.Map

func writeToSyncMap(key, value string) {
	smap.Store(key, value)
}

func readFromSyncMap(key string) (string, bool) {
	value, ok := smap.Load(key)
	if !ok {
		return "", false
	}
	return value.(string), true
}
