/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package restclient

import (
    "testing"
    "time"
)

type TestModel struct {
    Result []string
}

func TestGet(t *testing.T) {
    t.Run("get_string", func(t *testing.T) {
        c := New(SetTimeout(time.Second))
        str := ""
        _, err := c.Get(&str, "https://suggest.taobao.com/sug?code=utf-8", nil)
        if err != nil {
            t.Fatal(err)
        }
        t.Log(str)
    })

    t.Run("get_bytes", func(t *testing.T) {
        c := New(SetTimeout(time.Second))
        var str []byte
        _, err := c.Get(&str, "https://suggest.taobao.com/sug?code=utf-8", nil)
        if err != nil {
            t.Fatal(err)
        }
        t.Log(string(str))
    })
}

func TestWrapper(t *testing.T) {
    t.Run("get", func(t *testing.T) {
        o := New(SetTimeout(time.Second))
        c := NewWrapper(o, func(ex Exchange) Exchange {
            return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
                t.Logf("url: %v, method: %v, params: %v, body: %v\n", url, method, params, requestBody)
                n, err := ex(result, url, method, params, requestBody)
                t.Logf("result %v", result)
                return n, err
            }
        })
        str := ""
        _, err := c.Get(&str, "https://suggest.taobao.com/sug?code=utf-8", nil)
        if err != nil {
            t.Fatal(err)
        }
        t.Log(str)
    })
}

func TestBasicAuth(t *testing.T) {
    t.Run("get", func(t *testing.T) {
        o := New(SetTimeout(time.Second))
        auth := BasicAuth{Username:"user", Password:"password"}
        c := NewBasicAuthClient(o, &auth)
        str := ""
        _, err := c.Get(&str, "https://suggest.taobao.com/sug?code=utf-8", nil)
        if err != nil {
            t.Fatal(err)
        }
        t.Log(str)
    })
}
