// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package buffer

import (
	"bytes"
	"io"
	"sync"
)

const (
	InitialBufferSize = 1024
	MaxBufferSize     = 2056
)

type Pool interface {
	Get() *bytes.Buffer
	Put(*bytes.Buffer)
}

type defaultPool struct {
	initialSize int
	maxSize     int
	pool        sync.Pool
}

type Opt func(*defaultPool)

func NewPool(opts ...Opt) *defaultPool {
	ret := &defaultPool{
		initialSize: InitialBufferSize,
		maxSize:     MaxBufferSize,
	}
	for _, opt := range opts {
		opt(ret)
	}
	ret.pool = sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, ret.initialSize))
	}}
	return ret
}

func OptSetInitialBufferSize(size int) Opt {
	return func(pool *defaultPool) {
		pool.initialSize = size
	}
}

func OptSetMaxBufferSize(size int) Opt {
	return func(pool *defaultPool) {
		pool.maxSize = size
	}
}

func (p *defaultPool) Get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (p *defaultPool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	if buf.Len() > MaxBufferSize {
		return
	}
	p.pool.Put(buf)
}

type ContentLength interface {
	ContentLength() int64
}

type NopReadCloser struct {
	r   io.Reader
	len int64
}

func (rc *NopReadCloser) Read(p []byte) (n int, err error) {
	return rc.r.Read(p)
}

func (rc *NopReadCloser) Close() error {
	return nil
}

func (rc *NopReadCloser) ContentLength() int64 {
	return rc.len
}

func NewReadCloser(d []byte) *NopReadCloser {
	if d == nil {
		return nil
	}
	return &NopReadCloser{
		r:   bytes.NewReader(d),
		len: int64(len(d)),
	}
}

type ReadWriteCloser struct {
	pool Pool
	buf  *bytes.Buffer
	once sync.Once
}

func (rc *ReadWriteCloser) Bytes() []byte {
	return rc.buf.Bytes()
}

func (rc *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return rc.buf.Read(p)
}

func (rc *ReadWriteCloser) Write(p []byte) (n int, err error) {
	return rc.buf.Write(p)
}

func (rc *ReadWriteCloser) Close() error {
	// just return once
	rc.once.Do(func() {
		rc.pool.Put(rc.buf)
	})
	return nil
}

func (rc *ReadWriteCloser) ContentLength() int64 {
	return int64(rc.buf.Len())
}

func NewReadWriteCloser(pool Pool) *ReadWriteCloser {
	buf := pool.Get()
	return &ReadWriteCloser{
		pool: pool,
		buf:  buf,
	}
}

type MergeReaderWriter struct {
	r io.Reader
	w io.Writer
}

func NewMergeReaderWriter(r io.Reader, w io.Writer) *MergeReaderWriter {
	return &MergeReaderWriter{
		r: r,
		w: w,
	}
}

func (mrw *MergeReaderWriter) Read(p []byte) (int, error) {
	n, err := mrw.r.Read(p)
	_, _ = mrw.w.Write(p[:n])
	return n, err
}
