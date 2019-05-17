package error

import (
	"fmt"
	"runtime"

	"github.com/chrisxrepo/goutils/log"
)

func RecoverPanic() {
	if err := recover(); err != nil {
		var stack string
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		}
		//fmt.Println(stack)
		log.DefaultLog.Error(stack)
	}
}
