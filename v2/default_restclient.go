// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

const(
	KeyMethod = "method"
)

type defaultParam struct {
	method string
}

func (p *defaultParam) Set(key string, value interface{}) {
	switch key {
	case KeyMethod:
		p.method = value.(string)
	}
}

type defaultRestClient struct {

}

//func (c *defaultRestClient) Get(url string, opts ...Opt) error {
//	return c.Exchange(url, )
//}
//
//func (c *defaultRestClient) Post(url string, opts ...Opt) error {
//
//}
//
//func (c *defaultRestClient) Put(url string, opts ...Opt) error {
//
//}
//
//func (c *defaultRestClient) Delete(url string, opts ...Opt) error {
//
//}
//
//func (c *defaultRestClient) Head(url string, opts ...Opt) error {
//
//}
//
//func (c *defaultRestClient) Options(url string, opts ...Opt) error {
//
//}
//func (c *defaultRestClient) Patch(url string, opts ...Opt) error {
//
//}

func (c *defaultRestClient) Exchange(url string, opts ...Opt) error {
	return nil
}
