// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import "github.com/xfali/restclient/restutil"

type BasicAuth struct {
    Username string
    Password string
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

func NewBasicAuthClient(client RestClient, auth *BasicAuth) RestClient {
    return NewWrapper(client, auth.Exchange)
}
