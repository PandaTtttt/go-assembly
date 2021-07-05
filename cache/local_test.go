package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLocal(t *testing.T) {
	i := 0
	injectError := false
	l := NewLocal(time.Millisecond*50,
		func() (interface{}, error) {
			i++
			if injectError {
				return 0, fmt.Errorf("error %d", i)
			}
			return i, nil
		}, )
	v, err := l.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, v.(int))

	time.Sleep(time.Millisecond * 100)
	// This should returns the new value since old value has expired.
	v, err = l.Get()
	assert.NoError(t, err)
	assert.Equal(t, 2, v.(int))

	injectError = true
	// wait for expired.
	time.Sleep(time.Millisecond * 100)
	v, err = l.Get()
	// it should returns error.
	assert.Equal(t, "error 3", err.Error())

	// get again, since last one was having error.
	// the loader function will be called.
	v, err = l.Get()
	assert.Equal(t, "error 4", err.Error())

	// When the error is clear, the value will be
	// retrieved immediately.
	injectError = false
	v, err = l.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)

	// get again
	v, err = l.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
}

//BenchmarkLocal_Get   	 20000     79 ns/op
func BenchmarkLocal_Get(b *testing.B) {
	b.StopTimer()
	l := NewLocal(time.Minute,
		func() (interface{}, error) {
			time.Sleep(time.Millisecond * 5)
			return 1, nil
		}, )
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := l.Get()
		if err != nil {
			panic(err)
		}
	}
}
