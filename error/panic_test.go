package error

import (
	"fmt"
	"testing"
)

func TestRecoverPanic(t *testing.T) {
	panicFunc()
	fmt.Println("Hell Test Panic")
}

func panicFunc() {
	defer RecoverPanic()

	fmt.Println("panicFunc...")
	panic("hello painic")
}
