// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"context"
	"fmt"
	"github.com/xfali/restclient"
	"github.com/xfali/restclient/request"
	"github.com/xfali/restclient/restutil"
	"io"
	"net/http"
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
	http.HandleFunc("/auth", func(writer http.ResponseWriter, request *http.Request) {
		v := request.Header.Get(restutil.HeaderAuthorization)
		fmt.Println(v)
		writer.Header().Set(restutil.HeaderContentType, "application/json")
		writer.Write([]byte(`{ "result":["hello", "world"]}`))
	})

	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			writer.Write([]byte(`Get: `))
			body := request.Body
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		case http.MethodPost:
			writer.Write([]byte(`Post: `))
			body := request.Body
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		case http.MethodPut:
			writer.Write([]byte(`Put: `))
			body := request.Body
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		case http.MethodDelete:
			writer.Write([]byte(`Delete: `))
			body := request.Body
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		case http.MethodHead:
			writer.Write([]byte(`Head: `))
			body := request.Body
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		case http.MethodOptions:
			writer.Write([]byte(`Options: `))
			body := request.Body
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		case http.MethodPatch:
			writer.Write([]byte(`Patch: `))
			body := request.Body
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		}
	})

	server := &http.Server{Addr: ":8080", Handler: nil}
	go server.ListenAndServe()

	<-time.NewTimer(shutdown).C

	server.Shutdown(context.Background())
}

func TestRequest(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		ret := ""
		s, err := request.New().Get("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Post", func(t *testing.T) {
		ret := ""
		s, err := request.New().SetBody("hello world").Post("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Put", func(t *testing.T) {
		ret := ""
		s, err := request.New().Put("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Delete", func(t *testing.T) {
		ret := ""
		s, err := request.New().Delete("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Patch", func(t *testing.T) {
		ret := ""
		s, err := request.New().Patch("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Options", func(t *testing.T) {
		ret := ""
		s, err := request.New().Options("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Head", func(t *testing.T) {
		ret := ""
		s, err := request.New().Head("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})
}

func TestRequest2(t *testing.T) {
	t.Run("Get no result", func(t *testing.T) {
		s, err := request.New().Get("http://localhost:8080/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
	})

	t.Run("Post with body", func(t *testing.T) {
		ret := ""
		s, err := request.New(request.SetBody("This is Post body!")).Post("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Put with body", func(t *testing.T) {
		ret := ""
		s, err := request.New(request.SetBody("This is Put body!")).Put("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Delete", func(t *testing.T) {
		ret := ""
		s, err := request.New(request.SetBody("this is delete body")).Delete("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Patch", func(t *testing.T) {
		ret := ""
		s, err := request.New().Patch("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Options", func(t *testing.T) {
		ret := ""
		s, err := request.New().Options("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})

	t.Run("Head", func(t *testing.T) {
		ret := ""
		s, err := request.New().Head("http://localhost:8080/test", &ret)
		if err != nil {
			t.Fatal(err)
		}
		if s != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(s)
		}
		t.Log(ret)
	})
}

func TestConv(t *testing.T) {
	client := restclient.New(restclient.SetAutoAccept(restclient.AcceptUserOnly))
	v := &[]byte{}
	status, err := client.Exchange(v, "http://localhost:8080/test", http.MethodPost, map[string]interface{}{"Accept":"application/json"}, []byte("123"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(status)
	t.Log(string(*v))
}