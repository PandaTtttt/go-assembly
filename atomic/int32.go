package atomic

import "sync/atomic"

type Int32 struct {
	v int32
}

func (i *Int32) Add(n int32) int32 {
	return atomic.AddInt32(&i.v, n)
}

func (i *Int32) Swap(n int32) int32 {
	return atomic.SwapInt32(&i.v, n)
}

func (i *Int32) CAS(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&i.v, old, new)
}

func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}

func (i *Int32) Store(n int32) {
	atomic.StoreInt32(&i.v, n)
}
