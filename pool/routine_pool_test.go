package pool

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestNoDelayWork(t *testing.T) {
	var count int32 = 0
	task := &Task{
		Handle: func(i interface{}) {
			time.Sleep(time.Second * 2)
			c := i.(*int32)
			atomic.AddInt32(c, 1)
		},
		Arg: &count,
	}

	fmt.Println("pool count:", DefaultRoutinePool.Count())
	fmt.Println("pool free:", DefaultRoutinePool.Free())

	for i := 0; i < 1000; i++ {
		DefaultRoutinePool.DoTask(task)
	}

	fmt.Println("pool count:", DefaultRoutinePool.Count())
	fmt.Println("pool free:", DefaultRoutinePool.Free())

	for i := 0; i < 1000; i++ {
		DefaultRoutinePool.PushTask(task)
	}

	fmt.Println("pool count:", DefaultRoutinePool.Count())
	fmt.Println("pool free:", DefaultRoutinePool.Free())
	time.Sleep(time.Second * 4)

	fmt.Println("pool count:", DefaultRoutinePool.Count())
	fmt.Println("pool free:", DefaultRoutinePool.Free())
	fmt.Println("count:", count)
}
