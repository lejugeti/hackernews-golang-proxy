package cache

import (
	"testing"
	"time"
)

func TestAddShouldStoreValueInCache(t *testing.T) {
	// GIVEN
	var cache = NewTimeToLiveCache[int, int](10)
	const key = 1
	const value = 2

	// WHEN
	cache.Add(key, value)

	// THEN
	actual, isCached := cache.Get(key)

	if !isCached {
		t.Error("value should be cached")
	} else if actual != value {
		t.Errorf("expected cached value '%d' but got '%d' instead", key, actual)
	}
}

func TestAddShouldStoreValueInCacheEvenIfNil(t *testing.T) {
	// GIVEN
	var cache = NewTimeToLiveCache[int, *int](10)
	const key int = 1

	// WHEN
	cache.Add(key, nil)

	// THEN
	actual, isCached := cache.Get(key)

	if !isCached {
		t.Error("value should be cached")
	} else if actual != nil {
		t.Errorf("expected cached value to be nil but got '%d' instead", actual)
	}
}

func TestGetShouldNotReturnValueIfNotCached(t *testing.T) {
	// GIVEN
	var cache = NewTimeToLiveCache[int, int](10)
	const key int = 1

	// WHEN
	_, isCached := cache.Get(key)

	// THEN

	if isCached {
		t.Error("no value should be cached")
	}
}

func TestDeleteShouldRemoveValueFromCache(t *testing.T) {
	// GIVEN
	var cache = NewTimeToLiveCache[int, int](10)
	const key int = 1
	const value int = 2

	cache.Add(key, value)
	
	// WHEN
	cache.Delete(key)

	// THEN
	_, isCached := cache.Get(key)

	if isCached {
		t.Error("no value should be cached")
	}
}

func TestDeleteEntriesAfterTimeToLiveEllapsed(t *testing.T) {
	// GIVEN
	timeToLive := time.Nanosecond * 10
	cache := NewTimeToLiveCache[int, int](timeToLive)
	const key int = 1
	const value int = 2

	cache.Add(key, value)
	
	// WHEN
	time.Sleep(time.Microsecond * 500)

	_, isCached := cache.Get(key)
	if isCached {
		t.Error("no value should be cached")
	}

	_, timeExists := cache.timeToLiveTimers[key]
	if timeExists {
		t.Error("no timer should exist after time to live")
	}
}