/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package restclient

import (
	"context"
	"fmt"
	"github.com/xfali/restclient/restutil"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

type TestModel struct {
	Result []string
}

func init() {
	go startHttpServer(5 * time.Second)
	time.Sleep(time.Second)
}

func startHttpServer(shutdown time.Duration) {
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		v := request.Header.Get(restutil.HeaderAuthorization)
		fmt.Println(v)
		writer.Header().Set(restutil.HeaderContentType, "application/json")
		_, err := writer.Write([]byte(`{ "result":["hello", "world"]}`))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
	})
	server := &http.Server{Addr: ":8080", Handler: nil}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			os.Exit(-1)
		}
	}()

	<-time.NewTimer(shutdown).C

	err := server.Shutdown(context.Background())
	if err != nil {
		os.Exit(-1)
	}
}

func TestGet(t *testing.T) {
	t.Run("get_string", func(t *testing.T) {
		c := New(SetTimeout(time.Second))
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})

	t.Run("get_bytes", func(t *testing.T) {
		c := New(SetTimeout(time.Second))
		var str []byte
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(str))
	})
}

func TestPost(t *testing.T) {
	t.Run("get_string", func(t *testing.T) {
		c := New(SetTimeout(time.Second))
		str := ""
		_, err := c.Post(&str, "http://localhost:8080/test", nil, time.Time{})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})

	t.Run("get_struct", func(t *testing.T) {
		c := New(SetTimeout(time.Second))
		str := TestModel{}
		_, err := c.Post(&str, "http://localhost:8080/test", nil, time.Time{})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})

	t.Run("get_resp entity struct", func(t *testing.T) {
		c := New(SetTimeout(time.Second))
		str := TestModel{}
		_, err := c.Post(NewResponseEntity(&str), "http://localhost:8080/test", nil, time.Time{})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestWrapper(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		o := New(SetTimeout(time.Second))
		c := NewWrapper(o, func(ex Exchange) Exchange {
			return func(result interface{}, url string, method string, params map[string]interface{}, requestBody interface{}) (i int, e error) {
				now := time.Now()
				id := RandomId(10)
				t.Logf("[restclient request %s]: url: %v , method: %v , params: %v , body: %v \n",
					id, url, method, params, requestBody)
				n, err := ex(result, url, method, params, requestBody)
				v := reflect.ValueOf(result)
				v = reflect.Indirect(v)
				t.Logf("[restclient response %s]: use time: %d ms, result: %v ",
					id, time.Since(now)/time.Millisecond, v.Interface())
				return n, err
			}
		})
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestBasicAuth(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		o := New(SetTimeout(time.Second))
		auth := NewBasicAuth("user", "password")
		c := NewBasicAuthClient(o, auth)
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestAccessTokenAuth(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		o := New(SetTimeout(time.Second))
		auth := NewAccessTokenAuth("mytoken")
		c := NewAccessTokenAuthClient(o, auth)
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestDigestAuth(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		o := New(SetTimeout(time.Second))
		auth := NewDigestAuth("user", "pw")
		c := NewDigestAuthClient(o, auth)
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestLog(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		c := NewLogClient(New(SetTimeout(time.Second)), NewLog(t.Logf, "test"))
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})

	t.Run("get resp entity", func(t *testing.T) {
		c := NewLogClient(New(SetTimeout(time.Second)), NewLog(t.Logf, "test"))
		str := ""
		_, err := c.Get(NewResponseEntity(&str), "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestAccept(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		builder := Builder{}
		c := builder.Default().
			Log(NewLog(t.Logf, "Mytag")).
			Build()
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})

	t.Run("none_struct", func(t *testing.T) {
		builder := Builder{}
		c := builder.Default().
			Log(NewLog(t.Logf, "Mytag")).
			Build()
		m := TestModel{}
		_, err := c.Get(&m, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(m)
	})

	t.Run("json", func(t *testing.T) {
		builder := Builder{}
		c := builder.Default().
			Log(NewLog(t.Logf, "Mytag")).
			Build()
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", restutil.Headers().WithAccept(MediaTypeJson).Build())
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestBuilder(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		builder := Builder{}
		c := builder.Default().
			Log(NewLog(t.Logf, "Mytag")).
			BasicAuth(NewBasicAuth("user", "pw")).
			Build()
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}
