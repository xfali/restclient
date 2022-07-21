/*
 * Copyright 2022 Xiongfa Li.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
