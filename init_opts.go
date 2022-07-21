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
	"github.com/xfali/restclient/v2/buffer"
	"github.com/xfali/restclient/v2/filter"
	"net/http"
	"time"
)

// SetTimeout 设置读写超时
func SetTimeout(timeout time.Duration) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.timeout = timeout
	}
}

// SetConverters 配置初始转换器列表
func SetConverters(convs []Converter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.converters = convs
	}
}

// AddConverters 添加初始转换器列表
func AddConverters(convs ...Converter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.converters = append(client.converters, convs...)
	}
}

// SetRoundTripper 配置连接池
func SetRoundTripper(tripper http.RoundTripper) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.transport = tripper
	}
}

// SetClientCreator 配置http客户端创建器
func SetClientCreator(cliCreator HttpClientCreator) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.cliCreator = cliCreator
	}
}

// CookieJar 配置http.Client的CookieJar
func CookieJar(jar http.CookieJar) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.jar = jar
	}
}

// SetAutoAccept 配置是否自动添加accept
func SetAutoAccept(v AcceptFlag) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.acceptFlag = v
	}
}

// SetResponseBodyFlag 配置是否自动添加accept
func SetResponseBodyFlag(v ResponseBodyFlag) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.respFlag = v
	}
}

// SetBufferPool 配置内存池
func SetBufferPool(pool buffer.Pool) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.pool = pool
	}
}

// AddFilter 增加处理filter
func AddFilter(filters ...filter.Filter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		client.filterManager.Add(filters...)
	}
}

// AddIFilter 增加处理filter
func AddIFilter(filters ...filter.IFilter) func(client *defaultRestClient) {
	return func(client *defaultRestClient) {
		for _, v := range filters {
			if v != nil {
				client.filterManager.Add(v.Filter)
			}
		}
	}
}
