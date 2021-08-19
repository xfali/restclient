/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
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
