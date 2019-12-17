// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import (
    "bytes"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "time"
)

type DefaultRestClient struct {
    transport http.RoundTripper
    converter Converter
    timeout   int
}

func New() RestClient {
    return &DefaultRestClient{
        transport: transport,
        converter: &JsonConverter{},
    }
}

func (c *DefaultRestClient) Init(conv Converter, timeout int) {
    if conv != nil {
        c.converter = conv
    }
    c.timeout = timeout
}

func (c *DefaultRestClient) Get(result interface{}, url string) (int, error) {
    return c.Exchange(result, url, http.MethodGet, nil, nil)
}

func (c *DefaultRestClient) Post(result interface{}, url string, param interface{}) (int, error) {
    return c.Exchange(result, url, http.MethodPost, nil, param)
}

func (c *DefaultRestClient) Put(result interface{}, url string, param interface{}) (int, error) {
    return c.Exchange(result, url, http.MethodPut, nil, param)
}

func (c *DefaultRestClient) Delete(result interface{}, url string) (int, error) {
    return c.Exchange(result, url, http.MethodDelete, nil, nil)
}

func (c *DefaultRestClient) Exchange(result interface{}, url string, method string, header map[string]string,
    requestBody interface{}) (int, error) {
    var r io.Reader
    if requestBody != nil {
        b, err := c.converter.Serialize(requestBody)
        if err != nil {
            return http.StatusBadRequest, err
        }
        r = bytes.NewReader(b)
    }
    request, err := http.NewRequest(method, url, r)
    if err != nil {
        return http.StatusBadRequest, err
    }
    if header != nil {
        for k, v := range header {
            request.Header.Set(k, v)
        }
        request.Header.Set("Accept", c.converter.Accept())
        request.Header.Set("Content-Type", c.converter.ContentType())
    }
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
        buf := bytes.NewBuffer(nil)
        _, err := io.Copy(buf, resp.Body)
        if err != nil {
            return http.StatusBadRequest, err
        }

        d := buf.Bytes()
        if d != nil && len(d) > 0 {
            err := c.converter.Deserialize(d, result)
            if err != nil {
                return http.StatusBadRequest, err
            }
        }
    }

    return http.StatusOK, nil
}

func (c *DefaultRestClient) newClient(timeoutMsec int) *http.Client {
    return &http.Client{
        Transport: c.transport,
        Timeout:   time.Duration(timeoutMsec) * time.Millisecond,
    }
}

var transport = &http.Transport{
    Proxy: http.ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
        DualStack: true,
    }).DialContext,
    MaxIdleConns:          3000,
    MaxIdleConnsPerHost:   3000,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
}
