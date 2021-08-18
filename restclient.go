// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"context"
)

type RestClient interface {
	AddConverter(conv Converter)
	GetConverters() []Converter

	Get(result interface{}, url string, params map[string]interface{}) (int, error)
	GetContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error)

	Post(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
	PostContext(ctx context.Context, result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)

	Put(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
	PutContext(ctx context.Context, result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)

	Delete(result interface{}, url string, params map[string]interface{}) (int, error)
	DeleteContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error)

	Head(result interface{}, url string, params map[string]interface{}) (int, error)
	HeadContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error)

	Options(result interface{}, url string, params map[string]interface{}) (int, error)
	OptionsContext(ctx context.Context, result interface{}, url string, params map[string]interface{}) (int, error)

	Patch(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
	PatchContext(ctx context.Context, result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)

	Exchange(result interface{}, url string, method string, params map[string]interface{},
		requestBody interface{}) (int, error)
	ExchangeContext(ctx context.Context, result interface{}, url string, method string, params map[string]interface{},
		requestBody interface{}) (int, error)
}
