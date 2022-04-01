// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"context"
	"github.com/xfali/restclient/v2/filter"
	"github.com/xfali/restclient/v2/request"
	"github.com/xfali/restclient/v2/restutil"
	"net/http"
)

type defaultParam struct {
	ctx           context.Context
	method        string
	header        http.Header
	filterManager filter.FilterManager

	reqBody  interface{}
	result   interface{}
	response *http.Response
	respFlag bool
}

func emptyParam() *defaultParam {
	return &defaultParam{
		method: http.MethodGet,
		ctx:    context.Background(),
		header: http.Header{},
	}
}

func (p *defaultParam) Set(key string, value interface{}) {
	switch key {
	case request.KeyMethod:
		p.method = value.(string)
	case request.KeyAddFilter:
		p.filterManager.Add(value.([]filter.Filter)...)
	case request.KeyRequestContext:
		p.ctx = value.(context.Context)
	case request.KeyRequestHeader:
		p.header = value.(http.Header)
	case request.KeyRequestAddHeader:
		ss := value.([]string)
		p.header.Add(ss[0], ss[1])
	case request.KeyRequestBody:
		p.reqBody = value
	case request.KeyResult:
		p.result = value
	case request.KeyResponse:
		rs := value.([]interface{})
		p.response = rs[0].(*http.Response)
		p.respFlag = rs[1].(bool)
	}
}

func NewRequest() *defaultParam {
	return emptyParam()
}

func (p *defaultParam) Context(ctx context.Context) *defaultParam {
	p.ctx = ctx
	return p
}

func (p *defaultParam) Method(method string) *defaultParam {
	p.method = method
	return p
}

// 设置请求方法为GET
func (p *defaultParam) MethodGet() *defaultParam {
	return p.Method(http.MethodGet)
}

// 设置请求方法为POST
func (p *defaultParam) MethodPost() *defaultParam {
	return p.Method(http.MethodPost)
}

// 设置请求方法为PUT
func (p *defaultParam) MethodPut() *defaultParam {
	return p.Method(http.MethodPut)
}

// 设置请求方法为DELETE
func (p *defaultParam) MethodDelete() *defaultParam {
	return p.Method(http.MethodDelete)
}

// 设置请求方法为HEAD
func (p *defaultParam) MethodHead() *defaultParam {
	return p.Method(http.MethodHead)
}

// 设置请求方法为PATCH
func (p *defaultParam) MethodPatch() *defaultParam {
	return p.Method(http.MethodPatch)
}

// 设置请求方法为GET
func (p *defaultParam) MethodOptions() *defaultParam {
	return p.Method(http.MethodOptions)
}

// 设置请求方法为CONNECT
func (p *defaultParam) MethodConnect() *defaultParam {
	return p.Method(http.MethodConnect)
}

// 设置请求方法为TRACE
func (p *defaultParam) MethodTrace() *defaultParam {
	return p.Method(http.MethodTrace)
}

func (p *defaultParam) Header(header http.Header) *defaultParam {
	p.header = header
	return p
}

func (p *defaultParam) Accept(accept string) *defaultParam {
	p.header.Add(restutil.HeaderAccept, accept)
	return p
}

func (p *defaultParam) ContentType(contentType string) *defaultParam {
	p.header.Add(restutil.HeaderContentType, contentType)
	return p
}

func (p *defaultParam) RequestBody(reqBody interface{}) *defaultParam {
	p.reqBody = reqBody
	return p
}

func (p *defaultParam) Result(result interface{}) *defaultParam {
	p.result = result
	return p
}

func (p *defaultParam) Response(response *http.Response, withResponseBody bool) *defaultParam {
	p.response = response
	p.respFlag = withResponseBody
	return p
}

func (p *defaultParam) Filters(filters ...filter.Filter) *defaultParam {
	p.filterManager.Add(filters...)
	return p
}

func (p *defaultParam) self(setter request.Setter) {
	if s, ok := setter.(*defaultParam); ok {
		*s = *p
	}
}

func (p *defaultParam) Build() request.Opt {
	return p.self
}
