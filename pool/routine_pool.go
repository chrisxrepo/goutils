package pool

import (
	"sync/atomic"

	"github.com/chrisxrepo/goutils/error"
)

const (
	DefaultRoutinePoolSize = 10000
)

type Task struct {
	Handle func(interface{})
	Arg    interface{}
}

type RoutinePool struct {
	size         int32
	free         int32
	taskChan     chan *Task
	waitTaskChan chan *Task
	checkChan    chan struct{}
}

func NewRoutinePool(size int) *RoutinePool {
	rp := &RoutinePool{
		size:         int32(size),
		free:         int32(size),
		taskChan:     make(chan *Task),
		waitTaskChan: make(chan *Task, 100000),
		checkChan:    make(chan struct{}),
	}

	for i := 0; i < size; i++ {
		go rp.routineRun()
	}

	return rp
}

func (p *RoutinePool) DoTask(task *Task) bool {
	select {
	case p.taskChan <- task:
		return true
	default:
		return false
	}
}

func (p *RoutinePool) PushTask(task *Task) {
	p.waitTaskChan <- task
}

func (p *RoutinePool) Count() int {
	return int(p.size)
}

func (p *RoutinePool) Free() int {
	return int(p.free)
}

func (p *RoutinePool) routineRun() {
	for {
		select {
		case task := <-p.taskChan:
			p.handleTask(task)

		case task := <-p.waitTaskChan:
			p.handleTask(task)

		case <-p.checkChan:
			break

		}
	}
}

func (p *RoutinePool) handleTask(task *Task) {
	defer error.RecoverPanic()

	atomic.AddInt32(&p.free, -1)
	task.Handle(task.Arg)
	atomic.AddInt32(&p.free, 1)
}

func (p *RoutinePool) Stop() {
	close(p.checkChan)
}

var DefaultRoutinePool = NewRoutinePool(DefaultRoutinePoolSize)
