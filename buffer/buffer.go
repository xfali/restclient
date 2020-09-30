// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package buffer

import (
	"bytes"
	"io"
	"io/ioutil"
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

func NewReadCloser(d []byte) io.ReadCloser {
	if d == nil {
		return nil
	}
	return ioutil.NopCloser(bytes.NewReader(d))
}
