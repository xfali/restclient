// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
    "io"
)

type Serializer interface {
    Serialize(o interface{}) (io.Reader, error)
    CanSerialize(o interface{}, mediaType MediaType) bool
}

type Deserializer interface {
    Deserialize(io.Reader, interface{}) (int, error)
    CanDeserialize(o interface{}, mediaType MediaType) bool
}

type Converter interface {
    Serializer
    Deserializer
    SupportMediaType() []MediaType
}

type RestClient interface {
    AddConverter(conv Converter)
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
