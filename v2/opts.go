// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import "net/http"

func OptSetMethod(method string) Opt {
	return func(setter Setter) {
		setter.Set(KeyMethod, method)
	}
}

func MethodGet() Opt {
	return OptSetMethod(http.MethodGet)
}

func MethodPost() Opt {
	return OptSetMethod(http.MethodPost)
}

func MethodPut() Opt {
	return OptSetMethod(http.MethodPut)
}

func MethodDelete() Opt {
	return OptSetMethod(http.MethodDelete)
}

func MethodHead() Opt {
	return OptSetMethod(http.MethodHead)
}

func MethodPatch() Opt {
	return OptSetMethod(http.MethodPatch)
}

func MethodOptions() Opt {
	return OptSetMethod(http.MethodOptions)
}

func MethodConnect() Opt {
	return OptSetMethod(http.MethodConnect)
}

func MethodTrace() Opt {
	return OptSetMethod(http.MethodTrace)
}
