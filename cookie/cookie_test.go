// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package cookie

import (
	"context"
	"fmt"
	"github.com/xfali/restclient"
	"github.com/xfali/restclient/restutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestCookieSet(t *testing.T) {
	cache := NewCache()
	defer cache.Close()

	cache.Set("https://127.0.0.1:8080", &http.Cookie{
		Name:   "hello",
		Value:  "first",
	})

	cache.Set("https://127.0.0.1:8080/test", &http.Cookie{
		Name:   "hello",
		Value:  "second",
	})

	cache.Set("https://127.0.0.1:8080/test/key", &http.Cookie{
		Name:   "hello",
		Value:  "third",
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
		status, err := c.Get(func(s string) {
			fmt.Printf("%s\n", s)
		}, "http://localhost:8080/cookie", nil)
		if err != nil {
			t.Fatal(err)
		}
		if status != http.StatusUnauthorized {
			t.Fatal("must 401")
		}
		status, err = c.Get(func(s string) {
			fmt.Printf("%s\n", s)
		}, "http://localhost:8080/cookie", map[string]interface{}{
			restutil.HeaderAuthorization: "123",
		})

		if err != nil {
			t.Fatal(err)
		}
		if status != http.StatusOK {
			t.Fatal("must 200")
		}
		status, err = c.Get(func(s string) {
			fmt.Printf("%s\n", s)
		}, "http://localhost:8080/cookie", nil)
		if err != nil {
			t.Fatal(err)
		}
		if status != http.StatusOK {
			t.Fatal("must 200")
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
