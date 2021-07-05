package cache

import (
	"github.com/PandaTtttt/go-assembly/atomic"
	"sync"
	"time"
)

// KV provides a local cache that could fetch value based on a given key.
// When a cached value became stale, a new get will trigger a new fetch.
type KV struct {
	loader func(key interface{}) (interface{}, error)
	stale  time.Duration
	expire time.Duration

	entries map[interface{}]*Local
	mu      sync.RWMutex

	lastGC    time.Time
	gcRunning atomic.Int32
}

// NewKV creates a KV, if expire is smaller than stale, expire equals twice of stale instead.
func NewKV(stale, expire time.Duration, loader func(key interface{}) (interface{}, error)) *KV {
	if expire < stale {
		expire = stale * 2
	}
	return &KV{
		loader:  loader,
		stale:   stale,
		entries: make(map[interface{}]*Local),

		expire: expire,
		lastGC: time.Now(),
	}
}

// Get returns a value from the cache.
func (kv *KV) Get(key interface{}) (interface{}, error) {
	triggerGC := false
	kv.mu.RLock()
	entry, ok := kv.entries[key]
	kv.mu.RUnlock()
	if ok {
		loadedTime := entry.loadedTime()
		// if expired
		if !loadedTime.IsZero() && loadedTime.Before(time.Now().Add(-kv.expire)) {
			ok = false
		}
	}

	if !ok {
		func() {
			kv.mu.Lock()
			defer kv.mu.Unlock()

			triggerGC = true
			entry = NewLocal(kv.stale, func() (interface{}, error) { return kv.loader(key) })
			kv.entries[key] = entry
		}()
	}

	if triggerGC {
		kv.garbageCollect()
	}
	return entry.Get()
}

// garbageCollect performs garbage collection to remove expired items.
func (kv *KV) garbageCollect() {
	if kv.gcRunning.CAS(0, 1) {
		defer kv.gcRunning.Store(0)

		// GC's interval equals twice of stale.
		if kv.lastGC.After(time.Now().Add(-kv.stale * 2)) {
			return
		}
		kv.lastGC = time.Now()

		go func() {
			var expires []interface{}
			cutoff := time.Now().Add(-kv.expire)
			kv.mu.RLock()
			for k, entry := range kv.entries {
				t := entry.loadedTime()
				if t.IsZero() || t.After(cutoff) {
					continue
				}
				expires = append(expires, k)
			}
			kv.mu.RUnlock()

			if len(expires) > 0 {
				kv.mu.Lock()
				for _, k := range expires {
					delete(kv.entries, k)
				}
				kv.mu.Unlock()
			}
		}()
	}
}

// Size returns number of values in the cache.
func (kv *KV) Size() int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return len(kv.entries)
}
