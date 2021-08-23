// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"fmt"
	"net/http"
)

var DefaultErrorStatus = http.StatusBadRequest

type Error interface {
	error

	// 获得http status code
	StatusCode() int

	// 获得原始error
	Origin() error
}

type defaultError struct {
	status int
	err    error
}

func withErr(status int, err error) defaultError {
	return defaultError{
		status: status,
		err:    err,
	}
}

func withStatus(status int) defaultError {
	return defaultError{
		status: status,
		err:    fmt.Errorf("restclient status: [%d] %s", status, http.StatusText(status)),
	}
}

func (e defaultError) Origin() error {
	return e.err
}

func (e defaultError) Error() string {
	if e.err != nil {
		return e.err.Error()
	} else {
		return ""
	}
}

func (e defaultError) StatusCode() int {
	return e.status
}
