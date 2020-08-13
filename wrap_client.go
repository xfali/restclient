// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"bytes"
	"github.com/xfali/restclient/restutil"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

func NewBasicAuthClient(client RestClient, auth *BasicAuth) RestClient {
	return NewWrapper(client, auth.Exchange)
}

func (b *BasicAuth) Exchange(ex Exchange) Exchange {
	return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		if params == nil {
			params = map[string]interface{}{}
		}
		k, v := restutil.BasicAuthHeader(b.Username, b.Password)
		params[k] = v
		n, err := ex(result, url, method, params, requestBody)
		return n, err
	}
}

func NewAccessTokenAuthClient(client RestClient, auth *AccessTokenAuth) RestClient {
	return NewWrapper(client, auth.Exchange)
}

func (b *AccessTokenAuth) Exchange(ex Exchange) Exchange {
	return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		if params == nil {
			params = map[string]interface{}{}
		}
		if b.Type == restutil.Bearer {
			k, v := restutil.AccessTokenAuthHeader(b.Token)
			params[k] = v
		} else {
			params[b.Name] = b.Token
		}

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

func (dr *DigestReader) Reader(r io.Reader) io.Reader {
	_, err := io.Copy(&dr.buf, r)
	if err != nil {
		return nil
	}
	return bytes.NewReader(dr.buf.Bytes())
}

func (b *DigestAuth) Exchange(ex Exchange) Exchange {
	return func(result interface{}, uri string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		ent := entity(result)
		if ent == nil {
			ent = NewResponseEntity(result)
		}
		digestBuf := DigestReader{}
		if requestBody != nil {
			body := body(requestBody)
			if body == nil {
				body = NewRequestBody(requestBody, digestBuf.Reader)
			} else {
				originReader := body.Reader
				body.Reader = func(r io.Reader) io.Reader {
					return originReader(digestBuf.Reader(r))
				}
			}
			requestBody = body
		}
		n, err := ex(ent, uri, method, params, requestBody)
		if n == http.StatusUnauthorized {
			digest := findWWWAuth(ent.Header)
			wwwAuth := ParseWWWAuthenticate(digest)
			uriP, _ := url.Parse(uri)
			err := b.Refresh(method, uriP.RequestURI(), digestBuf.buf.Bytes(), wwwAuth)
			if err != nil {
				return n, err
			}
			auth, err := b.ToString()
			if err != nil {
				return n, err
			}
			if params == nil {
				params = map[string]interface{}{}
			}
			params["Authorization"] = auth
			return ex(result, uri, method, params, requestBody)
		}
		return n, err
	}
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
	Log LogFunc
	Tag string
}

func NewLog(log LogFunc, tag string) *Log {
	if tag == "" {
		tag = "restclient"
	}
	return &Log{
		Log: log,
		Tag: tag,
	}
}

func NewLogClient(client RestClient, log *Log) RestClient {
	return NewWrapper(client, log.Exchange)
}

func (log *Log) Exchange(ex Exchange) Exchange {
	return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
		now := time.Now()
		id := RandomId(10)
		log.Log("[%s request %s]: url: %v , method: %v , params: %v , body: %v \n",
			log.Tag, id, url, method, params, requestBody)
		n, err := ex(result, url, method, params, requestBody)
		entity := entity(result)
		if entity != nil {
			result = entity.Result
			v := reflect.ValueOf(result)
			v = reflect.Indirect(v)
			log.Log("[%s response %s]: use time: %d ms, status: %d , header: %v, result: %v ",
				log.Tag, id, time.Since(now)/time.Millisecond, n, entity.Header, v.Interface())
		} else {
			v := reflect.ValueOf(result)
			v = reflect.Indirect(v)
			log.Log("[%s response %s]: use time: %d ms, status: %d , result: %v ",
				log.Tag, id, time.Since(now)/time.Millisecond, n, v.Interface())
		}
		return n, err
	}
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

func (b *Builder) DigestAuth(auth *DigestAuth) *Builder {
	b.c = NewDigestAuthClient(b.c, auth)
	return b
}

func (b *Builder) Log(log *Log) *Builder {
	b.c = NewLogClient(b.c, log)
	return b
}

func (b *Builder) Build() RestClient {
	return b.c
}
