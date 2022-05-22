package queue

import (
	"fmt"
	"testing"
	"time"
)

func startTask() {
	fmt.Println("start task")
}

func endTask() {
	fmt.Println("end task")
}
func errTask(err error) {
	fmt.Println("err task", err)
}

func TestHandleTaskX(t *testing.T) {
	maxProcs := 1
	q := NewX(maxProcs, 1)
	q.Run()

	for i := 0; i < 2; i++ {
		q.Push("DoPanic", func() error {
			DoPanic()
			return nil
		}, startTask, endTask, errTask)
	}

	for i := 0; i < 2; i++ {
		q.Push("DoError", func() error {
			return DoError()
		}, startTask, endTask, errTask)
	}

	for i := 0; i < 10; i++ {
		q.Push("Do1", func() error {
			Do1()
			return nil
		}, startTask, endTask, errTask)
		q.Push("Do2", func() error {
			Do2()
			return nil
		}, startTask, endTask, errTask)
	}
	time.Sleep(time.Second * 1)
}
