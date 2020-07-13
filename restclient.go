// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"io"
)

type Encoder interface {
	Encode(o interface{}) (int64, error)
}

type Decoder interface {
	Decode(o interface{}) (int64, error)
}

type Converter interface {
	CreateEncoder(io.Writer) Encoder
	CreateDecoder(io.Reader) Decoder

	CanEncode(o interface{}, mediaType MediaType) bool
	CanDecode(o interface{}, mediaType MediaType) bool
	SupportMediaType() []MediaType
}

type RestClient interface {
	AddConverter(conv Converter)
	GetConverters() []Converter

	Get(result interface{}, url string, params map[string]interface{}) (int, error)
	Post(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
	Put(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
	Delete(result interface{}, url string, params map[string]interface{}) (int, error)
	Head(result interface{}, url string, params map[string]interface{}) (int, error)
	Options(result interface{}, url string, params map[string]interface{}) (int, error)
	Patch(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
	Exchange(result interface{}, url string, method string, params map[string]interface{},
		requestBody interface{}) (int, error)
}
