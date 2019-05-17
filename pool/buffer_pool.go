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
	return len(b.data)
}

func (b *ByteBuffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.data
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.data = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

func (b *ByteBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.data)
	return int64(n), err
}

func (b *ByteBuffer) Bytes() []byte {
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

func (b *ByteBuffer) Set(p []byte) {
	b.data = append(b.data[:0], p...)
}

func (b *ByteBuffer) SetString(s string) {
	b.data = append(b.data[:0], s...)
}

func (b *ByteBuffer) String() string {
	return string(b.data)
}

func (b *ByteBuffer) Reset() {
	b.data = b.data[:0]
	b.pos = 0
}

var DefaultBufferPool BufferPool