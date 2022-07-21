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

import "strings"

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
