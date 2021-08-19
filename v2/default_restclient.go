// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"context"
	"errors"
	"github.com/xfali/restclient/buffer"
	"github.com/xfali/restclient/reflection"
	"github.com/xfali/restclient/restutil"
	"github.com/xfali/restclient/transport"
	"github.com/xfali/restclient/v2/filter"
	"github.com/xfali/restclient/v2/param"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type AcceptFlag int

const (
	DefaultTimeout             = 0
	AcceptUserOnly  AcceptFlag = 1
	AcceptAutoFirst AcceptFlag = 1 << 1
	AcceptAutoAll   AcceptFlag = 1 << 2
)

var (
	defaultTransport  = transport.New()
	defaultConverters = []Converter{
		NewByteConverter(),
		NewStringConverter(),
		NewXmlConverter(),
		NewJsonConverter(),
	}
)

type defaultParam struct {
	ctx           context.Context
	method        string
	header        http.Header
	filterManager filter.FilterManager

	reqBody  interface{}
	result   interface{}
	response *http.Response
	respFlag bool
}

func emptyParam() *defaultParam {
	return &defaultParam{
		method: http.MethodGet,
		ctx:    context.Background(),
	}
}

func (p *defaultParam) Set(key string, value interface{}) {
	switch key {
	case param.KeyMethod:
		p.method = value.(string)
	case param.KeyAddFilter:
		p.filterManager.Add(value.([]filter.Filter)...)
	case param.KeyRequestContext:
		p.ctx = value.(context.Context)
	case param.KeyRequestHeader:
		p.header = value.(http.Header)
	case param.KeyRequestAddHeader:
		ss := value.([]string)
		p.header.Add(ss[0], ss[1])
	case param.KeyRequestBody:
		p.reqBody = value
	case param.KeyResult:
		p.result = value
	case param.KeyResponse:
		rs := value.([]interface{})
		p.response = rs[0].(*http.Response)
		p.respFlag = rs[1].(bool)
	}
}

type defaultRestClient struct {
	client        *http.Client
	converters    []Converter
	filterManager filter.FilterManager
	pool          buffer.Pool

	autoAccept AcceptFlag
	transport  http.RoundTripper
	timeout    time.Duration
}

type Opt func(client *defaultRestClient)

func New(opts ...Opt) *defaultRestClient {
	ret := &defaultRestClient{
		transport:  defaultTransport,
		converters: defaultConverters,
		pool:       buffer.NewPool(),
		timeout:    DefaultTimeout,
		autoAccept: AcceptAutoFirst,
	}
	ret.filterManager.Add(ret.filter)
	for _, opt := range opts {
		opt(ret)
	}
	ret.client = ret.newClient()
	return ret
}

