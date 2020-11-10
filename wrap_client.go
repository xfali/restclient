// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"bytes"
	"fmt"
	"github.com/xfali/restclient/buffer"
	"github.com/xfali/restclient/restutil"
	"github.com/xfali/xlog"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
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

func NewDigestAuthClient(client RestClient, auth *DigestAuth) RestClient {
	return NewWrapper(client, auth.Exchange)
}

type DigestReader struct {
	buf bytes.Buffer
}

func (dr *DigestReader) Reader(r io.ReadCloser) io.ReadCloser {
	dr.buf.Reset()
	_, err := io.Copy(&dr.buf, r)
	if err != nil {
		return nil
	}
	return buffer.NewReadCloser(dr.buf.Bytes())
}

func (b *DigestAuth) Exchange(ex Exchange) Exchange {
	return func(result interface{}, uri string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		ent := responseEntity(result)
		if ent == nil {
			ent = NewResponseEntity(result)
		}
		digestBuf := DigestReader{}
		if requestBody != nil {
			body := requestEntity(requestBody)
			if body == nil {
				body = NewRequestEntity(requestBody, digestBuf.Reader)
			} else {
				originReader := body.Reader
				body.Reader = func(r io.ReadCloser) io.ReadCloser {
					return originReader(digestBuf.Reader(r))
				}
			}
			requestBody = body
		}
		n, err := ex(ent, uri, method, params, requestBody)
		if n == http.StatusUnauthorized {
			da := b.newDigestData()
			digest := findWWWAuth(ent.Header)
			wwwAuth := ParseWWWAuthenticate(digest)
			uriP, _ := url.Parse(uri)
			err := da.Refresh(method, uriP.RequestURI(), digestBuf.buf.Bytes(), wwwAuth)
			if err != nil {
				return n, err
			}
			auth, err := da.ToString()
			if err != nil {
				return n, err
			}
			if params == nil {
				params = map[string]interface{}{}
			}
			params[restutil.HeaderAuthorization] = auth
			return ex(result, uri, method, params, requestBody)
		}
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
		// close old request body
		request.Body.Close()
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

func ContentLengthFilter(request *http.Request, fc FilterChain) (*http.Response, error) {
	if request.ContentLength > 0 {
		return fc.Filter(request)
	}
	lengthStr := request.Header.Get("Content-Length")
	if lengthStr != "" {
		l, err := strconv.ParseInt(lengthStr, 10, 64)
		if err == nil {
			request.ContentLength = l
		}
		return fc.Filter(request)
	}
	if cl, ok := request.Body.(buffer.ContentLength); ok {
		request.ContentLength = cl.ContentLength()
	}
	return fc.Filter(request)
}

type LogFunc func(format string, args ...interface{})
type Log struct {
	Log  xlog.Logger
	Tag  string
	pool buffer.Pool
}

func NewLog(log xlog.Logger, tag string) *Log {
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
	reqBuf := buffer.NewReadWriteCloser(log.pool)

	var reqData []byte
	if request.Body != nil {
		_, err := io.Copy(reqBuf, request.Body)
		if err != nil {
			return nil, err
		}
		reqData = reqBuf.Bytes()
		// close old request body
		request.Body.Close()
		request.Body = reqBuf
	}

	now := time.Now()
	id := RandomId(10)
	log.Log.Infof("[%s request %s]: url: %s , method: %s , header: %v , body: %s \n",
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
			respBuf := buffer.NewReadWriteCloser(log.pool)

			_, rspErr := io.Copy(respBuf, resp.Body)
			resp.Body.Close()
			if rspErr == nil {
				respData = respBuf.Bytes()
			}
			resp.Body = respBuf
		}
	}
	if err != nil {
		log.Log.Infof("[%s response %s]: use time: %d ms, status: %d , header: %v, result: %s, error: %v \n",
			log.Tag, id, time.Since(now)/time.Millisecond, status, header, string(respData), err)
	} else {
		log.Log.Infof("[%s response %s]: use time: %d ms, status: %d , header: %v, result: %s \n",
			log.Tag, id, time.Since(now)/time.Millisecond, status, header, string(respData))
	}

	return resp, err
}

func NewLogClient(client RestClient, log *Log) RestClient {
	return NewWrapper(client, log.Exchange)
}

func (log *Log) Exchange(ex Exchange) Exchange {
	return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		now := time.Now()
		id := RandomId(10)
		log.Log.Infof("[%s request %s]: url: %v , method: %v , params: %v , body: %v \n",
			log.Tag, id, url, method, params, requestBody)
		n, err := ex(result, url, method, params, requestBody)
		entity := responseEntity(result)
		if entity != nil {
			result = entity.Result
			v := reflect.ValueOf(result)
			v = reflect.Indirect(v)
			log.Log.Infof("[%s response %s]: use time: %d ms, status: %d , header: %v, result: %v ",
				log.Tag, id, time.Since(now)/time.Millisecond, n, entity.Header, v.Interface())
		} else {
			v := reflect.ValueOf(result)
			v = reflect.Indirect(v)
			log.Log.Infof("[%s response %s]: use time: %d ms, status: %d , result: %v ",
				log.Tag, id, time.Since(now)/time.Millisecond, n, v.Interface())
		}
		return n, err
	}
}

type RecoveryFilter struct {
	Log xlog.Logger
}

func NewRecovery(log xlog.Logger) *RecoveryFilter {
	return &RecoveryFilter{
		Log: log,
	}
}

func (rf *RecoveryFilter) Filter(request *http.Request, fc FilterChain) (resp *http.Response, err error) {
	defer func() {
		r := recover()
		if r != nil && rf.Log != nil {
			rf.Log.Infoln(r)
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
