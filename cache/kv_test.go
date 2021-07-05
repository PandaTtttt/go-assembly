package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type counter struct {
	mu        sync.Mutex
	numCalled int
}

func (c *counter) Inc() {
	c.mu.Lock()
	c.numCalled++
	c.mu.Unlock()
}

func (c *counter) Value() int {
	c.mu.Lock()
	v := c.numCalled
	c.mu.Unlock()
	return v
}

func TestKV(t *testing.T) {
	c := counter{}
	injectError := false
	cache := NewKV(time.Millisecond*50, time.Millisecond*200,
		func(key interface{}) (interface{}, error) {
			c.Inc()
			if injectError {
				return "", fmt.Errorf("error %d", c.Value())
			}
			return "echo " + key.(string), nil
		})

	v, err := cache.Get("tom")
	assert.NoError(t, err)
	assert.Equal(t, "echo tom", v.(string))
	assert.Equal(t, 1, c.Value())

	v, err = cache.Get("jerry")
	assert.NoError(t, err)
	assert.Equal(t, "echo jerry", v.(string))
	assert.Equal(t, 2, c.Value())

	// wait for staled.
	time.Sleep(time.Millisecond * 100)
	v, err = cache.Get("tom")
	assert.NoError(t, err)
	assert.Equal(t, "echo tom", v.(string))
	// since old value of tom has staled, loader should be recalled.
	assert.Equal(t, 3, c.Value())

	// wait for expired.
	time.Sleep(200 * time.Millisecond)
	v, err = cache.Get("rick")
	// wait for GC
	time.Sleep(time.Millisecond * 10)
	assert.NoError(t, err)
	assert.Equal(t, "echo rick", v.(string))
	assert.Equal(t, 4, c.Value())
	// since tom and jerry have expired,
	// they should be GC after called rick.
	assert.Equal(t, 1, cache.Size())

	v, err = cache.Get("tom")
	v, err = cache.Get("jerry")
	// wait for expired.
	time.Sleep(200 * time.Millisecond)
	v, err = cache.Get("rick")
	// wait for GC
	time.Sleep(time.Millisecond * 10)
	assert.NoError(t, err)
	assert.Equal(t, "echo rick", v.(string))
	assert.Equal(t, 7, c.Value())
	// since tom and jerry have expired,
	// they should be GC after called rick.
	assert.Equal(t, 1, cache.Size())

	injectError = true
	v, err = cache.Get("morty")
	assert.Equal(t, "error 8", err.Error())

	// get again, since last one was having error.
	// the loader function will be called.
	v, err = cache.Get("morty")
	assert.Equal(t, "error 9", err.Error())

	injectError = false
	v, err = cache.Get("morty")
	assert.NoError(t, err)
	assert.Equal(t, "echo morty", v.(string))
	assert.Equal(t, 10, c.Value())

	// get again
	v, err = cache.Get("morty")
	assert.NoError(t, err)
	assert.Equal(t, "echo morty", v.(string))
	assert.Equal(t, 10, c.Value())
}

// BenchmarkKV_Get 	   10000    191 ns/op
func BenchmarkKV_Get(b *testing.B) {
	b.StopTimer()
	l := NewKV(time.Minute, 0,
		func(key interface{}) (interface{}, error) {
			time.Sleep(time.Millisecond * 5)
			return 1, nil
		}, )
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := l.Get("test")
		if err != nil {
			panic(err)
		}
	}
}
