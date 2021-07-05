package queue

import (
	"github.com/PandaTtttt/go-assembly/util/unwind"
	"github.com/PandaTtttt/go-assembly/zlog"
	"go.uber.org/zap"
)

type Queue struct {
	maxProcs int

	goCh    chan struct{}
	taskQue chan struct {
		name string
		f    func() error
	}
}

// maxProcs sets the maximum number of goroutine
// queBuf sets the buffer size of Task channel
func New(maxProcs, queBuf int) *Queue {
	q := &Queue{
		maxProcs: maxProcs,
		goCh:     make(chan struct{}, maxProcs),
		taskQue: make(chan struct {
			name string
			f    func() error
		}, queBuf),
	}
	return q
}

func (q *Queue) Run() {
	for i := 0; i < q.maxProcs; i++ {
		go q.worker()
	}
	go func() {
		for {
			<-q.goCh
			go q.worker()
		}
	}()
}

func (q *Queue) worker() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Error("Unexpected panic",
				zap.Any("error", err),
				zap.String("stack", unwind.Stack(3)))
		}
		q.goCh <- struct{}{}
	}()
	for {
		t := <-q.taskQue
		err := t.f()
		if err != nil {
			zlog.Error(t.name, zap.Error(err))
		}
	}
}

func (q *Queue) Push(name string, f func() error) {
	q.taskQue <- struct {
		name string
		f    func() error
	}{name: name, f: f}
}
