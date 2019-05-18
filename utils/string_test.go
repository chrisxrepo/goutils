package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBytes2Str(t *testing.T) {

	Convey("bytes2str", t, func() {
		bytes := []byte("Hello World!")
		str := Bytes2Str(bytes)
		So(len(bytes), ShouldEqual, len(str))
		for i := 0; i < len(str); i++ {
			So(bytes[i], ShouldEqual, str[i])
		}

		Println()
		Println(str)
		bytes[0] = 'h'
		Println(str)
	})
}

func TestStr2Bytes(t *testing.T) {
	Convey("str2bytes", t, func() {
		str := "Hello World!"
		bytes := Str2Bytes(str)
		So(len(bytes), ShouldEqual, len(str))

		for i := 0; i < len(str); i++ {
			So(bytes[i], ShouldEqual, str[i])
		}

		//bytes[0] = 'h'   //panic
		bytes = append(bytes, 'a')
	})
}
