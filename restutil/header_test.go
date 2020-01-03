// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restutil

import (
    "github.com/xfali/restclient"
    "testing"
)

func TestBuilderHeader(t *testing.T) {
    x := Headers().
        WithAccept(restclient.MediaTypeJson).
        WithContentType(restclient.MediaTypeJson).
        WithBasicAuth("a", "b").
        WithKeyValue("key", "value").Builder()
    t.Log(x)
}
