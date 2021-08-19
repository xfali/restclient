// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package request

type Setter interface {
	Set(key string, value interface{})
}

type Opt func(Setter)
