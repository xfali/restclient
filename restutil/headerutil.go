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

import "encoding/base64"

const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderAccept        = "Accept"

	Bearer = "bearer"
)

func BasicAuthHeader(username, password string) (string, string) {
	return HeaderAuthorization, "Basic " + BasicAuth(username, password)
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func AccessTokenAuthHeader(token string) (string, string) {
	return HeaderAuthorization, Bearer + " " + token
}

type HeaderBuilder struct {
	header map[string]interface{}
}

func Headers() *HeaderBuilder {
	return &HeaderBuilder{map[string]interface{}{}}
}

func (b *HeaderBuilder) WithBasicAuth(username, password string) *HeaderBuilder {
	k, v := BasicAuthHeader(username, password)
	b.header[k] = v
	return b
}

func (b *HeaderBuilder) WithContentType(ct string) *HeaderBuilder {
	b.header[HeaderContentType] = ct
	return b
}

func (b *HeaderBuilder) WithAccept(ct string) *HeaderBuilder {
	b.header[HeaderAccept] = ct
	return b
}

func (b *HeaderBuilder) WithKeyValue(k, v string) *HeaderBuilder {
	b.header[k] = v
	return b
}

func (b *HeaderBuilder) Build() map[string]interface{} {
	return b.header
}
