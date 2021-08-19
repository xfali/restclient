// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

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
