// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
    "encoding/json"
)

type Serialize interface {
    Serialize(interface{}) ([]byte, error)
}

type Deserialize interface {
    Deserialize([]byte, interface{}) error
}

type Converter interface {
    Serialize
    Deserialize
    Accept() string
    ContentType() string
}

type JsonConverter struct {
}

func (c *JsonConverter) Serialize(i interface{}) ([]byte, error) {
    return json.Marshal(i)
}

func (c *JsonConverter) Deserialize(d []byte, r interface{}) error {
    return json.Unmarshal(d, r)
}

func (c *JsonConverter) Accept() string {
    return "application/json"
}

func (c *JsonConverter) ContentType() string {
    return "application/json"
}

type RestClient interface {
    Init(conv Converter, timeout int)
    Get(result interface{}, url string) (int, error)
    Post(result interface{}, url string, requestBody interface{}) (int, error)
    Put(result interface{}, url string, requestBody interface{}) (int, error)
    Delete(result interface{}, url string) (int, error)
    Exchange(result interface{}, url string, method string, header map[string]string,
        requestBody interface{}) (int, error)
}
