package utils

import (
	"reflect"
	"unsafe"
)

func Bytes2Str(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

//危险方法
//修改string转化后的[]byte内部数值会带来panic
func Str2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
