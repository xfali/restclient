// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"io"
	"net/http"
)

type WrapReader func(r io.ReadCloser) io.ReadCloser

type RequestEntity struct {
	Body   interface{}
	Reader WrapReader
}

func NewRequestEntity(body interface{}, r WrapReader) *RequestEntity {
	return &RequestEntity{
		Body:   body,
		Reader: r,
	}
}

type ResponseEntity struct {
	Result     interface{}
	Header     http.Header
	StatusCode int
}

func NewResponseEntity(result interface{}) *ResponseEntity {
	return &ResponseEntity{
		Result:     result,
		Header:     http.Header{},
		StatusCode: http.StatusOK,
	}
}

func requestEntity(ret interface{}) *RequestEntity {
	if r, ok := ret.(*RequestEntity); ok {
		return r
	}
	return nil
}

func responseEntity(ret interface{}) *ResponseEntity {
	if r, ok := ret.(*ResponseEntity); ok {
		return r
	}
	return nil
}

func (e *ResponseEntity) fill(resp *http.Response) {
	e.Header = resp.Header.Clone()
	e.StatusCode = resp.StatusCode
}
