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
	timers      map[string]*timingwheel.Timer
	timingWheel *timingwheel.TimingWheel
}

func NewKVStore(tick time.Duration, wheelSize int64) *KVStore {
	tw := timingwheel.NewTimingWheel(tick, wheelSize)
	tw.Start()
	return &KVStore{
		data:        make(map[string]int),
		timers:      make(map[string]*timingwheel.Timer),
		timingWheel: tw,
	}
}

// IndexIncIfExist 存在就 inc 值，不存在则设置为 0
func (kv *KVStore) IndexIncIfExist(key string) int {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	i, exist := kv.data[key]
	if !exist {
		// 键不存在，设置为 0，并启动计时器
		kv.data[key] = 0
		kv.startTimer(key, time.Minute/2)
		return 0
	} else {
		// 键已存在，增加值，但不影响计时器
		kv.data[key] = i + 1
		return i + 1
	}
}

func (kv *KVStore) SetDefault(key string, value int) {
	kv.Set(key, value, time.Minute/2)
}

func (kv *KVStore) Set(key string, value int, expiration time.Duration) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, exists := kv.data[key]
	kv.data[key] = value

	// 只有当键不存在时，才启动新的计时器
	if !exists {
		kv.startTimer(key, expiration)
	}
}

// TODO 测试
func (kv *KVStore) startTimer(key string, expiration time.Duration) {
	// 启动新的定时任务，指定时间后删除键值对
	timer := kv.timingWheel.AfterFunc(expiration, func() {
		kv.mu.Lock()
		defer kv.mu.Unlock()
		delete(kv.data, key)
		delete(kv.timers, key)
		fmt.Printf("Key '%s' has been deleted after %v\n", key, expiration)
	})

	//存储定时器句柄
	kv.timers[key] = timer
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
	if timer, exists := kv.timers[key]; exists {
		timer.Stop()
		delete(kv.timers, key)
	}
}

func (kv *KVStore) Close() {
	kv.timingWheel.Stop()
}
