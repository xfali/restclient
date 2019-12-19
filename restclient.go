// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
    "bytes"
    "encoding/json"
    "io"
)

type Serialize interface {
    Serialize(interface{}) (io.Reader, error)
}

type Deserialize interface {
    Deserialize(io.Reader, interface{}) (int, error)
}

type Converter interface {
    Serialize
    Deserialize
    Accept() string
    ContentType() string
}

type JsonConverter struct {
}

func (c *JsonConverter) Serialize(i interface{}) (io.Reader, error) {
    d, err := json.Marshal(i)
    if err != nil {
        return nil, err
    }
    return bytes.NewReader(d), nil
}

func (c *JsonConverter) Deserialize(r io.Reader, result interface{}) (int, error) {
    buf := bytes.NewBuffer(nil)
    n, err := io.Copy(buf, r)
    if err != nil {
        return int(n), err
    }

    d := buf.Bytes()
    return int(n), json.Unmarshal(d, result)
}

func (c *JsonConverter) Accept() string {
    return "application/json"
}

func (c *JsonConverter) ContentType() string {
    return "application/json"
}

type RestClient interface {
    AddConverter(conv Converter)
    Get(result interface{}, url string) (int, error)
    Post(result interface{}, url string, requestBody interface{}) (int, error)
    Put(result interface{}, url string, requestBody interface{}) (int, error)
    Delete(result interface{}, url string) (int, error)
    Exchange(result interface{}, url string, method string, header map[string]string,
        requestBody interface{}) (int, error)
}
