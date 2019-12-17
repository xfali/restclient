/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package restclient

import "testing"

func TestGet(t *testing.T) {
    t.Run("get", func(t *testing.T) {
        c := New()
        str := ""
        _, err := c.Get(&str, "https://suggest.taobao.com/sug?code=utf-8")
        if err !=  nil {
            t.Fatal(err)
        }
        t.Log(str)
    })
}
