// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"fmt"
	"github.com/xfali/restclient/buffer"
	"github.com/xfali/restclient/restutil"
	"io"
	"net/http"
	"time"
)

func NewBasicAuthClient(client RestClient, auth *BasicAuth) RestClient {
	return NewWrapper(client, auth.Exchange)
}

func (auth *BasicAuth) ResetCredentials(username, password string) {
	auth.lock.Lock()
	defer auth.lock.Unlock()

	auth.username = username
	auth.password = password
}

func (auth *BasicAuth) Filter(request *http.Request, fc FilterChain) (*http.Response, error) {
	auth.lock.Lock()
	k, v := restutil.BasicAuthHeader(auth.username, auth.password)
	request.Header.Set(k, v)
	auth.lock.Unlock()

	return fc.Filter(request)
}

func (auth *BasicAuth) Exchange(ex Exchange) Exchange {
	return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		if params == nil {
			params = map[string]interface{}{}
		}
		auth.lock.Lock()
		k, v := restutil.BasicAuthHeader(auth.username, auth.password)
		auth.lock.Unlock()

		params[k] = v
		n, err := ex(result, url, method, params, requestBody)
		return n, err
	}
}

func NewAccessTokenAuthClient(client RestClient, auth *AccessTokenAuth) RestClient {
	return NewWrapper(client, auth.Exchange)
}

func (auth *AccessTokenAuth) ResetCredentials(token string) {
	auth.lock.Lock()
	defer auth.lock.Unlock()

	auth.token = token
}

func (auth *AccessTokenAuth) Filter(request *http.Request, fc FilterChain) (*http.Response, error) {
	auth.lock.Lock()
	k, v := auth.tokenBuilder(auth.token)
	request.Header.Set(k, v)
	auth.lock.Unlock()

	return fc.Filter(request)
}

func (b *AccessTokenAuth) Exchange(ex Exchange) Exchange {
	return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		if params == nil {
			params = map[string]interface{}{}
		}
		b.lock.Lock()
		k, v := b.tokenBuilder(b.token)
		b.lock.Unlock()

		params[k] = v

		n, err := ex(result, url, method, params, requestBody)
		return n, err
	}
}

func (auth *DigestAuth) Filter(request *http.Request, fc FilterChain) (*http.Response, error) {
	buf := auth.pool.Get()
	defer auth.pool.Put(buf)

	var reqData []byte
	if request.Body != nil {
		_, err := io.Copy(buf, request.Body)
		if err != nil {
			return nil, err
		}
		reqData = buf.Bytes()
		request.Body = buffer.NewReadCloser(reqData)
	}

	resp, err := fc.Filter(request)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		da := auth.newDigestData()
		digest := findWWWAuth(resp.Header)
		wwwAuth := ParseWWWAuthenticate(digest)
		err := da.Refresh(request.Method, request.RequestURI, reqData, wwwAuth)
		if err != nil {
			return nil, err
		}
		auth, err := da.ToString()
		if err != nil {
			return nil, err
		}
		request.Header.Set(restutil.HeaderAuthorization, auth)
		if reqData != nil {
			request.Body = buffer.NewReadCloser(reqData)
		}
		return fc.Filter(request)
	}
	return resp, err
}

func findWWWAuth(header http.Header) string {
	if header == nil {
		return ""
	}

	digest := header.Get("WWW-Authenticate")
	if digest != "" {
		return digest
	}

	digest = header.Get("Www-Authenticate")
	if digest != "" {
		return digest
	}

	digest = header.Get("www-Authenticate")
	if digest != "" {
		return digest
	}

	return ""
}

type LogFunc func(format string, args ...interface{})
type Log struct {
	Log  LogFunc
	Tag  string
	pool buffer.Pool
}

func NewLog(log LogFunc, tag string) *Log {
	if tag == "" {
		tag = "restclient"
	}
	return &Log{
		Log:  log,
		Tag:  tag,
		pool: buffer.NewPool(),
	}
}

func (log *Log) Filter(request *http.Request, fc FilterChain) (*http.Response, error) {
	buf := log.pool.Get()
	defer log.pool.Put(buf)

	var reqData []byte
	if request.Body != nil {
		_, err := io.Copy(buf, request.Body)
		if err != nil {
			return nil, err
		}
		reqData = buf.Bytes()
		request.Body = buffer.NewReadCloser(reqData)
	}

	now := time.Now()
	id := RandomId(10)
	log.Log("[%s request %s]: url: %s , method: %s , header: %v , body: %s \n",
		log.Tag, id, request.URL.String(), request.Method, request.Header, string(reqData))

	resp, err := fc.Filter(request)
	var (
		status   int
		header   http.Header
		respData []byte
	)
	if resp != nil {
		status = resp.StatusCode
		header = resp.Header

		if resp.Body != nil {
			respBuf := log.pool.Get()
			defer log.pool.Put(respBuf)

			_, rspErr := io.Copy(respBuf, resp.Body)
			resp.Body.Close()
			if rspErr == nil {
				respData = respBuf.Bytes()
			}
			resp.Body = buffer.NewReadCloser(respData)
		}
	}
	if err != nil {
		log.Log("[%s response %s]: use time: %d ms, status: %d , header: %v, result: %s, error: %v \n",
			log.Tag, id, time.Since(now)/time.Millisecond, status, header, string(respData), err)
	} else {
		log.Log("[%s response %s]: use time: %d ms, status: %d , header: %v, result: %s \n",
			log.Tag, id, time.Since(now)/time.Millisecond, status, header, string(respData))
	}

	return resp, err
}

type RecoveryFilter struct {
	Log LogFunc
}

func NewRecovery(log LogFunc) *RecoveryFilter {
	return &RecoveryFilter{
		Log: log,
	}
}

func (rf *RecoveryFilter) Filter(request *http.Request, fc FilterChain) (resp *http.Response, err error) {
	defer func() {
		r := recover()
		if r != nil && rf.Log != nil {
			rf.Log("%v\n", r)
		}
		err = fmt.Errorf("RestClient panic :%v\n", r)
	}()
	return fc.Filter(request)
}

type Builder struct {
	c RestClient
}

func (b *Builder) Default(opts ...Opt) *Builder {
	b.c = New(opts...)
	return b
}

func (b *Builder) BasicAuth(auth *BasicAuth) *Builder {
	b.c = NewBasicAuthClient(b.c, auth)
	return b
}

func (b *Builder) Build() RestClient {
	return b.c
}