func (c *defaultRestClient) Exchange(url string, opts ...param.Parameter) error {
	param := emptyParam()
	for _, opt := range opts {
		opt(param)
	}

	r, err := c.processRequest(param.reqBody, param.header)
	if r != nil {
		defer r.Close()
	}
	if err != nil {
		return err
	}

	nilResult := reflection.IsNil(param.result)
	if !nilResult {
		param.header = c.addAccept(param.result, param.header)
	}

	request := defaultRequestCreator(param.ctx, param.method, url, r, param.header)
	fm := c.filterManager
	if param.filterManager.Valid() {
		fm = filter.MergeFilterManager(c.filterManager, param.filterManager)
	}
	response, err := fm.RunFilter(request)
	if err != nil {
		return err
	}

	if response.Body != nil {
		defer response.Body.Close()
		// need response
		if param.response != nil {
			copyResponse(param.response, response)
			// need response's body
			if param.respFlag {
				// get buffer form pool
				buf := buffer.NewReadWriteCloser(c.pool)
				// 封装reader，在读取response body数据时写入到buffer中
				reader := buffer.NewMergeReaderWriter(response.Body, buf)
				// 替换response body
				response.Body = ioutil.NopCloser(reader)
				// 调用者response设置body为buffer
				param.response.Body = buf
			}
		}
		if nilResult {
			// 如果用户没设置result，则直接读取body到discard
			_, err = io.Copy(ioutil.Discard, response.Body)
			return err
		}

		// 处理response
		err = c.processResponse(response, param.result)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (c *defaultRestClient) filter(request *http.Request, fc filter.FilterChain) (*http.Response, error) {
	return c.client.Do(request)
}

func copyResponse(dst, src *http.Response) {
	*dst = *src
	dst.Header = src.Header.Clone()
	// dst default without body
	dst.Body = nil
}

func defaultRequestCreator(ctx context.Context, method, url string, r io.Reader, header http.Header) *http.Request {
	request, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return nil
	}

	if len(header) > 0 {
		for k, vs := range header {
			for _, v := range vs {
				request.Header.Add(k, v)
			}
		}
	}
	return request
}

func (c *defaultRestClient) newClient() *http.Client {
	return &http.Client{
		Transport: c.transport,
		Timeout:   c.timeout,
	}
}

func (c *defaultRestClient) processRequest(requestBody interface{}, header http.Header) (io.ReadCloser, error) {
	if requestBody != nil {
		mtStr := getContentMediaType(header)
		mediaType := ParseMediaType(mtStr)
		conv, err := chooseEncoder(c.converters, requestBody, mediaType)
		if err != nil {
			return nil, err
		}
		if mtStr == "" {
			header.Set(restutil.HeaderContentType, getDefaultMediaType(conv).String())
		}
		// 从池中获得一个buffer
		buf := buffer.NewReadWriteCloser(c.pool)
		encoder := conv.CreateEncoder(buf)
		// 将序列化数据写入buffer
		_, err = encoder.Encode(requestBody)
		if err != nil {
			// 归还buffer
			_ = buf.Close()
			return nil, err
		}
		return buf, nil
	}
	return nil, nil
}

func (c *defaultRestClient) processResponse(resp *http.Response, result interface{}) error {
	mediaType := getResponseMediaType(resp)
	t := reflect.TypeOf(result)
	if t.Kind() != reflect.Func {
		conv, err := chooseDecoder(c.converters, result, mediaType)
		if err != nil {
			return err
		}
		decoder := conv.CreateDecoder(resp.Body)
		_, err = decoder.Decode(result)
		if err == io.EOF {
			return nil
		}
		return err
	} else {
		fn := reflect.ValueOf(result)
		if fn.Type().NumIn() != 1 || fn.Type().NumOut() != 0 {
			return errors.New("Function must be of type func(type) ")
		}
		inType := fn.Type().In(0)
		obj := reflect.New(inType).Interface()
		conv, err := chooseDecoder(c.converters, obj, mediaType)
		if err != nil {
			return err
		}
		decoder := conv.CreateDecoder(resp.Body)
		for {
			n, err := decoder.Decode(obj)
			if err != nil && err != io.EOF {
				return err
			}
			if n > 0 {
				var param [1]reflect.Value
				param[0] = reflect.ValueOf(obj).Elem()
				fn.Call(param[:])
			}
			if err == io.EOF {
				return nil
			}
		}
	}
	return nil
}

func (c *defaultRestClient) addAccept(result interface{}, header http.Header) http.Header {
	userAccept := getAcceptMediaType(header)
	mt := ParseMediaType(userAccept)
	typeMap := map[string]bool{}
	var acceptList []string
	if c.autoAccept != AcceptUserOnly {
		index := len(c.converters)
		for index > 0 {
			index--
			conv := c.converters[index]
			if conv.CanDecode(result, mt) {
				mts := conv.SupportMediaType()
				for _, v := range mts {
					if !v.isWildcardInnerSub() {
						mtStr := v.String()
						if _, have := typeMap[mtStr]; !have {
							acceptList = append(acceptList, mtStr)
							typeMap[mtStr] = true
							if c.autoAccept == AcceptAutoFirst {
								break
							}
						}
					}
				}
			}
			if c.autoAccept == AcceptAutoFirst && len(typeMap) > 0 {
				break
			}
		}
	} else {
		if userAccept != "" {
			acceptList = append(acceptList, userAccept)
		}
	}

	buf := strings.Builder{}
	l := len(acceptList)
	if l > 0 {
		for i := range acceptList {
			buf.WriteString(acceptList[i])
			l--
			if l > 0 {
				buf.WriteString(",")
			}
		}
	}

	header.Set(restutil.HeaderAccept, buf.String())
	return header
}

func getAcceptMediaType(header http.Header) string {
	if header != nil {
		if c := header.Get(restutil.HeaderAccept); c != "" {
			return c
		}
	}
	return ""
}

func getContentMediaType(header http.Header) string {
	if header != nil {
		if c := header.Get(restutil.HeaderContentType); c != "" {
			return c
		}
	}
	return ""
}

func getResponseMediaType(resp *http.Response) MediaType {
	mediaType := ""
	if resp != nil {
		mediaType = resp.Header.Get(restutil.HeaderContentType)
	}
	return ParseMediaType(mediaType)
}
