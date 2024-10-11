package kvutil

import (
	"fmt"
	"sync"
	"time"

	"github.com/RussellLuo/timingwheel"
)

type KVStore struct {
	mu          sync.RWMutex
	data        map[string]int
	timingWheel *timingwheel.TimingWheel
}

func NewKVStore(tick time.Duration, wheelSize int64) *KVStore {
	tw := timingwheel.NewTimingWheel(tick, wheelSize)
	tw.Start()
	return &KVStore{
		data:        make(map[string]int),
		timingWheel: tw,
	}
}

// IndexIncIfExist 存在就inc值 不存在set值
func (kv *KVStore) IndexIncIfExist(key string) int {
	i, exist := kv.data[key]
	if !exist {
		kv.SetDefault(key, 0)
		return 0
	} else {
		kv.SetDefault(key, i+1)
		return i + 1
	}
}

func (kv *KVStore) SetDefault(key string, value int) {
	kv.Set(key, value, time.Minute/2)
}

func (kv *KVStore) Set(key string, value int, expiration time.Duration) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[key] = value

	// 设置定时任务，指定时间后删除键值对
	kv.timingWheel.AfterFunc(expiration, func() {
		kv.mu.Lock()
		defer kv.mu.Unlock()
		delete(kv.data, key)
		fmt.Printf("Key '%s' has been deleted after %v\n", key, expiration)
	})
}

func (kv *KVStore) Get(key string) (int, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[key]
	return val, ok
}

func (kv *KVStore) Delete(key string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.data, key)
}

func (kv *KVStore) Close() {
	kv.timingWheel.Stop()
}
