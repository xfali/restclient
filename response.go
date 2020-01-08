// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import "net/http"

type ResponseEntity struct {
    Data       interface{}
    Headers    map[string]string
    StatusCode int
}

func NewResponseEntity(result interface{}) *ResponseEntity {
    return &ResponseEntity{
        Data:    result,
        Headers: map[string]string{},
        StatusCode: http.StatusOK,
    }
}

func entity(ret interface{}) *ResponseEntity {
    if r, ok := ret.(*ResponseEntity); ok {
        return r
    }
    return nil
}
