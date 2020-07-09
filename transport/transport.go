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
	CONNECT_TIMEOUT         = 30 * time.Second
	KEEPALIVE_TIME          = 30 * time.Second
	MAX_IDLE_CONN           = 500
	MAX_IDLE_CONN_PER_HOST  = 100
	IDLE_CONN_TIMEOUT       = 90 * time.Second
	TLS_HANDSHAKE_TIMEOUT   = 10 * time.Second
	EXPECT_CONTINUE_TIMEOUT = 1 * time.Second
)

type Opt func(*http.Transport)

func New(opts ...Opt) *http.Transport {
	ret := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   CONNECT_TIMEOUT,
			KeepAlive: KEEPALIVE_TIME,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          MAX_IDLE_CONN,
		MaxIdleConnsPerHost:   MAX_IDLE_CONN_PER_HOST,
		IdleConnTimeout:       IDLE_CONN_TIMEOUT,
		TLSHandshakeTimeout:   TLS_HANDSHAKE_TIMEOUT,
		ExpectContinueTimeout: EXPECT_CONTINUE_TIMEOUT,
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
