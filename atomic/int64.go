package atomic

import "sync/atomic"

type Int64 struct {
	v int64
}

func (i *Int64) Add(n int64) int64 {
	return atomic.AddInt64(&i.v, n)
}

func (i *Int64) Swap(n int64) int64 {
	return atomic.SwapInt64(&i.v, n)
}

func (i *Int64) CAS(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&i.v, old, new)
}

func (i *Int64) Load() int64 {
	return atomic.LoadInt64(&i.v)
}

func (i *Int64) Store(n int64) {
	atomic.StoreInt64(&i.v, n)
}
