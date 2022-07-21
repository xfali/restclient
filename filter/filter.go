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

package filter

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

func (fc FilterManager) Valid() bool {
	return len(fc) > 0
}

func (fc FilterManager) RunFilter(request *http.Request) (*http.Response, error) {
	return FilterChain(fc).Filter(request)
}

func MergeFilterManager(fms ...FilterManager) FilterManager {
	ret := make([]Filter, 0, 64)
	for _, v := range fms {
		ret = append(ret, v...)
	}
	return ret
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
