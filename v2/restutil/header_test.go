// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restutil

import (
	"testing"
)

func TestBuilderHeader(t *testing.T) {
	x := Headers().
		WithAccept("application/json").
		WithContentType("application/json").
		WithBasicAuth("a", "b").
		WithKeyValue("key", "value").Build()
	t.Log(x)
}
