// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import (
    "bytes"
    "github.com/xfali/restclient/restutil"
    "io"
    "net/http"
    "net/url"
)

func NewBasicAuthClient(client RestClient, auth *BasicAuth) RestClient {
    return NewWrapper(client, auth.Exchange)
}

func (b *BasicAuth) Exchange(ex Exchange) Exchange {
    return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
        if params == nil {
            params = map[string]interface{}{}
        }
        k, v := restutil.BasicAuthHeader(b.Username, b.Password)
        params[k] = v
        n, err := ex(result, url, method, params, requestBody)
        return n, err
    }
}

func NewDigestAuthClient(client RestClient, auth *DigestAuth) RestClient {
    return NewWrapper(client, auth.Exchange)
}

type DigestReader struct {
    buf bytes.Buffer
}

func (dr *DigestReader) Reader(r io.Reader) io.Reader {
    io.Copy(&dr.buf, r)
    return bytes.NewReader(dr.buf.Bytes())
}

func (b *DigestAuth) Exchange(ex Exchange) Exchange {
    return func(result interface{}, uri string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
        ent := entity(result)
        if ent == nil {
            ent = NewResponseEntity(result)
        }
        digestBuf := DigestReader{}
        if requestBody != nil {
            body := body(requestBody)
            if body == nil {
                body = NewRequestBody(requestBody, digestBuf.Reader)
            } else {
                originReader := body.Reader
                body.Reader = func(r io.Reader) io.Reader {
                    return originReader(digestBuf.Reader(r))
                }
            }
        }
        n, err := ex(ent, uri, method, params, requestBody)
        if n == http.StatusUnauthorized {
            digest := findWWWAuth(ent.Headers)
            wwwAuth := ParseWWWAuthenticate(digest)
            uriP, _ := url.Parse(uri)
            err := b.Refresh(method, uriP.RequestURI(), digestBuf.buf.Bytes(), wwwAuth)
            if err != nil {
                return n, err
            }
            auth, err := b.ToString()
            if err != nil {
                return n, err
            }
            if params == nil {
                params = map[string]interface{}{}
            }
            params["Authorization"] = auth
            return ex(result, uri, method, params, requestBody)
        }
        return n, err
    }
}

func findWWWAuth(headers map[string]string) string{
    if digest, ok := headers["WWW-Authenticate"]; ok {
        return digest
    }

    if digest, ok := headers["Www-Authenticate"]; ok {
        return digest
    }

    if digest, ok := headers["www-Authenticate"]; ok {
        return digest
    }

    return ""
}
