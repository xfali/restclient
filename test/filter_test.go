// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
    "github.com/xfali/restclient"
    "net/http"
    "testing"
)

func TestFilter(t *testing.T) {
    fm := restclient.FilterManager{}
    filter1 := func(request *http.Request, fc restclient.FilterChain) (*http.Response, error) {
        t.Logf("filter1: %s\n", request.URL)
        return fc.Filter(request)
    }

    filter2 := func(request *http.Request, fc restclient.FilterChain) (*http.Response, error) {
        t.Logf("filter2: %s\n", request.URL)
        return fc.Filter(request)
    }

    filter3 := func(request *http.Request, fc restclient.FilterChain) (*http.Response, error) {
        t.Logf("filter3: %s\n", request.URL)
        return fc.Filter(request)
    }

    fm.Add(filter3, filter2, filter1)
    req, _ := http.NewRequest(http.MethodGet, "http://test.org", nil)
    _, err1 := fm.RunFilter(req)
    if err1 != nil {
        t.Fatal(err1)
    }
}

func TestFilter2(t *testing.T) {
    fm := restclient.FilterManager{}
    filter1 := func(request *http.Request, fc restclient.FilterChain) (*http.Response, error) {
        t.Logf("filter1: %s\n", request.URL)
        return  fc.Filter(request)
    }

    filter2 := func(request *http.Request, fc restclient.FilterChain) (*http.Response, error) {
        t.Logf("filter2: %s\n", request.URL)
        return nil, nil
    }

    filter3 := func(request *http.Request, fc restclient.FilterChain) (*http.Response, error) {
        t.Fatal("cannot be here!")
        return  fc.Filter(request)
    }

    fm.Add(filter3, filter2, filter1)

    fm.Add(filter3, filter2, filter1)
    req, _ := http.NewRequest(http.MethodGet, "http://test.org", nil)
    _, err1 := fm.RunFilter(req)
    if err1 != nil {
        t.Fatal(err1)
    }

    fm.Add(filter3, filter2, filter1)
    req, _ = http.NewRequest(http.MethodGet, "http://test.org", nil)
    _, err2 := fm.RunFilter(req)
    if err1 != nil {
        t.Fatal(err2)
    }
}
