// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import (
    "fmt"
    "github.com/xfali/restclient/transport"
    "io"
    "io/ioutil"
    "net/http"
    "time"
)

const (
    defaultTimeout = 15 * time.Second
)

type RequestCreator func(method, url string, r io.Reader, params map[string]interface{}) *http.Request

type DefaultRestClient struct {
    transport  http.RoundTripper
    converter  Converter
    timeout    time.Duration
    reqCreator RequestCreator
}

type Opt func(client *DefaultRestClient)

func New(opts ...Opt) RestClient {
    ret := &DefaultRestClient{
        transport:  defaultTransport,
        converter:  JsonConverter,
        timeout:    defaultTimeout,
        reqCreator: DefaultRequestCreator,
    }

    if opts != nil {
        for i := range opts {
            opts[i](ret)
        }
    }
    return ret
}

func (c *DefaultRestClient) AddConverter(conv Converter) {
    c.converter = conv
}

func (c *DefaultRestClient) Get(result interface{}, url string, params map[string]interface{}) (int, error) {
    return c.Exchange(result, url, http.MethodGet, nil, nil)
}

func (c *DefaultRestClient) Post(result interface{}, url string, params map[string]interface{}, param interface{}) (int, error) {
    return c.Exchange(result, url, http.MethodPost, nil, param)
}

func (c *DefaultRestClient) Put(result interface{}, url string, params map[string]interface{}, param interface{}) (int, error) {
    return c.Exchange(result, url, http.MethodPut, nil, param)
}

func (c *DefaultRestClient) Delete(result interface{}, url string, params map[string]interface{}) (int, error) {
    return c.Exchange(result, url, http.MethodDelete, nil, nil)
}

func (c *DefaultRestClient) Exchange(result interface{}, url string, method string, params map[string]interface{},
    requestBody interface{}) (int, error) {
    var r io.Reader
    if requestBody != nil {
        b, err := c.converter.Serialize(requestBody)
        if err != nil {
            return http.StatusBadRequest, err
        }
        r = b
    }

    request := c.reqCreator(method, url, r, params)
    request.Header.Set("Accept", c.converter.Accept())
    request.Header.Set("Content-Type", c.converter.ContentType())

    cli := c.newClient(c.timeout)
    resp, err := cli.Do(request)
    if err != nil {
        return http.StatusBadRequest, err
    }
    if resp.Body != nil {
        defer resp.Body.Close()
        if result == nil {
            io.Copy(ioutil.Discard, resp.Body)
            return http.StatusOK, nil
        }
        _, err := c.converter.Deserialize(resp.Body, result)
        if err != nil {
            return http.StatusBadRequest, err
        }
    }

    return http.StatusOK, nil
}

func (c *DefaultRestClient) newClient(timeout time.Duration) *http.Client {
    return &http.Client{
        Transport: c.transport,
        Timeout:   timeout,
    }
}

var defaultTransport = transport.New()

func SetTimeout(timeout time.Duration) func(client *DefaultRestClient) {
    return func(client *DefaultRestClient) {
        client.timeout = timeout
    }
}

func SetConverter(conv Converter) func(client *DefaultRestClient) {
    return func(client *DefaultRestClient) {
        client.AddConverter(conv)
    }
}

func SetRoundTripper(tripper http.RoundTripper) func(client *DefaultRestClient) {
    return func(client *DefaultRestClient) {
        client.transport = tripper
    }
}

func SetRequestCreator(f RequestCreator) func(client *DefaultRestClient) {
    return func(client *DefaultRestClient) {
        client.reqCreator = f
    }
}

func DefaultRequestCreator(method, url string, r io.Reader, params map[string]interface{}) *http.Request {
    request, err := http.NewRequest(method, url, r)
    if err != nil {
        return nil
    }
    if len(params) > 0 {
        for k, v := range params {
            request.Header.Set(k, interface2String(v))
        }
    }
    return request
}

func interface2String(v interface{}) string {
    if s, ok := v.(string); ok {
        return s
    }
    return fmt.Sprint(v)
}
