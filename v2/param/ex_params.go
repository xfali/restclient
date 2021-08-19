// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package param

import (
	"context"
	"github.com/xfali/restclient/v2/filter"
	"net/http"
)

const (
	KeyMethod           = "method"
	KeyAddFilter        = "add_filter"
	KeyRequestHeader    = "request_header"
	KeyRequestContext   = "request_context"
	KeyRequestAddHeader = "request_add_header"
	KeyRequestBody      = "request_body"
	KeyResult           = "result"
	KeyResponse         = "response"
)

// 设置请求方法，请使用http包中的常量配置，如http.MethodPost
// 如果不设置默认为MethodGet
func Method(method string) Parameter {
	return func(setter Setter) {
		setter.Set(KeyMethod, method)
	}
}

// 设置请求体
func RequestBody(body interface{}) Parameter {
	return func(setter Setter) {
		setter.Set(KeyRequestBody, body)
	}
}

// 设置请求的context，在filter中可以通过request.Context()获取
func RequestContext(ctx context.Context) Parameter {
	return func(setter Setter) {
		setter.Set(KeyRequestContext, ctx)
	}
}

// 设置接收应答数据的目的对象，对象分为：
// 1、结构体指针，直接将response数据反序列化为目的对象
// 2、函数，类型为func(Type) 接收序列化完成后的对象，其中type不可为指针
func Result(result interface{}) Parameter {
	return func(setter Setter) {
		setter.Set(KeyResult, result)
	}
}

// 获得请求的应答
// response：请求完成后会自动将应答填充到response中
// withBody：是否填充应答的body
//   如果为true则填充，!!注意：填充后调用者需手工close即response.Body.Close()，否则可能会引起内存泄漏!!
//   如果为false则不填充，Body为nil，如果读取数据会引发panic
func ResponseReceiver(response *http.Response, withBody bool) Parameter {
	return func(setter Setter) {
		setter.Set(KeyResponse, []interface{}{response, withBody})
	}
}

// 设置请求header
func RequestHeader(header http.Header) Parameter {
	return func(setter Setter) {
		setter.Set(KeyRequestHeader, header)
	}
}

// 添加请求header
func AddRequestHeader(key, value string) Parameter {
	return func(setter Setter) {
		setter.Set(KeyRequestAddHeader, []string{key, value})
	}
}

func AddIFilter(filters ...filter.IFilter) Parameter {
	return func(setter Setter) {
		fs := make([]filter.Filter, 0, len(filters))
		for i, v := range filters {
			if v != nil {
				fs[i] = v.Filter
			}
		}
		setter.Set(KeyAddFilter, fs)
	}
}

func AddFilter(filters ...filter.Filter) Parameter {
	return func(setter Setter) {
		setter.Set(KeyAddFilter, filters[:])
	}
}

// 设置请求方法为GET
func MethodGet() Parameter {
	return Method(http.MethodGet)
}

// 设置请求方法为POST
func MethodPost() Parameter {
	return Method(http.MethodPost)
}

// 设置请求方法为PUT
func MethodPut() Parameter {
	return Method(http.MethodPut)
}

// 设置请求方法为DELETE
func MethodDelete() Parameter {
	return Method(http.MethodDelete)
}

// 设置请求方法为HEAD
func MethodHead() Parameter {
	return Method(http.MethodHead)
}

// 设置请求方法为PATCH
func MethodPatch() Parameter {
	return Method(http.MethodPatch)
}

// 设置请求方法为GET
func MethodOptions() Parameter {
	return Method(http.MethodOptions)
}

// 设置请求方法为CONNECT
func MethodConnect() Parameter {
	return Method(http.MethodConnect)
}

// 设置请求方法为TRACE
func MethodTrace() Parameter {
	return Method(http.MethodTrace)
}
