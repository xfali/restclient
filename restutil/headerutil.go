// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restutil

import "encoding/base64"

const(
    HeaderAuthorization = "Authorization"
    HeaderContentType = "Content-Type"
    HeaderAccept = "Accept"
)

func BasicAuthHeader(username, password string) (string, string){
    return HeaderAuthorization, "Basic " + BasicAuth(username, password)
}

func BasicAuth(username, password string) string {
    auth := username + ":" + password
    return base64.StdEncoding.EncodeToString([]byte(auth))
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

func (b *HeaderBuilder) Builder() map[string]interface{} {
    return b.header
}