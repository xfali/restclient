// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import "net/http"

type Exchange func(result interface{}, url string, method string,
    params map[string]interface{}, requestBody interface{}) (int, error)

type Wrapper func(ex Exchange) Exchange

func NewWrapper(c RestClient, wrapper Wrapper) RestClient {
    return &ClientWrapper{c, wrapper}
}

type ClientWrapper struct {
    c       RestClient
    wrapper Wrapper
}

func (w *ClientWrapper) AddConverter(conv Converter) {
    w.c.AddConverter(conv)
}

func (w *ClientWrapper) Get(result interface{}, url string, params map[string]interface{}) (int, error) {
    return w.Exchange(result, url, http.MethodGet, params, nil)
}

func (w *ClientWrapper) Post(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error) {
    return w.Exchange(result, url, http.MethodPost, params, requestBody)
}

func (w *ClientWrapper) Put(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error) {
    return w.Exchange(result, url, http.MethodPut, params, requestBody)
}

func (w *ClientWrapper) Delete(result interface{}, url string, params map[string]interface{}) (int, error) {
    return w.Exchange(result, url, http.MethodDelete, params, nil)
}

func (w *ClientWrapper) Head(result interface{}, url string, params map[string]interface{}) (int, error) {
    return w.Exchange(result, url, http.MethodHead, params, nil)
}

func (w *ClientWrapper) Options(result interface{}, url string, params map[string]interface{}) (int, error) {
    return w.Exchange(result, url, http.MethodOptions, params, nil)
}

func (w *ClientWrapper) Patch(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error) {
    return w.Exchange(result, url, http.MethodPatch, params, requestBody)
}

func (w *ClientWrapper) Exchange(
    result interface{},
    url string,
    method string,
    params map[string]interface{},
    requestBody interface{}) (int, error) {
    return w.wrapper(w.c.Exchange)(result, url, method, params, requestBody)
}
