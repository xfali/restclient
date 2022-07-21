# restclient

## 介绍 

  restclient是一个简单易用的RESTFUL客户端（连接池）。
  
  内置interface转换器：
  - bytes
  - string
  - xml
  - json
  - yaml
  
  内置支持认证方式：
  1. Basic Auth
  2. Digest Auth
  3. Token Auth
  
## 安装

使用命令安装：

```
go get github.com/xfali/restclient/v2
```

## 配置

### 基础配置

可以在创建默认client时对其进行配置，支持的options请参照[options](init_opts.go)的API说明

使用为
```
// restclient.New(Option1, Option2 ... OptionN), 例:
client := restclient.New(restclient.SetTimeout(10*time.Second))
```
```
//设置读写超时
restclient.SetTimeout(timeout time.Duration)
```
```
//配置初始转换器列表
restclient.SetConverters(convs []Converter)
```
```
//配置连接池
restclient.SetRoundTripper(tripper http.RoundTripper)
```
```
// 增加处理filter
restclient.AddFilter(filters ...Filter)
```
```
//配置request创建器
restclient.SetRequestCreator(f RequestCreator)
```
```
// CookieJar 配置http.Client的CookieJar
restclient.CookieJar(jar http.CookieJar)
```
```
// 配置http客户端创建器
restclient.SetClientCreator(cliCreator HttpClientCreator)
```
### 连接池配置

请参照http.transport的API说明

## 使用
1. 使用request传递http请求参数
```
//使用默认配置
client := restclient.New()
resp := &Response{}
err := client.Exchange("http://localhost:8080/test",
    request.WithResult(&ret),
    request.WithResponse(resp, false))
```
2. 使用request builder创建和传递http请求参数
```
err := client.Exchange("http://localhost:8080/error",
    restclient.NewRequest().
        MethodPost().
        RequestBody(req).
        Result(&resp).
        Build())
```

## 扩展

使用filter.Filter进行行为控制和扩展功能，如增加client的输入输出日志：
```
client := restclient.New(restclient.AddIFilter(filter.NewLog(xlog.GetLogger(), "")))
resp := &Response{}
err := client.Exchange("http://localhost:8080/test",
    request.WithResult(&ret),
    request.WithResponse(resp, false))
```
可以自行实现IFilter接口，并注册到restclient扩展其功能

## 认证

### Basic Auth

```
o := restclient.New(restclient.SetTimeout(time.Second))
auth := filter.NewBasicAuth("user", "password")
client := restclient.New(restclient.AddIFilter(auth))
resp := &Response{}
err := client.Exchange("http://localhost:8080/test",
    request.WithResult(&ret),
    request.WithResponse(resp, false))

//change username and password
auth.ResetCredentials(username, password)
```

### Digest Auth

```
auth := filter.NewDigestAuth("user", "password")
client := restclient.New(restclient.AddIFilter(auth))
resp := &Response{}
err := client.Exchange("http://localhost:8080/test",
    request.WithResult(&ret),
    request.WithResponse(resp, false))

//change username and password
auth.ResetCredentials(username, password)
```

### Token Auth

```
auth := filter.NewAccessTokenAuth("{TOKEN}")
client := restclient.New(restclient.AddIFilter(auth))
resp := &Response{}
err := client.Exchange("http://localhost:8080/test",
    request.WithResult(&ret),
    request.WithResponse(resp, false))

//change username and password
auth.ResetCredentials("{TOKEN}")
```

### 带日志client
```
client := restclient.New(restclient.AddIFilter(filter.NewLog(xlog.GetLogger(), "")))
resp := &Response{}
err := client.Exchange("http://localhost:8080/test",
    request.WithResult(&ret),
    request.WithResponse(resp, false))
```

### 捕捉panic
```
client := restclient.New(restclient.AddIFilter(filter.NewRecovery(xlog.GetLogger())))
resp := &Response{}
err := client.Exchange("http://localhost:8080/test",
    request.WithResult(&ret),
    request.WithResponse(resp, false))
```

## UrlBuilder
可以使用restclient.NewUrlBuilder为url添加参数，快速构建请求路径
```
builder := restclient.NewUrlBuilder("x/:a/tt/:b?")
builder.PathVariable("a", "1")
builder.PathVariable("b", 2)
builder.QueryVariable("c", 100)
builder.QueryVariable("d", 1.1)
url := builder.Build()
```