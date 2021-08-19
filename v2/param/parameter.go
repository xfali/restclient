// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package param

type Setter interface {
	Set(key string, value interface{})
}

type Parameter func(Setter)
