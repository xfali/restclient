// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package restclient

import "github.com/xfali/restclient/v2/request"

type RestClient interface {
	// 发起请求
	// url：请求路径
	// params：请求参数，见ex_params.go具体定义
	Exchange(url string, opts ...request.Opt) error
}
