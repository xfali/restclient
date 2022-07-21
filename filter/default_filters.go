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

package filter

import (
	"bytes"
	"fmt"
	"github.com/xfali/restclient/v2/buffer"
	"github.com/xfali/restclient/v2/restutil"
	"github.com/xfali/xlog"
	"io"
	"net/http"
	"strconv"
	"time"
)

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

type DigestReader struct {
	buf bytes.Buffer
}

func (dr *DigestReader) Reader(r io.ReadCloser) io.ReadCloser {
	dr.buf.Reset()
	_, err := io.Copy(&dr.buf, r)
	r.Close()
	if err != nil {
		return nil
	}
	return buffer.NewReadCloser(dr.buf.Bytes())
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
