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
	"context"
	"errors"
	"github.com/xfali/restclient/v2/buffer"
	"github.com/xfali/restclient/v2/filter"
	"github.com/xfali/restclient/v2/reflection"
	"github.com/xfali/restclient/v2/request"
	"github.com/xfali/restclient/v2/restutil"
	"github.com/xfali/restclient/v2/transport"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type AcceptFlag int
type ResponseBodyFlag int

const (
	DefaultTimeout = 0

	// 仅接受用户指定的header
	AcceptUserOnly AcceptFlag = 1
	// 根据Converter支持的类型，自动添加第一个支持的类型（默认）
	AcceptAutoFirst AcceptFlag = 1 << 1
	// 根据Converter支持的类型，自动添加所有支持的类型
	AcceptAutoAll AcceptFlag = 1 << 2

	// 处理所有的response body
	ResponseBodyAll ResponseBodyFlag = 1
	// 不处理http status 400及以上的response的body
	ResponseBodyIgnoreBad ResponseBodyFlag = 1 << 1
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

type HttpClientCreator func() *http.Client

type defaultRestClient struct {
	client        *http.Client
	converters    []Converter
	filterManager filter.FilterManager
	pool          buffer.Pool

	cliCreator HttpClientCreator
	acceptFlag AcceptFlag
	respFlag   ResponseBodyFlag
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
		acceptFlag: AcceptAutoFirst,
		respFlag:   ResponseBodyAll,
	}
	ret.filterManager.Add(ret.filter)
	for _, opt := range opts {
		opt(ret)
	}
	if ret.cliCreator == nil {
		ret.cliCreator = ret.newClient
	}
	ret.client = ret.cliCreator()
	return ret
}

func (c *defaultRestClient) Exchange(url string, opts ...request.Opt) Error {
	param := emptyParam()
	for _, opt := range opts {
		opt(param)
	}

	// 序列化request body
	r, err := c.encodeRequest(param.reqBody, param.header)
	if r != nil {
		defer r.Close()
	}
	if err != nil {
		return withErr(DefaultErrorStatus, err)
	}

	nilResult := reflection.IsNil(param.result)
	if !nilResult {
		// 根据反序列化response body的目的result类型添加header Accept
		param.header = c.addAccept(param.result, param.header)
	}

	// 创建http.Request
	req := defaultRequestCreator(param.ctx, param.method, url, r, param.header)
	fm := c.filterManager
	if param.filterManager.Valid() {
		fm = filter.MergeFilterManager(c.filterManager, param.filterManager)
	}
	response, err := fm.RunFilter(req)
	if err != nil {
		return withErr(DefaultErrorStatus, err)
	}

	return c.processResponse(response, param, nilResult)
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

func (c *defaultRestClient) encodeRequest(requestBody interface{}, header http.Header) (io.ReadCloser, error) {
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

func (c *defaultRestClient) processResponse(response *http.Response, param *defaultParam, nilResult bool) Error {
	errStatus := response.StatusCode
	if response.StatusCode < http.StatusBadRequest {
		errStatus = DefaultErrorStatus
	} else if c.respFlag == ResponseBodyIgnoreBad {
		return withStatus(response.StatusCode)
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
			_, err := io.Copy(ioutil.Discard, response.Body)
			if err != nil {
				return withErr(errStatus, err)
			}
		} else {
			// 处理response
			err := c.decodeResponse(response, param.result)
			if err != nil {
				return withErr(errStatus, err)
			}
		}
	}

	if response.StatusCode >= http.StatusBadRequest {
		return withStatus(response.StatusCode)
	}

	return nil
}

func (c *defaultRestClient) decodeResponse(resp *http.Response, result interface{}) error {
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
	if c.acceptFlag != AcceptUserOnly {
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
							if c.acceptFlag == AcceptAutoFirst {
								break
							}
						}
					}
				}
			}
			if c.acceptFlag == AcceptAutoFirst && len(typeMap) > 0 {
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
