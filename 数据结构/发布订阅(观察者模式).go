package main

import (
	"fmt"
	"sync"
)

func main() {
	// 第一种: 单key对应单value
	singleKVSubscribe("卖火柴", func(param string) interface{} {
		fmt.Println(param + "卖火柴的小女孩")
		return nil
	})
	singleKVSubscribe("卖油条", func(param string) interface{} {
		fmt.Println(param + "卖油条的小男孩")
		return nil
	})
	singleKVPublish("卖油条", "帅气的")

	// 第二种: 单key对应多value
	multiKVWathcer("吃早饭", func(param string) interface{} {
		fmt.Println(param + ", 吃面包")
		return nil
	})
	multiKVWathcer("吃早饭", func(param string) interface{} {
		fmt.Println(param + ", 吃卤蛋")
		return nil
	})
	multiKVWathcer("吃早饭", func(param string) interface{} {
		fmt.Println(param + ", 吃泡菜")
		return nil
	})
	multiKVNotify("吃早饭", "在家里")
}

// 第一种: 单key对应单value
var singleKVMap = make(map[string]interface{})
var singleLock sync.RWMutex

type valueFunc func(string) interface{}

func singleKVSubscribe(key string, fn valueFunc) {
	singleLock.Lock()
	defer singleLock.Unlock()

	singleKVMap[key] = fn
}

func singleKVPublish(key string, param string) interface{} {
	singleLock.RLock()
	defer singleLock.RUnlock()

	fn, ok := singleKVMap[key]
	if !ok {
		return nil
	}

	return fn.(valueFunc)(param)
}

// 第二种: 单key对应多value
var multiKVMap = make(map[string][]interface{})
var multiLock sync.RWMutex

func multiKVWathcer(key string, fn valueFunc) {
	multiLock.Lock()
	defer multiLock.Unlock()

	multiKVMap[key] = append(multiKVMap[key], fn)
}

func multiKVNotify(key string, param string) {
	multiLock.RLock()
	defer multiLock.RUnlock()

	fns, ok := multiKVMap[key]
	if !ok {
		return
	}

	for _, fn := range fns {
		fn.(valueFunc)(param)
	}
}
