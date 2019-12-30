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

type Accept interface {
    Accept() string
}

type ContentType interface {
    ContentType() string
}

type Converter struct {
    Serialize   func(interface{}) (io.Reader, error)
    Deserialize func(io.Reader, interface{}) (int, error)
    Accept      func() string
    ContentType func() string
}

var JsonConverter = Converter{
    Serialize:   JsonSerialize,
    Deserialize: JsonDeserialize,
    Accept:      JsonAccept,
    ContentType: JsonContentType,
}

func JsonSerialize(i interface{}) (io.Reader, error) {
    d, err := json.Marshal(i)
    if err != nil {
        return nil, err
    }
    return bytes.NewReader(d), nil
}

func JsonDeserialize(r io.Reader, result interface{}) (int, error) {
    buf := bytes.NewBuffer(nil)
    n, err := io.Copy(buf, r)
    if err != nil {
        return int(n), err
    }

    d := buf.Bytes()
    return int(n), json.Unmarshal(d, result)
}

func JsonAccept() string {
    return "application/json"
}

func JsonContentType() string {
    return "application/json"
}

type RestClient interface {
    AddConverter(conv Converter)
    Get(result interface{}, url string, params map[string]interface{}) (int, error)
    Post(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
    Put(result interface{}, url string, params map[string]interface{}, requestBody interface{}) (int, error)
    Delete(result interface{}, url string, params map[string]interface{}) (int, error)
    Exchange(result interface{}, url string, method string, params map[string]interface{},
        requestBody interface{}) (int, error)
}
