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

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xfali/restclient/v2"
	"github.com/xfali/restclient/v2/filter"
	"github.com/xfali/restclient/v2/request"
	"github.com/xfali/restclient/v2/restutil"
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

type testStruct struct {
	Id         int64
	Name       string
	Value      float64
	CreateTime time.Time
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

	http.HandleFunc("/struct", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			body := request.Body
			if body != nil {
				defer body.Close()
			}
			d, _ := json.Marshal(testStruct{
				Id:         1,
				Name:       "test",
				Value:      3.1415926,
				CreateTime: time.Now(),
			})
			writer.Header().Set(restutil.HeaderContentType, restclient.MediaTypeJson)
			writer.Write(d)
			break
		case http.MethodPost:
			body := request.Body
			writer.Header().Set(restutil.HeaderContentType, restclient.MediaTypeJson)
			if body != nil {
				defer body.Close()
				io.Copy(writer, body)
			}
			break
		}
	})

	http.HandleFunc("/error", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			body := request.Body
			if body != nil {
				defer body.Close()
			}
			d, _ := json.Marshal(testStruct{
				Id:         1,
				Name:       "test",
				Value:      3.1415926,
				CreateTime: time.Now(),
			})
			writer.Header().Set(restutil.HeaderContentType, restclient.MediaTypeJson)
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write(d)
			break
		case http.MethodPost:
			body := request.Body
			writer.Header().Set(restutil.HeaderContentType, restclient.MediaTypeJson)
			writer.WriteHeader(http.StatusBadRequest)
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
			t.Log(resp.StatusCode)
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
			t.Log(resp.StatusCode)
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
			t.Log(resp.StatusCode)
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
			t.Log(resp.StatusCode)
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
			t.Log(resp.StatusCode)
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
			t.Log(resp.StatusCode)
		}
		t.Log(ret)
	})

	t.Run("Get param", func(t *testing.T) {
		ret := ""
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/test",
			restclient.NewRequest().Method(http.MethodPost).Response(resp, false).RequestBody("hello world").Build())
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200 ", resp.StatusCode)
		} else {
			t.Log(resp.StatusCode)
		}
		t.Log(ret)
	})
}

func TestStruct(t *testing.T) {
	client := restclient.New(restclient.AddIFilter(filter.NewLog(xlog.GetLogger(), "")))
	t.Run("Get", func(t *testing.T) {
		ret := testStruct{}
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/struct",
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
		if ret.Id != 1 {
			t.Fatal("expect id 1 but get ", ret.Id)
		}
		t.Log(ret)
	})

	t.Run("Get func", func(t *testing.T) {
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/struct",
			request.WithResult(func(ret testStruct) {
				if ret.Id != 1 {
					t.Fatal("expect id 1 but get ", ret.Id)
				}
				t.Log(ret)
			}),
			request.WithResponse(resp, false))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(resp.StatusCode)
		}
	})

	t.Run("Post", func(t *testing.T) {
		ret := testStruct{
			Id:         2,
			Name:       "test2",
			Value:      1.0,
			CreateTime: time.Now(),
		}
		resp := new(http.Response)
		err := client.Exchange("http://localhost:8080/struct",
			request.MethodPost(),
			request.WithResult(&ret),
			request.WithResponse(resp, false),
			request.WithRequestBody(ret))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(resp.StatusCode)
		}
		if ret.Id != 2 {
			t.Fatal("expect id 2 but get ", ret.Id)
		}
		t.Log(ret)
	})
}

func TestErrorStruct(t *testing.T) {
	client := restclient.New(restclient.AddIFilter(filter.NewLog(xlog.GetLogger(), "")))
	t.Run("Get", func(t *testing.T) {
		ret := testStruct{}
		err := client.Exchange("http://localhost:8080/error",
			request.WithResult(&ret))
		if err == nil {
			t.Fatal(err)
		} else {
			t.Log(err)
		}
		if err.StatusCode() == http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(err.StatusCode())
		}
		if ret.Id != 1 {
			t.Fatal("expect id 1 but get ", ret.Id)
		}
		t.Log(ret)
	})

	t.Run("Get not found", func(t *testing.T) {
		ret := testStruct{}
		err := client.Exchange("http://localhost:8080/404",
			request.WithResult(&ret))
		if err == nil {
			t.Fatal(err)
		} else {
			t.Log(err)
		}
		if err.StatusCode() != http.StatusNotFound {
			t.Fatal("not 404")
		} else {
			t.Log(err.StatusCode())
		}
		if ret.Id == 1 {
			t.Fatal("expect id 0 but get ", ret.Id)
		}
		t.Log(ret)
	})

	t.Run("Get func", func(t *testing.T) {
		err := client.Exchange("http://localhost:8080/error",
			request.WithResult(func(ret testStruct) {
				if ret.Id != 1 {
					t.Fatal("expect id 1 but get ", ret.Id)
				}
				t.Log(ret)
			}))
		if err == nil {
			t.Fatal(err)
		} else {
			t.Log(err)
		}
		if err.StatusCode() == http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(err.StatusCode())
		}
	})

	t.Run("Post", func(t *testing.T) {
		ret := testStruct{
			Id:         2,
			Name:       "test2",
			Value:      1.0,
			CreateTime: time.Now(),
		}
		err := client.Exchange("http://localhost:8080/error",
			request.MethodPost(),
			request.WithResult(&ret),
			request.WithRequestBody(ret))
		if err == nil {
			t.Fatal(err)
		} else {
			t.Log(err)
		}
		if err.StatusCode() == http.StatusOK {
			t.Fatal("not 200")
		} else {
			t.Log(err.StatusCode())
		}
		if ret.Id != 2 {
			t.Fatal("expect id 2 but get ", ret.Id)
		}
		t.Log(ret)
	})
}

func TestUrlBuilder(t *testing.T) {
	b := restclient.NewUrlBuilder("x/:a/tt/:b?")
	b.PathVariable("a", "1")
	b.PathVariable("b", 2)
	b.QueryVariable("c", 100)
	b.QueryVariable("d", 1.1)
	url := b.Build()
	if url != "x/1/tt/2?c=100&d=1.1" {
		t.Fatal(url)
	}

	t.Log(url)
}
