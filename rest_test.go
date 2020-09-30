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
	Result []string `json:"result"`
}

func init() {
	go startHttpServer(100 * time.Second)
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
	http.HandleFunc("/cookie", func(writer http.ResponseWriter, request *http.Request) {
		if c, err := request.Cookie("SESSION"); err != nil || c.Value == "" {
			v := request.Header.Get(restutil.HeaderAuthorization)
			fmt.Println(v)
			if v == "" {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				http.SetCookie(writer, &http.Cookie{Name: "SESSION", Value: v, MaxAge: 0})
				return
			}
		}

		writer.Header().Set(restutil.HeaderContentType, "application/json")
		_, err := writer.Write([]byte(`{ "result":["hello", "world"]}`))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
	})
	http.HandleFunc("/test/chunk", func(writer http.ResponseWriter, request *http.Request) {
		v := request.Header.Get(restutil.HeaderAuthorization)
		fmt.Println(v)
		writer.Header().Set(restutil.HeaderContentType, "application/json")
		writer.Header().Set("Transfer-Encoding", "chunked")
		for i := 0; i < 5; i++ {
			_, err := writer.Write([]byte(`{ "result":["hello", "world"]}`))
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
			}
			writer.(http.Flusher).Flush()
			time.Sleep(time.Second)
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

func TestFilter(t *testing.T) {
	t.Run("get_string", func(t *testing.T) {
		c := New(SetTimeout(time.Second), AddFilter(func(request *http.Request, fc FilterChain) (response *http.Response, e error) {
			t.Log(request.URL)
			resp, err := fc.Filter(request)
			t.Log(resp.Status)
			return resp, err
		}))
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestPost(t *testing.T) {
	t.Run("get_string", func(t *testing.T) {
		c := New(SetTimeout(time.Second), AddFilter(NewLog(t.Logf, "").Filter))
		str := ""
		_, err := c.Post(&str, "http://localhost:8080/test", nil, TestModel{})
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
		o := New(SetTimeout(time.Second), AddFilter(
			NewLog(t.Logf, "").Filter,
			NewDigestAuth("user", "pw").Filter))
		str := ""
		_, err := o.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestLog(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		c := New(SetTimeout(time.Second), AddFilter(NewLog(t.Logf, "test").Filter))
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})

	t.Run("get resp entity", func(t *testing.T) {
		c := New(SetTimeout(time.Second), AddFilter(NewLog(t.Logf, "test").Filter))
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
		c := New(SetTimeout(time.Second), AddFilter(NewLog(t.Logf, "test").Filter))
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})

	t.Run("none_struct", func(t *testing.T) {
		c := New(SetTimeout(time.Second), AddFilter(NewLog(t.Logf, "test").Filter))
		m := TestModel{}
		_, err := c.Get(&m, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(m)
	})

	t.Run("json", func(t *testing.T) {
		c := New(SetTimeout(time.Second), AddFilter(NewLog(t.Logf, "test").Filter))
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
		c := New(SetTimeout(time.Second),
			AddFilter(NewLog(t.Logf, "test").Filter, NewBasicAuth("user", "pw").Filter),
		)
		str := ""
		_, err := c.Get(&str, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(str)
	})
}

func TestChunkGet(t *testing.T) {
	t.Run("get_string_chunked", func(t *testing.T) {
		c := New(SetTimeout(0))
		_, err := c.Get(func(s string) {
			fmt.Printf("%s\n", s)
		}, "http://localhost:8080/test/chunk", nil)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("get_struct_chunked", func(t *testing.T) {
		c := New(SetTimeout(0))
		i := 0
		v := TestModel{Result: []string{"hello", "world"}}
		_, err := c.Get(func(s TestModel) {
			i++
			fmt.Printf("%v\n", s)
			if !reflect.DeepEqual(v, s) {
				t.Fatalf("expect %v but get: %v ", v, s)
			}
		}, "http://localhost:8080/test/chunk", nil)
		if err != nil {
			t.Fatal(err)
		}
		if i != 5 {
			t.Fatal("expect 5 but get: ", i)
		}
	})
	t.Run("get_string", func(t *testing.T) {
		c := New(SetTimeout(0))
		_, err := c.Get(func(s string) {
			fmt.Printf("%s\n", s)
		}, "http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
	})
}
