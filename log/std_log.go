package log

import "fmt"

type stdLog int

func (l *stdLog) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (l *stdLog) Error(args ...interface{}) {
	fmt.Println(args...)
}
