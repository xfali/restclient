# restclient

## 介绍 

  restclient是一个简单易用的RESTFUL客户端（连接池）。
  
  内置interface转换器：
  - bytes
  - string
  - xml
  - json
  
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

