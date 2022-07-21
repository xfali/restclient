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
	"strings"
)

const (
	MediaTypeAll               = "*/*"
	MediaTypeAtom              = "application/atom"
	MediaTypeAtomXml           = "application/atom+xml"
	MediaTypeFormUrlencoded    = "application/x-www-form-urlencoded"
	MediaTypeJson              = "application/json"
	MediaTypeJsonUtf8          = "application/json;charset=UTF-8"
	MediaTypeYaml              = "application/yaml"
	MediaTypeYamlUtf8          = "application/yaml;charset=UTF-8"
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

func BuildMediaType(t string, subType string) MediaType {
	return MediaType{t, subType}
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

func (t *MediaType) IsWildcard() bool {
	return t.t == "*"
}

func (t *MediaType) IsWildcardSub() bool {
	return t.sub == "*"
}

func (t *MediaType) isWildcardInnerSub() bool {
	if len(t.sub) > 0 && t.sub[:1] == "*" {
		return true
	}
	return false
}

func (t *MediaType) subEqual(o MediaType) bool {
	return strings.Index(o.sub, t.sub) == 0
}

func (t *MediaType) Includes(o MediaType) bool {
	if t.IsWildcard() {
		return true
	} else {
		if t.t == o.t {
			if t.IsWildcardSub() {
				return true
			}

			if t.subEqual(o) {
				return true
			}

			if t.isWildcardInnerSub() {
				wildSubType := t.sub[1:]
				if len(o.sub) >= len(wildSubType) {
					oSubType := o.sub[:len(wildSubType)]
					if wildSubType == oSubType {
						return true
					}
				}
			}
		}
	}
	return false
}

func (t MediaType) String() string {
	return t.t + "/" + t.sub
}
