package queue

import (
	"github.com/PandaTtttt/go-assembly/util/unwind"
	"github.com/PandaTtttt/go-assembly/zlog"
	"go.uber.org/zap"
)

type QueueX struct {
	maxProcs int

	goCh    chan struct{}
	taskQue chan struct {
		name   string
		startF func()
		f      func() error
		errF   func(error)
		endF   func()
	}
}

func NewX(maxProcs, queBuf int) *QueueX {
	q := &QueueX{
		maxProcs: maxProcs,
		goCh:     make(chan struct{}, maxProcs),
		taskQue: make(chan struct {
			name   string
			startF func()
			f      func() error
			errF   func(error)
			endF   func()
		}, queBuf),
	}
	return q
}

func (q *QueueX) Run() {
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

func (q *QueueX) worker() {
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
		if t.startF != nil {
			t.startF()
		}
		err := t.f()
		if err != nil {
			if t.errF != nil {
				t.errF(err)
			}
			zlog.Error(t.name, zap.Error(err))
		} else {
			if t.endF != nil {
				t.endF()
			}
		}
	}
}

func (q *QueueX) Push(name string, f func() error, startF func(), endF func(), errF func(error)) {
	q.taskQue <- struct {
		name   string
		startF func()
		f      func() error
		errF   func(error)
		endF   func()
	}{name: name, f: f, startF: startF, endF: endF, errF: errF}
}
