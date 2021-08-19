// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package reflection

import (
	"reflect"
)

func IsNil(o interface{}) bool {
	if o == nil {
		return true
	}

	return reflect.ValueOf(o).IsNil()
}
