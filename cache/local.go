package cache

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Local struct {
	loader func() (interface{}, error)
	ttl    time.Duration

	v atomic.Value

	mu sync.Mutex
}

type localV struct {
	v      interface{}
	err    error
	panicV error
	loaded time.Time
}

func NewLocal(ttl time.Duration, loader func() (interface{}, error)) *Local {
	return &Local{
		loader: loader,
		ttl:    ttl,
	}
}

// Get returns the cached value.
func (l *Local) Get() (interface{}, error) {
	if l.needLoad() {
		l.loadInternal()
	}
	v := l.v.Load().(*localV)

	if v.panicV != nil {
		panic(v.panicV)
	}
	if v.err != nil {
		return nil, v.err
	}
	return v.v, nil
}

func (l *Local) loadInternal() {
	l.mu.Lock()
	defer l.mu.Unlock()
	// duplicate suppression
	if !l.needLoad() {
		return
	}

	lv := &localV{}
	// call loader() and handle the panic.
	func() {
		// if panic, record error
		defer func() {
			v := recover()
			if v != nil {
				err := errors.New(fmt.Sprintf("recover panic:%v", v))
				lv.panicV = err
			}
		}()
		// if loader returns error，record error but keep value empty
		v, err := l.loader()
		if err != nil {
			lv.err = err
			return
		}
		lv.v = v
	}()
	lv.loaded = time.Now()
	l.v.Store(lv)
}

// needLoad returns true if there is no value of v,
// or there is value but with error, panic or expiration.
func (l *Local) needLoad() bool {
	v := l.v.Load()
	if v == nil {
		return true
	}
	lv := v.(*localV)
	if lv.err != nil || lv.panicV != nil ||
		lv.loaded.Before(time.Now().Add(-l.ttl)) {
		return true
	}
	return false
}

// loadedTime returns when the current value is loaded, or zero if never.
func (l *Local) loadedTime() time.Time {
	v := l.v.Load()
	if v == nil {
		return time.Time{}
	}
	return v.(*localV).loaded
}
