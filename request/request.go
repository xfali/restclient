// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package request

import (
    "github.com/xfali/restclient"
    "net/http"
)

var defaultClient = restclient.New()

type Request struct {
    C restclient.RestClient

    params map[string]interface{}
    body   interface{}
    method string
}

type Opt func(req *Request)

func NewRequest(opts ...Opt) *Request {
    ret := &Request{
        C:      defaultClient,
        method: http.MethodGet,
    }
    for _, opt := range opts {
        opt(ret)
    }
    return ret
}

func SetClient(client restclient.RestClient) Opt {
    return func(req *Request) {
        req.C = client
    }
}

func SetHeaders(headers map[string]string) Opt {
    return func(req *Request) {
        for k, v := range headers {
            req.params[k] = v
        }
    }
}

func SetBody(body interface{}) Opt {
    return func(req *Request) {
        req.body = body
    }
}

func SetMethod(method string) Opt {
    return func(req *Request) {
        req.method = method
    }
}

func (req *Request) Exchange(url string, result interface{}) (int, error) {
    return req.C.Exchange(result, url, req.method, req.params, req.body)
}

func (req *Request) Get(url string, result interface{}) (int, error) {
    req.method = http.MethodGet
    return req.Exchange(url, result)
}

func (req *Request) Post(url string, result interface{}) (int, error) {
    req.method = http.MethodPost
    return req.Exchange(url, result)
}

func (req *Request) Put(url string, result interface{}) (int, error) {
    req.method = http.MethodPut
    return req.Exchange(url, result)
}

func (req *Request) Delete(url string, result interface{}) (int, error) {
    req.method = http.MethodDelete
    return req.Exchange(url, result)
}

func (req *Request) Head(url string, result interface{}) (int, error) {
    req.method = http.MethodHead
    return req.Exchange(url, result)
}

func (req *Request) Options(url string, result interface{}) (int, error) {
    req.method = http.MethodOptions
    return req.Exchange(url, result)
}

func (req *Request) Patch(url string, result interface{}) (int, error) {
    req.method = http.MethodPatch
    return req.Exchange(url, result)
}
