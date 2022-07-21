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

package restutil

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func QueryUrl(url string, params map[string]string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return url
	}
	buf := strings.Builder{}
	buf.WriteString(url)
	index := strings.LastIndex(url, "?")
	if index == -1 {
		buf.WriteString("?")
	} else {
		if index != len(url)-1 && url[len(url)-1:] != "&" {
			buf.WriteString("&")
		}
	}

	length := len(params)
	for k, v := range params {
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
		length--
		if length > 0 {
			buf.WriteString("&")
		}
	}
	return buf.String()
}

func PlaceholderUrl(url string, params map[string]string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return url
	}

	for k, v := range params {
		url = strings.Replace(url, "${"+k+"}", v, -1)
	}

	return url
}

func Query(keyAndValue ...interface{}) (string, error) {
	if len(keyAndValue) == 0 {
		return "", nil
	}
	if len(keyAndValue)%2 != 0 {
		return "", fmt.Errorf("Query parameter missing value, size is %d ", len(keyAndValue))
	}
	m := make(map[string]interface{}, len(keyAndValue)/2)
	for i := 0; i < len(keyAndValue); i += 2 {
		if k, ok := keyAndValue[i].(string); ok {
			m[k] = keyAndValue[i+1]
		} else {
			fmt.Errorf("Query key must be string, but get %s ", reflect.TypeOf(keyAndValue[i]).String())
		}
	}
	return EncodeQuery(m), nil
}

func EncodeQuery(keyAndValue map[string]interface{}) string {
	if len(keyAndValue) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for k, v := range keyAndValue {
		buf.WriteString(url.QueryEscape(fmt.Sprintf("%v", k)))
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(fmt.Sprintf("%v", v)))
		buf.WriteString("&")
	}
	format := buf.String()
	format = format[:len(format)-1]
	return format
}

func ReplaceUrl(uri string, leftDelim string, rightDelim string, keyAndValue map[string]interface{}) string {
	if len(keyAndValue) == 0 {
		return uri
	}
	for k, v := range keyAndValue {
		uri = strings.Replace(uri, fmt.Sprintf("%s%v%s", leftDelim, k, rightDelim), url.QueryEscape(fmt.Sprintf("%v", v)), -1)
	}
	return uri
}

type UrlBuilder struct {
	url        string
	leftDelim  string
	rightDelim string
	path       map[string]interface{}
	query      map[string]interface{}
}

// NewUrlBuilder URL构造器
func NewUrlBuilder(url string) *UrlBuilder {
	return &UrlBuilder{
		url:       url,
		leftDelim: ":",
	}
}

// Delims delimiters设置占位符，用于替换url的path参数
func (b *UrlBuilder) Delims(leftDelim, rightDelim string) *UrlBuilder {
	b.leftDelim = leftDelim
	b.rightDelim = rightDelim
	return b
}

// PathVariable 增加path变量参数
func (b *UrlBuilder) PathVariable(key string, value interface{}) *UrlBuilder {
	if b.path == nil {
		b.path = map[string]interface{}{}
	}
	b.path[key] = value
	return b
}

// QueryVariable 增加query参数
func (b *UrlBuilder) QueryVariable(key string, value interface{}) *UrlBuilder {
	if b.query == nil {
		b.query = map[string]interface{}{}
	}
	b.query[key] = value
	return b
}

// Build 创建url
func (b *UrlBuilder) Build() string {
	buf := strings.Builder{}
	if len(b.path) > 0 {
		buf.WriteString(ReplaceUrl(b.url, b.leftDelim, b.rightDelim, b.path))
	} else {
		buf.WriteString(b.url)
	}
	if len(b.query) > 0 {
		query := EncodeQuery(b.query)
		if b.url[len(b.url)-1] == '?' {
			buf.WriteString(query)
		} else {
			buf.WriteString("?")
			buf.WriteString(query)
		}
	}
	return buf.String()
}

func (b *UrlBuilder) String() string {
	return b.Build()
}
