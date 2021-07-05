package queue

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
	"time"
)

// only for testing stat
var taskCount int
var mu sync.Mutex

func Do1() {
	time.Sleep(time.Millisecond * 30)

	mu.Lock()
	taskCount++
	mu.Unlock()
}

func Do2() {
	time.Sleep(time.Millisecond * 50)

	mu.Lock()
	taskCount++
	mu.Unlock()
}

func DoPanic() {
	panic("panic")
}

func DoError() error {
	return errors.New("this is an exec error")
}

func TestHandleTask(t *testing.T) {
	maxProcs := 200
	q := New(maxProcs, 1000)
	q.Run()

	for i := 0; i < 500; i++ {
		q.Push("DoPanic", func() error {
			DoPanic()
			return nil
		})
	}

	for i := 0; i < 500; i++ {
		q.Push("DoError", func() error {
			return DoError()
		})
	}

	for i := 0; i < 5000; i++ {
		q.Push("Do1", func() error {
			Do1()
			return nil
		})
		q.Push("Do2", func() error {
			Do2()
			return nil
		})
	}

	time.Sleep(time.Second * 1)
	assert.Equal(t, 10000, taskCount)
	assert.Equal(t, maxProcs+3, runtime.NumGoroutine())
}
