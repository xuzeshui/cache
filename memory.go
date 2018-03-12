package cache

import (
	"errors"
	"sync"
	"time"
)

type MemoryConf struct {
	CheckPeriodMs int //周期性检查,单位ms
	TimeoutFn     func(key string, value interface{})
}

// MemoryItem store memory cache item.
type MemoryItem struct {
	val         interface{}
	createdTime time.Time
	lifespan    time.Duration
}

func (mi *MemoryItem) isExpire() bool {
	// 0 means forever
	if mi.lifespan == 0 {
		return false
	}
	return time.Now().Sub(mi.createdTime) > mi.lifespan
}

// MemoryCache is Memory cache adapter.
// it contains a RW locker for safe map storage.
type MemoryCache struct {
	sync.RWMutex
	items  map[string]*MemoryItem
	config *MemoryConf
}

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() Cache {
	cache := MemoryCache{items: make(map[string]*MemoryItem)}
	return &cache
}

// Get cache from memory.
// if non-existed or expired, return nil.
func (mc *MemoryCache) Get(name string) interface{} {
	mc.RLock()
	defer mc.RUnlock()
	if itm, ok := mc.items[name]; ok {
		if itm.isExpire() {
			return nil
		}
		return itm.val
	}
	return nil
}

// GetMulti gets caches from memory.
// if non-existed or expired, return nil.
func (mc *MemoryCache) GetMulti(names []string) []interface{} {
	var rc []interface{}
	for _, name := range names {
		rc = append(rc, mc.Get(name))
	}
	return rc
}

// Put cache to memory.
// if lifespan is 0, it will be forever till restart.
func (mc *MemoryCache) Put(name string, value interface{}, lifespan time.Duration) error {
	mc.Lock()
	defer mc.Unlock()
	mc.items[name] = &MemoryItem{
		val:         value,
		createdTime: time.Now(),
		lifespan:    lifespan,
	}
	return nil
}

// Delete cache in memory.
func (mc *MemoryCache) Delete(name string) error {
	mc.Lock()
	defer mc.Unlock()
	if _, ok := mc.items[name]; !ok {
		return errors.New("key not exist")
	}
	delete(mc.items, name)
	if _, ok := mc.items[name]; ok {
		return errors.New("delete key error")
	}
	return nil
}

// Incr increase cache counter in memory.
// it supports int,int32,int64,uint,uint32,uint64.
func (mc *MemoryCache) Incr(key string) error {
	mc.RLock()
	defer mc.RUnlock()
	itm, ok := mc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) + 1
	case int32:
		itm.val = itm.val.(int32) + 1
	case int64:
		itm.val = itm.val.(int64) + 1
	case uint:
		itm.val = itm.val.(uint) + 1
	case uint32:
		itm.val = itm.val.(uint32) + 1
	case uint64:
		itm.val = itm.val.(uint64) + 1
	default:
		return errors.New("item val is not (u)int (u)int32 (u)int64")
	}
	return nil
}

// Decr decrease counter in memory.
func (mc *MemoryCache) Decr(key string) error {
	mc.RLock()
	defer mc.RUnlock()
	itm, ok := mc.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	switch itm.val.(type) {
	case int:
		itm.val = itm.val.(int) - 1
	case int64:
		itm.val = itm.val.(int64) - 1
	case int32:
		itm.val = itm.val.(int32) - 1
	case uint:
		if itm.val.(uint) > 0 {
			itm.val = itm.val.(uint) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint32:
		if itm.val.(uint32) > 0 {
			itm.val = itm.val.(uint32) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	case uint64:
		if itm.val.(uint64) > 0 {
			itm.val = itm.val.(uint64) - 1
		} else {
			return errors.New("item val is less than 0")
		}
	default:
		return errors.New("item val is not int int64 int32")
	}
	return nil
}

// IsExist check cache exist in memory.
func (mc *MemoryCache) IsExist(name string) bool {
	mc.RLock()
	defer mc.RUnlock()
	if v, ok := mc.items[name]; ok {
		return !v.isExpire()
	}
	return false
}

// ClearAll will delete all cache in memory.
func (mc *MemoryCache) ClearAll() error {
	mc.Lock()
	defer mc.Unlock()
	mc.items = make(map[string]*MemoryItem)
	return nil
}

// StartAndGC start memory cache. it will check expiration in every clock time.
func (mc *MemoryCache) StartAndGC(config interface{}) error {
	mc.config = config.(*MemoryConf)
	go mc.checkPeriod()
	return nil
}

// check expiration.
func (mc *MemoryCache) checkPeriod() {
	if mc.config.CheckPeriodMs < 40 {
		mc.config.CheckPeriodMs = 40
	}
	dur := time.Duration(mc.config.CheckPeriodMs) * time.Millisecond
	for {
		<-time.After(dur)
		if mc.items == nil {
			return
		}
		if keys := mc.expiredKeys(); len(keys) > 0 {
			mc.clearItems(keys)
		}
	}
}

// expiredKeys returns key list which are expired.
func (mc *MemoryCache) expiredKeys() (keys []string) {
	mc.RLock()
	defer mc.RUnlock()
	for key, itm := range mc.items {
		if itm.isExpire() {
			keys = append(keys, key)
			if mc.config.TimeoutFn != nil {
				go mc.config.TimeoutFn(key, itm.val)
			}
		}
	}
	return
}

// clearItems removes all the items which key in keys.
func (mc *MemoryCache) clearItems(keys []string) {
	mc.Lock()
	defer mc.Unlock()
	for _, key := range keys {
		delete(mc.items, key)
	}
}

func init() {
	Register("memory", NewMemoryCache)
}
