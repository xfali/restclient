// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"context"
	"net/http"
)

type Exchange func(result interface{}, url string, method string,
	params map[string]interface{}, requestBody interface{}) (int, error)

type Wrapper func(ex Exchange) Exchange

type ExchangeContext func(ctx context.Context, result interface{}, url string, method string,
	params map[string]interface{}, requestBody interface{}) (int, error)

type WrapperContext func(ex ExchangeContext) ExchangeContext

func NewWrapper(c RestClient, wrapper Wrapper) *clientWrapper {
	return &clientWrapper{c, wrapper}
}

func NewWrapperContext(c RestClient, wrapper WrapperContext) *clientWrapper {
	return &clientWrapper{c, wrapper}
}

type clientWrapper struct {
	c       RestClient
	wrapper interface{}
}

func (w *clientWrapper) AddConverter(conv Converter) {
	w.c.AddConverter(conv)
}

func (w *clientWrapper) GetConverters() []Converter {
	return w.c.GetConverters()
}

func (w *clientWrapper) Get(result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.Exchange(result, url, http.MethodGet, params, nil)
}

func (w *clientWrapper) GetContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.ExchangeContext(ctx, result, url, http.MethodGet, params, nil)
}

func (w *clientWrapper) Post(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error) {
	return w.Exchange(result, url, http.MethodPost, params, requestBody)
}

func (w *clientWrapper) PostContext(ctx context.Context, result interface{}, url string, params map[string]interface{}, body interface{}) (int, error) {
	return w.ExchangeContext(ctx, result, url, http.MethodPost, params, body)
}

func (w *clientWrapper) Put(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error) {
	return w.Exchange(result, url, http.MethodPut, params, requestBody)
}

func (w *clientWrapper) PutContext(ctx context.Context, result interface{}, url string, params map[string]interface{}, body interface{}) (int, error) {
	return w.ExchangeContext(ctx, result, url, http.MethodPut, params, body)
}

func (w *clientWrapper) Delete(result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.Exchange(result, url, http.MethodDelete, params, nil)
}

func (w *clientWrapper) DeleteContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.ExchangeContext(ctx, result, url, http.MethodDelete, params, nil)
}

func (w *clientWrapper) Head(result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.Exchange(result, url, http.MethodHead, params, nil)
}

func (w *clientWrapper) HeadContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.ExchangeContext(ctx, result, url, http.MethodHead, params, nil)
}

func (w *clientWrapper) Options(result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.Exchange(result, url, http.MethodOptions, params, nil)
}

func (w *clientWrapper) OptionsContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error) {
	return w.ExchangeContext(ctx, result, url, http.MethodOptions, params, nil)
}

func (w *clientWrapper) Patch(result interface{}, url string, params map[string]interface{}, body interface{}) (int, error) {
	return w.Exchange(result, url, http.MethodPatch, params, body)
}

func (w *clientWrapper) PatchContext(ctx context.Context, result interface{}, url string, params map[string]interface{}, body interface{}) (int, error) {
	return w.ExchangeContext(ctx, result, url, http.MethodPatch, params, body)
}

func (w *clientWrapper) Exchange(
	result interface{},
	url string,
	method string,
	params map[string]interface{},
	requestBody interface{}) (int, error) {
	return w.ExchangeContext(context.Background(), result, url, method, params, requestBody)
}

func (w *clientWrapper) ExchangeContext(
	ctx context.Context,
	result interface{},
	url string,
	method string,
	params map[string]interface{},
	requestBody interface{}) (int, error) {
	if ex, ok := w.wrapper.(WrapperContext); ok {
		return ex(w.c.ExchangeContext)(ctx, result, url, method, params, requestBody)
	} else if ex, ok := w.wrapper.(Wrapper); ok {
		return ex(w.c.Exchange)(result, url, method, params, requestBody)
	}
	panic("Wrapper is not support. ")
}
