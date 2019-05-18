package pool

import (
	"encoding/binary"
	"io"
	"sync"
	"sync/atomic"
)

const (
	DefaultBufferSize = 4096
	FreeBufferSize    = 1024 * 1024
)

type BufferPool struct {
	pool sync.Pool
	len  int32
	cap  int32
}

func (p *BufferPool) Get() *ByteBuffer {
	v := p.pool.Get()
	if v != nil {
		atomic.AddInt32(&p.len, -1)
		return v.(*ByteBuffer)
	}

	atomic.AddInt32(&p.cap, 1)

	return &ByteBuffer{
		data: make([]byte, 0, DefaultBufferSize),
		pos:  0,
	}
}

func (p *BufferPool) Put(buf *ByteBuffer) {
	if buf == nil || buf.data == nil {
		return
	}

	if cap(buf.data) >= FreeBufferSize {
		buf.data = make([]byte, 0, DefaultBufferSize)
	}

	atomic.AddInt32(&p.len, 1)
	buf.Reset()
	p.pool.Put(buf)
}

func (p *BufferPool) Cap() int {
	return int(p.cap)
}

func (p *BufferPool) Len() int {
	return int(p.len)
}

type ByteBuffer struct {
	data []byte
	pos  int
}

func (b *ByteBuffer) Len() int {
	return len(b.data) - b.pos
}

func (b *ByteBuffer) Cap() int {
	return cap(b.data)
}

func (b *ByteBuffer) Used() int {
	return len(b.data)
}

func (b *ByteBuffer) ReadAll(r io.Reader) (int, error) {
	var nn int
	for {
		if len(b.data) >= cap(b.data) {
			newData := make([]byte, cap(b.data)*2)
			copy(newData, b.data)
			b.data = newData[:len(b.data)]
		}

		n, e := r.Read(b.data[len(b.data):cap(b.data)])
		if n > 0 {
			b.data = b.data[:len(b.data)+n]
		}

		nn += n
		if e != nil {
			if e == io.EOF {
				return nn, nil
			}
			return nn, e
		}
	}
}

func (b *ByteBuffer) ReadFrom(r io.Reader) (int, error) {
	if len(b.data) >= cap(b.data) {
		newData := make([]byte, cap(b.data)*2)
		copy(newData, b.data)
		b.data = newData[:len(b.data)]
	}

	n, e := r.Read(b.data[len(b.data):cap(b.data)])
	if n > 0 {
		b.data = b.data[:len(b.data)+n]
	}

	return n, e
}

func (b *ByteBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.data)
	return int64(n), err
}

func (b *ByteBuffer) Bytes() []byte {
	return b.data[b.pos:]
}

func (b *ByteBuffer) Data() []byte {
	return b.data
}

func (b *ByteBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *ByteBuffer) WriteByte(c byte) error {
	b.data = append(b.data, c)
	return nil
}

func (b *ByteBuffer) WriteString(s string) (int, error) {
	b.data = append(b.data, s...)
	return len(s), nil
}

func (b *ByteBuffer) WriteUint16(v uint16) {
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, v)
	b.Write(tmp)
}

func (b *ByteBuffer) WriteUint32(v uint32) {
	tmp := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, v)
	b.Write(tmp)
}

func (b *ByteBuffer) WriteUint64(v uint64) {
	tmp := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp, v)
	b.Write(tmp)
}

func (b *ByteBuffer) PickByte() byte {
	if 1 > len(b.data)-b.pos {
		return 0
	}

	return b.data[b.pos]
}

func (b *ByteBuffer) ReadByte() byte {
	if 1 > len(b.data)-b.pos {
		return 0
	}

	v := b.data[b.pos]
	b.pos += 1
	return v
}

func (b *ByteBuffer) PickUint16() uint16 {
	if 2 > len(b.data)-b.pos {
		return 0
	}

	return binary.BigEndian.Uint16(b.data[b.pos:])
}

func (b *ByteBuffer) ReadUint16() uint16 {
	if 2 > len(b.data)-b.pos {
		return 0
	}

	v := binary.BigEndian.Uint16(b.data[b.pos:])
	b.pos += 2
	return v
}

func (b *ByteBuffer) PickUint32() uint32 {
	if 4 > len(b.data)-b.pos {
		return 0
	}

	return binary.BigEndian.Uint32(b.data[b.pos:])
}

func (b *ByteBuffer) ReadUint32() uint32 {
	if 4 > len(b.data)-b.pos {
		return 0
	}

	v := binary.BigEndian.Uint32(b.data[b.pos:])
	b.pos += 4
	return v
}

func (b *ByteBuffer) PickUint64() uint64 {
	if 8 > len(b.data)-b.pos {
		return 0
	}

	return binary.BigEndian.Uint64(b.data[b.pos:])
}

func (b *ByteBuffer) ReadUint64() uint64 {
	if 8 > len(b.data)-b.pos {
		return 0
	}

	v := binary.BigEndian.Uint64(b.data[b.pos:])
	b.pos += 8
	return v
}

func (b *ByteBuffer) PickBytes(size int) []byte {
	if size > len(b.data)-b.pos {
		return nil
	}

	return b.data[b.pos : b.pos+size]
}

func (b *ByteBuffer) ReadBytes(size int) []byte {
	if size > len(b.data)-b.pos {
		return nil
	}

	v := b.data[b.pos : b.pos+size]
	b.pos += size
	return v
}

func (b *ByteBuffer) ReadLine(sp string) []byte {
	if len(sp) == 0 {
		sp = "\n"
	}

	for i := b.pos + len(sp) - 1; i < len(b.data); i++ {
		hit := true
		for j := 0; j < len(sp); j++ {
			if b.data[i-len(sp)+1+j] != sp[j] {
				hit = false
				break
			}
		}

		if hit {
			v := b.data[b.pos : i-len(sp)+1]
			b.pos = i + 1
			return v
		}
	}

	return nil
}

func (b *ByteBuffer) Drain(l int) {
	b.pos += l
	if b.pos < 0 {
		b.pos = 0
	} else if b.pos > len(b.data) {
		b.pos = len(b.data)
	}
}

func (b *ByteBuffer) Set(p []byte) {
	b.data = append(b.data[:0], p...)
}

func (b *ByteBuffer) SetString(s string) {
	b.data = append(b.data[:0], s...)
}

func (b *ByteBuffer) String() string {
	return string(b.data[b.pos:])
}

func (b *ByteBuffer) Reset() {
	b.data = b.data[:0]
	b.pos = 0
}

func (b *ByteBuffer) Compact() {
	if b.pos == 0 {
		return
	}

	if b.pos >= len(b.data) {
		b.Reset()
	}

	v := b.data[b.pos:len(b.data)]
	copy(b.data, v)
	b.data = b.data[:len(v)]
	b.pos = 0
}

var DefaultBufferPool BufferPool
