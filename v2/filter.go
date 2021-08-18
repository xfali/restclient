// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"net/http"
)

type IFilter interface {
	Filter(request *http.Request, fc FilterChain) (*http.Response, error)
}

type Filter func(request *http.Request, fc FilterChain) (*http.Response, error)

type FilterChain []Filter
type FilterManager FilterChain

func (fc *FilterManager) Add(filter ...Filter) {
	*fc = append(*fc, filter...)
}

func (fc FilterManager) RunFilter(request *http.Request) (*http.Response, error) {
	return FilterChain(fc).Filter(request)
}

func (fc FilterChain) Filter(request *http.Request) (*http.Response, error) {
	if len(fc) > 0 {
		filter := fc[len(fc)-1]
		chain := fc.next()
		return filter(request, chain)
	}
	return nil, nil
}

func (fc FilterChain) next() FilterChain {
	if len(fc) > 0 {
		return fc[:len(fc)-1]
	} else {
		return FilterChain{}
	}
}
