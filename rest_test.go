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
    t.Run("get", func(t *testing.T) {
        c := New(SetTimeout(time.Second))
        str := TestModel{}
        _, err := c.Get(&str, "https://suggest.taobao.com/sug?code=utf-8")
        if err != nil {
            t.Fatal(err)
        }
        t.Log(str)
    })
}

func TestWrapper(t *testing.T) {
    t.Run("get", func(t *testing.T) {
        o := New(SetTimeout(time.Second))
        c := NewWrapper(o, func(ex Exchange) Exchange {
            return func(result interface{}, url string, method string, header map[string]string, requestBody interface{}) (i int, e error) {
                t.Logf("url: %v, method: %v, header: %v, body: %v\n", url, method, header, requestBody)
                n, err := ex(result, url, method, header, requestBody)
                t.Logf("result %v", result)
                return n, err
            }
        })
        str := TestModel{}
        _, err := c.Get(&str, "https://suggest.taobao.com/sug?code=utf-8")
        if err != nil {
            t.Fatal(err)
        }
        t.Log(str)
    })
}
