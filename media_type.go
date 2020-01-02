// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package restclient

import "strings"

const (
    MediaTypeAll               = "*/*"
    MediaTypeAtom              = "application/atom"
    MediaTypeAtomXml           = "application/atom+xml"
    MediaTypeFormUrlencoded    = "application/x-www-form-urlencoded"
    MediaTypeJson              = "application/json"
    MediaTypeJsonUtf8          = "application/json;charset=UTF-8"
    MediaTypeOctetStream       = "application/octet-stream"
    MediaTypePdf               = "application/pdf"
    MediaTypeProblemJson       = "application/problem+json"
    MediaTypeProblemJsonUtf8   = "application/problem+json;charset=UTF-8"
    MediaTypeXml               = "application/xml"
    MediaTypeProblemXml        = "application/problem+xml"
    MediaTypeRssXml            = "application/rss+xml"
    MediaTypeStreamJson        = "application/stream+json"
    MediaTypeXhtmlXml          = "application/xhtml+xml"
    MediaTypeImageAll          = "image/*"
    MediaTypeImageGif          = "image/gif"
    MediaTypeImageJpeg         = "image/jpeg"
    MediaTypeImagePng          = "image/png"
    MediaTypeMultipartFormData = "multipart/form-data"
    MediaTypeTextEventStream   = "text/event-stream"
    MediaTypeTextHtml          = "text/html"
    MediaTypeTextMarkdown      = "text/markdown"
    MediaTypeTextPlain         = "text/plain"
    MediaTypeTextXml           = "text/xml"
)

type MediaType struct {
    t   string
    sub string
}

func ParseMediaType(s string) MediaType {
    if s == "" {
        s = MediaTypeAll
    }
    s = strings.ToLower(strings.TrimSpace(s))
    strs := strings.Split(s, "/")
    if strs[0] == "" {
        strs[0] = "*"
    }
    if len(strs) == 1 {
            return MediaType{strs[0], "*"}
    } else if len(strs) > 1 {
        if strs[1] == "" {
            strs[1] = "*"
        }
    }
    return MediaType{strs[0], strs[1]}
}

func (t *MediaType) Includes(o MediaType) bool {
    if t.t == "*" {
        return true
    } else {
        if t.t == o.t {
            if t.sub == o.sub {
                return true
            }
            if t.sub == "*" {
                return true
            }

            if len(t.sub) > 1 && t.sub[:1] == "*" {
                wildSubType := t.sub[1:]
                oSubType := o.sub[:len(wildSubType)]
                if wildSubType == oSubType {
                    return true
                }
            }
        }
    }
    return false
}

