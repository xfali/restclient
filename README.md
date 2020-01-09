# restclient

## 介绍 

  restclient是一个简单易用的RESTFUL客户端（连接池）。
  
  内置interface转换器：
  - bytes
  - string
  - xml
  - json
  
  内置支持认证方式：
  1. Basic Auth
  2. Digest Auth
  
## 安装

使用命令安装：

```
go get github.com/xfali/restclient
```

## 配置

### 基础配置

请参照DefaultRestClient的API说明
```cassandraql
//设置读写超时
func SetTimeout(timeout time.Duration)
```
```cassandraql
//配置初始转换器列表
func SetConverters(convs []Converter)
```
```cassandraql
//配置连接池
func SetRoundTripper(tripper http.RoundTripper)
```
```cassandraql
//配置request创建器
func SetRequestCreator(f RequestCreator)
```
### 连接池配置

请参照transport的API说明

## 使用

```cassandraql
//使用默认配置
client := restclient.New()
str := ""
n, err := c.Get(&str, "https://${ADDRESS}", nil)
n, err := c.Post(&str, "https://${ADDRESS}", 
            restutil.Headers().WithContentType(MediaTypeJson).Build(), Entity{})
```

## 扩展

使用ClientWrapper进行行为控制和扩展功能，如增加client的输入输出日志：
```cassandraql
o := restclient.New(restclient.SetTimeout(time.Second))
c := restclient.NewWrapper(o, func(ex restclient.Exchange) restclient.Exchange {
    return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
        t.Logf("url: %v, method: %v, params: %v, body: %v\n", url, method, params, requestBody)
        n, err := ex(result, url, method, params, requestBody)
        t.Logf("result %v", result)
        return n, err
    }
})
str := ""
n, err := c.Get(&str, "https://${ADDRESS}", nil)
```

## 认证

### Basic Auth

```cassandraql
o := restclient.New(restclient.SetTimeout(time.Second))
auth := restclient.NewBasicAuth("user", "password")
c := restclient.NewBasicAuthClient(o, auth)
str := ""
_, err := c.Get(&str, "https://${ADDRESS}", nil)

//change username and password
auth.Username = "other_user"
auth.Password = "other_password"
```

### Digest Auth

```cassandraql
o := restclient.New(restclient.SetTimeout(time.Second))
auth := restclient.NewDigestAuth("user", "password")
c := restclient.NewDigestAuthClient(o, auth)
str := ""
_, err := c.Get(&str, "https://${ADDRESS}", nil)

//change username and password
auth.Username = "other_user"
auth.Password = "other_password"
```

### 带日志client
```$xslt
c := restclient.NewLogClient(restclient.New(), restclient.NewLog(t.Logf, "test"))
str := ""
_, err := c.Get(&str, "http://${ADDRESS}", nil)
```
