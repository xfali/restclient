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

type WrapReader func(r io.Reader) io.Reader

type RequestBody struct {
	Body   interface{}
	Reader WrapReader
}

type ResponseEntity struct {
	Result     interface{}
	Header     http.Header
	StatusCode int
}

func NewRequestBody(body interface{}, r WrapReader) *RequestBody {
	return &RequestBody{
		Body:   body,
		Reader: r,
	}
}

func NewResponseEntity(result interface{}) *ResponseEntity {
	return &ResponseEntity{
		Result:     result,
		Header:     http.Header{},
		StatusCode: http.StatusOK,
	}
}

func body(ret interface{}) *RequestBody {
	if r, ok := ret.(*RequestBody); ok {
		return r
	}
	return nil
}

func entity(ret interface{}) *ResponseEntity {
	if r, ok := ret.(*ResponseEntity); ok {
		return r
	}
	return nil
}

func (e *ResponseEntity) fill(resp *http.Response) {
	e.Header = resp.Header.Clone()
	e.StatusCode = resp.StatusCode
}
