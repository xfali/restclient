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

package transport

import (
	"net"
	"net/http"
	"time"
)

const (
	ConnectTimeout        = 30 * time.Second
	KeepaliveTime         = 30 * time.Second
	MaxIdleConn           = 100
	MaxIdleConnPerHost    = 5
	IdleConnTimeout       = 90 * time.Second
	TlsHandshakeTimeout   = 10 * time.Second
	ExpectContinueTimeout = 1 * time.Second
)

type Opt func(*http.Transport)

func New(opts ...Opt) *http.Transport {
	ret := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   ConnectTimeout,
			KeepAlive: KeepaliveTime,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          MaxIdleConn,
		MaxIdleConnsPerHost:   MaxIdleConnPerHost,
		IdleConnTimeout:       IdleConnTimeout,
		TLSHandshakeTimeout:   TlsHandshakeTimeout,
		ExpectContinueTimeout: ExpectContinueTimeout,
	}
	for i := range opts {
		opts[i](ret)
	}
	return ret
}

func SetDialContext(connTimeout, keepAlive time.Duration) Opt {
	return func(transport *http.Transport) {
		transport.DialContext = (&net.Dialer{
			Timeout:   connTimeout,
			KeepAlive: keepAlive,
			DualStack: true,
		}).DialContext
	}
}

func SetMaxIdleConnects(size int) Opt {
	return func(transport *http.Transport) {
		transport.MaxIdleConns = size
	}
}

func SetMaxIdleConnectsPerHost(size int) Opt {
	return func(transport *http.Transport) {
		transport.MaxIdleConnsPerHost = size
	}
}

func SetIdleConnectTimeout(time time.Duration) Opt {
	return func(transport *http.Transport) {
		transport.IdleConnTimeout = time
	}
}

func SetTlsShakeTimeout(time time.Duration) Opt {
	return func(transport *http.Transport) {
		transport.TLSHandshakeTimeout = time
	}
}

func SetExpectContinueTimeout(time time.Duration) Opt {
	return func(transport *http.Transport) {
		transport.ExpectContinueTimeout = time
	}
}
