// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xfali/restclient/restutil"
	"github.com/xfali/restclient/transport"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	defaultTimeout             = 15 * time.Second
	AcceptUserOnly  AcceptFlag = 1
	AcceptAutoFirst AcceptFlag = 1 << 1
	AcceptAutoAll   AcceptFlag = 1 << 2
)

type AcceptFlag int

type RequestCreator func(method, url string, r io.Reader, params map[string]interface{}) *http.Request

type DefaultRestClient struct {
	transport  http.RoundTripper
	converters []Converter
	timeout    time.Duration
	reqCreator RequestCreator
	autoAccept AcceptFlag
}

type Opt func(client *DefaultRestClient)

var defaultConverters = []Converter{
	NewByteConverter(),
	NewStringConverter(),
	NewXmlConverter(),
	NewJsonConverter(),
}

func New(opts ...Opt) RestClient {
	ret := &DefaultRestClient{
		transport:  defaultTransport,
		converters: defaultConverters,
		timeout:    defaultTimeout,
		reqCreator: DefaultRequestCreator,
		autoAccept: AcceptAutoFirst,
	}

	for i := range opts {
		opts[i](ret)
	}
	return ret
}

func (c *DefaultRestClient) AddConverter(conv Converter) {
	c.converters = append(c.converters, conv)
}

func (c *DefaultRestClient) GetConverters() []Converter {
	return c.converters
}

func (c *DefaultRestClient) Get(result interface{}, url string, params map[string]interface{}) (int, error) {
	return c.Exchange(result, url, http.MethodGet, params, nil)
}

func (c *DefaultRestClient) Post(result interface{}, url string, params map[string]interface{}, body interface{}) (int, error) {
	return c.Exchange(result, url, http.MethodPost, params, body)
}

func (c *DefaultRestClient) Put(result interface{}, url string, params map[string]interface{}, body interface{}) (int, error) {
	return c.Exchange(result, url, http.MethodPut, params, body)
}

func (c *DefaultRestClient) Delete(result interface{}, url string, params map[string]interface{}) (int, error) {
	return c.Exchange(result, url, http.MethodDelete, params, nil)
}

func (c *DefaultRestClient) Head(result interface{}, url string, params map[string]interface{}) (int, error) {
	return c.Exchange(result, url, http.MethodHead, params, nil)
}

func (c *DefaultRestClient) Options(result interface{}, url string, params map[string]interface{}) (int, error) {
	return c.Exchange(result, url, http.MethodOptions, params, nil)
}

func (c *DefaultRestClient) Patch(result interface{}, url string, params map[string]interface{}, body interface{}) (int, error) {
	return c.Exchange(result, url, http.MethodPatch, params, body)
}

func (c *DefaultRestClient) Exchange(result interface{}, url string, method string, params map[string]interface{},
	requestBody interface{}) (int, error) {
	if params == nil {
		params = map[string]interface{}{}
	}
	r, err := c.processRequest(requestBody, params)
	if err != nil {
		return http.StatusBadRequest, err
	}

	entity := entity(result)
	if entity != nil {
		result = entity.Result
	}

	nilResult := IsNil(result)
	if !nilResult {
		c.addAccept(result, &params)
	}

	request := c.reqCreator(method, url, r, params)
	if request == nil {
		return http.StatusBadRequest, fmt.Errorf("Request nil. method: %s , url: %s , params: %v\n", method, url, params)
	}

	cli := c.newClient(c.timeout)
	resp, err := cli.Do(request)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if entity != nil {
		entity.fill(resp)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		if nilResult {
			_, err := io.Copy(ioutil.Discard, resp.Body)
			return resp.StatusCode, err
		}

		err = c.processResponse(resp, result)
		if err != nil {
			if resp.StatusCode >= 400 {
				return resp.StatusCode, err
			} else {
				return http.StatusBadRequest, err
			}
		}
	}

	return resp.StatusCode, nil
}

func (c *DefaultRestClient) processRequest(requestBody interface{}, params map[string]interface{}) (io.Reader, error) {
	if requestBody != nil {
		reqBody := body(requestBody)
		if reqBody != nil {
			requestBody = reqBody.Body
		}
		if requestBody != nil {
			mtStr := getContentMediaType(params)
			mediaType := ParseMediaType(mtStr)
			conv, err := chooseEncoder(c.converters, requestBody, mediaType)
			if err != nil {
				return nil, err
			}
			if mtStr == "" {
				params[restutil.HeaderContentType] = getDefaultMediaType(conv).String()
			}
			buf := bytes.NewBuffer(nil)
			encoder := conv.CreateEncoder(buf)
			_, err = encoder.Encode(requestBody)
			if err != nil {
				return nil, err
			}
			if reqBody != nil && reqBody.Reader != nil {
				return reqBody.Reader(buf), nil
			}
			return buf, nil
		}
	}
	return nil, nil
}

func (c *DefaultRestClient) processResponse(resp *http.Response, result interface{}) error {
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

func (c *DefaultRestClient) addAccept(result interface{}, params *map[string]interface{}) {
	userAccept := getAcceptMediaType(*params)
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

	(*params)[restutil.HeaderAccept] = buf.String()
}

func (c *DefaultRestClient) newClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: c.transport,
		Timeout:   timeout,
	}
}

var defaultTransport = transport.New()

//设置读写超时
func SetTimeout(timeout time.Duration) func(client *DefaultRestClient) {
	return func(client *DefaultRestClient) {
		client.timeout = timeout
	}
}

//配置初始转换器列表
func SetConverters(convs []Converter) func(client *DefaultRestClient) {
	return func(client *DefaultRestClient) {
		client.converters = convs
	}
}

//配置连接池
func SetRoundTripper(tripper http.RoundTripper) func(client *DefaultRestClient) {
	return func(client *DefaultRestClient) {
		client.transport = tripper
	}
}

//配置request创建器
func SetRequestCreator(f RequestCreator) func(client *DefaultRestClient) {
	return func(client *DefaultRestClient) {
		client.reqCreator = f
	}
}

//配置是否自动添加accept
func SetAutoAccept(v AcceptFlag) func(client *DefaultRestClient) {
	return func(client *DefaultRestClient) {
		client.autoAccept = v
	}
}

func DefaultRequestCreator(method, url string, r io.Reader, params map[string]interface{}) *http.Request {
	request, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil
	}
	if len(params) > 0 {
		for k, v := range params {
			request.Header.Set(k, interface2String(v))
		}
	}
	return request
}

func interface2String(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

func getAcceptMediaType(params map[string]interface{}) string {
	if params != nil {
		if c, ok := params[restutil.HeaderAccept]; ok && c != nil {
			if t, ok := c.(string); ok {
				return t
			}
		}
	}
	return ""
}

func getContentMediaType(params map[string]interface{}) string {
	if params != nil {
		if c, ok := params[restutil.HeaderContentType]; ok && c != nil {
			if t, ok := c.(string); ok {
				return t
			}
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
