// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package cookie

import (
	"fmt"
	"github.com/xfali/restclient/v2/filter"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	// 默认清理时间间隔
	DefaultPurgeInterval = 50 * time.Millisecond
	DefaultDepth         = 5
)

type Cache struct {
	purgeInterval time.Duration
	depth         int

	db   map[string]*[]*http.Cookie
	stop chan bool
	lock sync.Locker
}

type cookieCacheOpt func(*Cache)

// 创建一个cookie缓存，注意需要Close
func NewCache(opts ...cookieCacheOpt) *Cache {
	ret := &Cache{
		purgeInterval: DefaultPurgeInterval,
		depth:         DefaultDepth,
		db:            map[string]*[]*http.Cookie{},
		stop:          make(chan bool),
		lock:          &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(ret)
	}

	ret.run()
	return ret
}

func OptSetPurgeInterval(interval time.Duration) cookieCacheOpt {
	return func(cache *Cache) {
		cache.purgeInterval = interval
	}
}

func OptSetDepth(depth int) cookieCacheOpt {
	return func(cache *Cache) {
		cache.depth = depth
	}
}

func OptSetLocker(lock sync.Locker) cookieCacheOpt {
	return func(cache *Cache) {
		cache.lock = lock
	}
}

type pos struct {
	k string
	i int
}

func (dm *Cache) purge() {
	dm.lock.Lock()
	defer dm.lock.Unlock()

	now := time.Now()
	var cookie2Del []pos
	for k, cookies := range dm.db {
		for i, v := range *cookies {
			if v.Expires.IsZero() {
				continue
			}
			if !v.Expires.After(now) {
				pos := pos{k, i}
				cookie2Del = append(cookie2Del, pos)
			}
		}
	}
	dm.delete(cookie2Del)
}

func (dm *Cache) delete(cookie2Del []pos) {
	for _, v := range cookie2Del {
		fmt.Println(v, "deleted")
		cs := dm.db[v.k]
		*cs = append((*cs)[:v.i], (*cs)[v.i+1:]...)
		if len(*cs) == 0 {
			delete(dm.db, v.k)
		}
	}
}

//初始化并开启回收线程，必须调用
func (dm *Cache) run() {
	if dm.purgeInterval <= 0 {
		dm.purgeInterval = 0
	}

	go func() {
		if dm.purgeInterval > 0 {
			timer := time.NewTicker(dm.purgeInterval)
			defer timer.Stop()
			for {
				select {
				case <-dm.stop:
					return
				case <-timer.C:
					dm.purge()
				}
			}
		} else {
			for {
				select {
				case <-dm.stop:
					return
				default:
				}
				dm.purge()

				runtime.Gosched()
			}
		}
	}()
}

//关闭
func (dm *Cache) Close() error {
	close(dm.stop)
	return nil
}

// 设置一个值，含过期时间
// 如果expireIn设置为-1，则永不过期
func (dm *Cache) Set(path string, cookie *http.Cookie) error {
	if cookie == nil {
		return nil
	}

	uri, err := url.Parse(path)
	if err != nil {
		return err
	}
	domain := cookie.Domain
	if domain == "" {
		domain = uri.Hostname()
	}

	checkCookiePath(uri.Path, cookie)

	if cookie.MaxAge > 0 {
		if cookie.Expires.IsZero() {
			cookie.Expires = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
		}
	}

	key := domain + cookie.Path

	dm.lock.Lock()
	defer dm.lock.Unlock()

	// delete cookie
	if cookie.MaxAge < 0 {
		if v, ok := dm.db[key]; ok {
			for i, old := range *v {
				if cookie.Name == old.Name {
					*v = append((*v)[:i], (*v)[i+1:]...)
					return nil
				}
			}
		}
		return nil
	}

	if v, ok := dm.db[key]; ok {
		for i, old := range *v {
			// found and return
			if old.Name == cookie.Name {
				(*v)[i] = cookie
				return nil
			}
		}
		*v = append(*v, cookie)
	} else {
		dm.db[key] = &[]*http.Cookie{cookie}
	}

	return nil
}

func checkCookiePath(path string, cookie *http.Cookie) {
	if cookie.Path == "" {
		if path == "" || path== "/" {
			cookie.Path = "/"
		} else {
			index := strings.LastIndex(path, "/")
			if index != -1 {
				cookie.Path = path[:index]
			}
			if cookie.Path == "" {
				cookie.Path = "/"
			}
		}
	}
}

//根据key获取value
func (dm *Cache) Get(path string) []*http.Cookie {
	uri, err := url.Parse(path)
	if err != nil {
		return nil
	}

	dm.lock.Lock()
	defer dm.lock.Unlock()

	found := 0
	var ret []*http.Cookie
	domain := uri.Hostname()
	path = uri.Path
	dm.foundAndPurge(domain+"/", &ret)
	for i := 1; i < len(path); i++ {
		if path[i] == '/' {
			dm.foundAndPurge(domain+path[:i], &ret)
			found++
			if found >= dm.depth {
				break
			}
		}
	}

	return ret
}

func (dm *Cache) foundAndPurge(key string, ret *[]*http.Cookie) bool {
	var cookie2Del []pos
	found := false
	if v, ok := dm.db[key]; ok {
		now := time.Now()
		for index, cookie := range *v {
			if cookie.Expires.IsZero() {
				found = true
				*ret = append(*ret, cookie)
			} else if !cookie.Expires.After(now) {
				cookie2Del = append(cookie2Del, pos{key, index})
			} else {
				found = true
				*ret = append(*ret, cookie)
			}
		}
	}
	dm.delete(cookie2Del)
	return found
}

func (dm *Cache) Filter(request *http.Request, fc filter.FilterChain) (*http.Response, error) {
	path := request.URL.String()
	for _, v := range dm.Get(path) {
		request.AddCookie(v)
	}
	resp, err := fc.Filter(request)
	if resp != nil {
		for _, v := range resp.Cookies() {
			dm.Set(path, v)
		}
	}

	return resp, err
}
