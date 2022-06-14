// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"github.com/xfali/restclient/v2/buffer"
	"github.com/xfali/restclient/v2/filter"
	"net/http"
	"time"
)

// 设置读写超时
func SetTimeout(timeout time.Duration) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.timeout = timeout
	}
}

// 配置初始转换器列表
func SetConverters(convs []Converter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.converters = convs
	}
}

// 添加初始转换器列表
func AddConverters(convs ...Converter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.converters = append(client.converters, convs...)
	}
}

// 配置连接池
func SetRoundTripper(tripper http.RoundTripper) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.transport = tripper
	}
}

// 配置http客户端创建器
func SetClientCreator(cliCreator HttpClientCreator) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.cliCreator = cliCreator
	}
}

// 配置是否自动添加accept
func SetAutoAccept(v AcceptFlag) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.acceptFlag = v
	}
}

// 配置是否自动添加accept
func SetResponseBodyFlag(v ResponseBodyFlag) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.respFlag = v
	}
}

// 配置内存池
func SetBufferPool(pool buffer.Pool) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.pool = pool
	}
}

// 增加处理filter
func AddFilter(filters ...filter.Filter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.filterManager.Add(filters...)
	}
}

// 增加处理filter
func AddIFilter(filters ...filter.IFilter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		for _, v := range filters {
			if v != nil {
				client.filterManager.Add(v.Filter)
			}
		}
	}
}
