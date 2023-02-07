/*
 * Copyright 2022 Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
		header: make(http.Header),
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
	case request.KeyRequestAddCookie:
		p.addCookies(value.([]*http.Cookie))
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

func (p *defaultParam) addCookies(cookies []*http.Cookie) {
	if len(cookies) == 0 {
		return
	}
	req := &http.Request{
		Header: make(http.Header),
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}

	for k, vs := range req.Header {
		for _, v := range vs {
			p.header.Add(k, v)
		}
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

func (p *defaultParam) AddCookies(cookies ...*http.Cookie) *defaultParam {
	p.addCookies(cookies)
	return p
}

func (p *defaultParam) AddHeaders(key string, values ...string) *defaultParam {
	for _, v := range values {
		p.header.Add(key, v)
	}
	return p
}

func (p *defaultParam) SetHeader(key string, value string) *defaultParam {
	p.header.Set(key, value)
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

func NewUrlBuilder(url string) *restutil.UrlBuilder {
	return restutil.NewUrlBuilder(url)
}
