// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"context"
	"fmt"
	"github.com/xfali/restclient/v2/request"
	"github.com/xfali/restclient/restutil"
	"github.com/xfali/restclient/v2"
	"github.com/xfali/restclient/v2/filter"
	"github.com/xfali/xlog"
	"io"
	"net/http"
	"testing"
	"time"
)

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
	client := restclient.New(restclient.AddIFilter(filter.NewLog(xlog.GetLogger(), "")))
	t.Run("Get", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			request.WithResult(&ret),
			request.WithResponse(resp, false))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Post", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			request.MethodPost(),
			request.WithResult(&ret),
			request.WithResponse(resp, false),
			request.WithRequestBody("hello world"))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log( resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Put", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			request.MethodPut(),
			request.WithResult(&ret),
			request.WithResponse(resp, false),
			request.WithRequestBody("hello world"))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log( resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Delete", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			request.MethodDelete(),
			request.WithResult(&ret),
			request.WithResponse(resp, false))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log( resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Patch", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			request.MethodPatch(),
			request.WithResult(&ret),
			request.WithResponse(resp, false),
			request.WithRequestBody("hello world"))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log( resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Options", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			request.MethodOptions(),
			request.WithResult(&ret),
			request.WithResponse(resp, false))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log( resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Head", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			request.MethodHead(),
			request.WithResult(&ret),
			request.WithResponse(resp, false))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log( resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Get param", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			restclient.NewRequest().WitMethod(http.MethodPost).WithResponse(resp, false).WithRequestBody("hello world").Build())
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200 ", resp.StatusCode)
		} else {
			t.Log( resp.StatusCode)
		}
		t.Log(ret)
	})
}
