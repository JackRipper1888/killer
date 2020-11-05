package mapkit

import (
	"sync"
	"sync/atomic"
)

type syncMap struct {
	ConcurrentMap []*ConcurrentSyncMapShared
	SHARE_COUNT   int
}

type ConcurrentSyncMapShared struct {
	items     *sync.Map // 本分片内的map
	MAP_COUNT int64     //长度
}

// 新建一个map
func NewConcurrentSyncMap(share_cont int) *syncMap {
	m := make([]*ConcurrentSyncMapShared, share_cont)
	for i := 0; i < share_cont; i++ {
		m[i] = &ConcurrentSyncMapShared{
			items:     &sync.Map{},
			MAP_COUNT: 0,
		}
	}
	body := syncMap{
		ConcurrentMap: m,
		SHARE_COUNT:   share_cont}

	return &body
}

// GetSharedMap 获取key对应的map分片
func (m *syncMap) getSharedMap(key string) *ConcurrentSyncMapShared {
	return m.ConcurrentMap[uint(fnv32(key))%uint(m.SHARE_COUNT)]
}

// Set 设置key,value
func (m *syncMap) Set(key string, value interface{}) {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	_, isOk := sharedMap.items.LoadOrStore(key, value)
	if isOk {
		//存在添加
		sharedMap.items.Store(key, value)
	} else {
		atomic.AddInt64(&sharedMap.MAP_COUNT, 1)
	}
}
func (m *syncMap) LoadOrStore(key string, value interface{}) {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	_, isOk := sharedMap.items.LoadOrStore(key, value)
	if !isOk {
		atomic.AddInt64(&sharedMap.MAP_COUNT, 1)
	}
}
func (m *syncMap) Get(key string) (value interface{}, ok bool) {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	return sharedMap.items.Load(key)
}

func (m *syncMap) Range(f func(k string, v interface{}) bool) {
	for _, mp := range m.ConcurrentMap {
		mp.items.Range(func(key, value interface{}) bool {
			return f(key.(string), value)
		})
	}
}

// Get 获取key对应的value删除对应的data
func (m *syncMap) Delete(key string) {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	atomic.AddInt64(&sharedMap.MAP_COUNT, -1)
	sharedMap.items.Delete(key)
}

// Count 统计key个数
func (m *syncMap) Count() int64 {
	var count int64
	for i := 0; i < m.SHARE_COUNT; i++ {
		count += atomic.LoadInt64(&m.ConcurrentMap[i].MAP_COUNT)
	}
	return count
}
