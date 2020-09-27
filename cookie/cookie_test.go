// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package cookie

import (
	"net/http"
	"testing"
	"time"
)

func TestCookie(t *testing.T) {
	cache := NewCache()
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
