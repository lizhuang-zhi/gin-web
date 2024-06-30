package option

import "sync"

type MetaData interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}

type metaData struct {
	sync.RWMutex // 读写锁

	name string

	m map[string]interface{}
}

func NewMetaData(name string) MetaData {
	return &metaData{
		name: name,
		m:    make(map[string]interface{}),
	}
}

func (md *metaData) Set(key string, value interface{}) {
	md.Lock()
	defer md.Unlock()

	md.m[key] = value
}

func (md *metaData) Get(key string) interface{} {
	md.RLock()
	defer md.RUnlock()

	return md.m[key]
}
