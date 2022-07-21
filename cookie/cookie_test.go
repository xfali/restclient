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

package cookie

import (
	"context"
	"fmt"
	"github.com/xfali/restclient/v2/request"
	"github.com/xfali/restclient/v2/restutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestCookieSet(t *testing.T) {
	cache := NewCache()
	defer cache.Close()

	cache.Set("https://127.0.0.1:8080", &http.Cookie{
		Name:  "hello",
		Value: "first",
	})

	cache.Set("https://127.0.0.1:8080/test", &http.Cookie{
		Name:  "hello",
		Value: "second",
	})

	cache.Set("https://127.0.0.1:8080/test/key", &http.Cookie{
		Name:  "hello",
		Value: "third",
	})

	cookies := cache.Get("https://127.0.0.1:8080/test/key")
	if len(cookies) != 2 {
		t.Fatal("expect 2 but get ", len(cookies))
	} else {
		for _, v := range cookies {
			t.Log(*v)
		}
	}
}

func TestCookie(t *testing.T) {
	cache := NewCache()
	cache.AutoPurge()
	defer cache.Close()

	cache.Set("https://127.0.0.1:8080", &http.Cookie{
		Name:   "hello",
		Value:  "world",
		Path:   "/",
		MaxAge: 2,
	})

	cache.Set("https://127.0.0.1:8080", &http.Cookie{
		Name:   "hello",
		Value:  "first",
		Path:   "/test",
		MaxAge: 1,
	})

	cache.Set("https://127.0.0.1:8080", &http.Cookie{
		Name:   "hello",
		Value:  "second",
		Path:   "/test",
		MaxAge: 1,
	})

	cookies := cache.Get("https://127.0.0.1:8080")
	if len(cookies) != 1 {
		t.Fatal("expect 3 but get ", len(cookies))
	}

	cookies = cache.Get("https://127.0.0.1:8080/test/key")
	if len(cookies) != 2 {
		t.Fatal("expect 2 but get ", len(cookies))
	} else {
		for _, v := range cookies {
			t.Log(*v)
		}
	}

	time.Sleep(time.Second)
	cookies = cache.Get("https://127.0.0.1:8080/test/key")
	if len(cookies) != 1 {
		t.Fatal("expect 1 but get ", len(cookies))
	} else {
		for _, v := range cookies {
			t.Log(*v)
		}
	}

	time.Sleep(time.Second)
	cookies = cache.Get("https://127.0.0.1:8080/test/key")
	if len(cookies) != 0 {
		t.Fatal("expect 0 but get ", len(cookies))
	} else {
		for _, v := range cookies {
			t.Log(*v)
		}
	}

	cache.Set("https://127.0.0.1:8080", &http.Cookie{
		Name:  "hello",
		Value: "third",
		Path:  "/test",
	})
	cookies = cache.Get("https://127.0.0.1:8080/test/key")
	if len(cookies) != 1 {
		t.Fatal("expect 1 but get ", len(cookies))
	} else {
		for _, v := range cookies {
			t.Log(*v)
		}
	}

	cache.Set("https://127.0.0.1:8080", &http.Cookie{
		Name:   "hello",
		Value:  "third",
		Path:   "/test",
		MaxAge: -1,
	})
	cookies = cache.Get("https://127.0.0.1:8080/test/key")
	if len(cookies) != 0 {
		t.Fatal("expect 0 but get ", len(cookies))
	} else {
		for _, v := range cookies {
			t.Log(*v)
		}
	}
}

func startHttpServer(shutdown time.Duration) {
	http.HandleFunc("/cookie", func(writer http.ResponseWriter, request *http.Request) {
		if c, err := request.Cookie("SESSION"); err != nil || c.Value == "" {
			v := request.Header.Get(restutil.HeaderAuthorization)
			fmt.Println(v)
			if v == "" {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				http.SetCookie(writer, &http.Cookie{Name: "SESSION", Value: v, MaxAge: 0})
			}
		}

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

func TestCookieRequest(t *testing.T) {
	go startHttpServer(100 * time.Second)
	time.Sleep(time.Second)

	t.Run("get string", func(t *testing.T) {
		cache := NewCache()
		defer cache.Close()
		c := restclient.New(restclient.AddFilter(cache.Filter))
		err := c.Exchange("http://localhost:8080/cookie", request.WithResult(func(s string) {
			fmt.Printf("%s\n", s)
		}))
		if err == nil {
			t.Fatal(err)
		}
		if err.StatusCode() != http.StatusUnauthorized {
			t.Fatal("must 401")
		}
		err = c.Exchange("http://localhost:8080/cookie", request.WithResult(func(s string) {
			fmt.Printf("%s\n", s)
		}), request.WithRequestHeader(http.Header{
			restutil.HeaderAuthorization: []string{"123"},
		}))

		if err != nil {
			t.Fatal(err)
		}
		err = c.Exchange("http://localhost:8080/cookie", request.WithResult(func(s string) {
			fmt.Printf("%s\n", s)
		}))

		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestSlice(t *testing.T) {
	sl := []int{}
	sl = append(sl[:0], sl[0:]...)
	t.Log(sl)

	sl = []int{1, 2}
	t.Log(sl[:1])
	sl = append(sl[:1], sl[1+1:]...)
	t.Log(sl)
}
