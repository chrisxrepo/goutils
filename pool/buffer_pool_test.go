package pool

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBufferPool(t *testing.T) {
	Convey("buffer_pool", t, func() {
		allocCount := 1000
		bufs := make([]*ByteBuffer, allocCount)
		for i := 0; i < allocCount; i++ {
			buffer := DefaultBufferPool.Get()
			bufs[i] = buffer
		}
		Println("buffer pool cap1:", DefaultBufferPool.Cap())
		Println("buffer pool len1:", DefaultBufferPool.Len())

		for i := 0; i < allocCount; i++ {
			DefaultBufferPool.Put(bufs[i])
		}
		Println("buffer pool cap2:", DefaultBufferPool.Cap())
		Println("buffer pool len2:", DefaultBufferPool.Len())

		buffer := DefaultBufferPool.Get()
		So(buffer, ShouldNotBeNil)

		buffer.SetString("Hello World")
		fmt.Println(buffer.String())
		DefaultBufferPool.Put(buffer)

	})

}

func TestByteBuffer(t *testing.T) {
	Convey("buffer_pool", t, func() {
		buf := DefaultBufferPool.Get()
		buf.WriteUint16(100)
		buf.WriteUint32(2000)
		buf.WriteUint64(30000)
		buf.WriteString("hello")

		So(buf.PickUint16(), ShouldEqual, 100)
		So(buf.ReadUint16(), ShouldEqual, 100)
		So(buf.PickUint32(), ShouldEqual, 2000)
		So(buf.ReadUint32(), ShouldEqual, 2000)
		So(buf.PickUint64(), ShouldEqual, 30000)
		So(buf.ReadUint64(), ShouldEqual, 30000)
		So(string(buf.PickBytes(5)), ShouldEqual, "hello")
		So(string(buf.ReadBytes(5)), ShouldEqual, "hello")
	})
}
