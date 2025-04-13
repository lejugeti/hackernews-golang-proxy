package cache

import (
	"sync"
	"time"
)

// Cache that automatically delete its entries after a time to live period
type TimeToLiveCache[K comparable, V any] struct {
	data map[K]V
	mutex sync.RWMutex
	timeToLive time.Duration
	timeToLiveTimers map[K]*time.Timer
}

func NewTimeToLiveCache[K comparable, V any](timeToLive time.Duration) *TimeToLiveCache[K, V] {
	return &TimeToLiveCache[K, V]{
		data: make(map[K]V),
		timeToLive: timeToLive,
		timeToLiveTimers: make(map[K]*time.Timer),
	}
}

func (cache *TimeToLiveCache[K, V]) Add(key K, value V) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data[key] = value

	cache.timeToLiveTimers[key] = time.AfterFunc(cache.timeToLive, func() {
		cache.Delete(key)
	})
}

func (cache *TimeToLiveCache[K, V]) Get(key K) (V, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	value, ok := cache.data[key]
	return value, ok
}

func (cache *TimeToLiveCache[K, V]) Delete(key K) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.data, key)

	cache.stopTtlTimer(key)
}

// Stops and deletes timer associated to key. 
// Used to prevent timer from trigerring entry deletion while it has already been done
func (cache *TimeToLiveCache[K, V]) stopTtlTimer(key K) {
	ttlTimer, timerExists := cache.timeToLiveTimers[key]
	
	if timerExists {
		ttlTimer.Stop()
		delete(cache.timeToLiveTimers, key)
	}
}