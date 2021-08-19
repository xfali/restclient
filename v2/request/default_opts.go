// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package request

import (
	"context"
	"github.com/xfali/restclient/v2/filter"
	"net/http"
)

const (
	KeySelf             = "self.set"
	KeyMethod           = "self.method.set"
	KeyAddFilter        = "self.filter.add"
	KeyRequestHeader    = "self.request.header.set"
	KeyRequestContext   = "self.request.context.set"
	KeyRequestAddHeader = "self.request.header.add"
	KeyRequestBody      = "self.request.body.set"
	KeyResult           = "self.result.set"
	KeyResponse         = "self.response.set"
)

// 设置请求方法，请使用http包中的常量配置，如http.MethodPost
// 如果不设置默认为MethodGet
func WithMethod(method string) Opt {
	return func(setter Setter) {
		setter.Set(KeyMethod, method)
	}
}

// 设置请求体
func WithRequestBody(body interface{}) Opt {
	return func(setter Setter) {
		setter.Set(KeyRequestBody, body)
	}
}

// 设置请求的context，在filter中可以通过request.Context()获取
func WithRequestContext(ctx context.Context) Opt {
	return func(setter Setter) {
		setter.Set(KeyRequestContext, ctx)
	}
}

// 设置接收应答数据的目的对象，对象分为：
// 1、结构体指针，直接将response数据反序列化为目的对象
// 2、函数，类型为func(Type) 接收序列化完成后的对象，其中type不可为指针
func WithResult(result interface{}) Opt {
	return func(setter Setter) {
		setter.Set(KeyResult, result)
	}
}

// 获得请求的应答
// response：请求完成后会自动将应答填充到response中
// withBody：是否填充应答的body
//   如果为true则填充，!!注意：填充后调用者需手工close即response.Body.Close()，否则可能会引起内存泄漏!!
//   如果为false则不填充，Body为nil，如果读取数据会引发panic
func WithResponse(response *http.Response, withBody bool) Opt {
	return func(setter Setter) {
		setter.Set(KeyResponse, []interface{}{response, withBody})
	}
}

// 设置请求header
func WithRequestHeader(header http.Header) Opt {
	return func(setter Setter) {
		setter.Set(KeyRequestHeader, header)
	}
}

// 添加请求header
func AddRequestHeader(key, value string) Opt {
	return func(setter Setter) {
		setter.Set(KeyRequestAddHeader, []string{key, value})
	}
}

func AddIFilter(filters ...filter.IFilter) Opt {
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

func AddFilter(filters ...filter.Filter) Opt {
	return func(setter Setter) {
		setter.Set(KeyAddFilter, filters[:])
	}
}

// 设置请求方法为GET
func MethodGet() Opt {
	return WithMethod(http.MethodGet)
}

// 设置请求方法为POST
func MethodPost() Opt {
	return WithMethod(http.MethodPost)
}

// 设置请求方法为PUT
func MethodPut() Opt {
	return WithMethod(http.MethodPut)
}

// 设置请求方法为DELETE
func MethodDelete() Opt {
	return WithMethod(http.MethodDelete)
}

// 设置请求方法为HEAD
func MethodHead() Opt {
	return WithMethod(http.MethodHead)
}

// 设置请求方法为PATCH
func MethodPatch() Opt {
	return WithMethod(http.MethodPatch)
}

// 设置请求方法为GET
func MethodOptions() Opt {
	return WithMethod(http.MethodOptions)
}

// 设置请求方法为CONNECT
func MethodConnect() Opt {
	return WithMethod(http.MethodConnect)
}

// 设置请求方法为TRACE
func MethodTrace() Opt {
	return WithMethod(http.MethodTrace)
}
