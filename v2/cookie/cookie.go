// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package cookie

import (
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
	DefaultPurgeInterval = 15 * time.Second
	DefaultDepth         = 5
)

type Cache interface {
	// 设置Cookie
	Set(path string, cookie *http.Cookie) error

	// 获得Cookie，会触发自动回收同域过期Cookie
	Get(path string) []*http.Cookie

	// 回收过期Cookie
	Purge()

	// 自动回收过期Cookie
	// 注意只能调用一次，且需要调用Close关闭
	AutoPurge()

	// 关闭缓存
	// 注意只能调用一次，AutoPurge调用后如不再需要Cache则必须调用此方法
	Close() error
}

type defaultCache struct {
	purgeInterval time.Duration
	depth         int

	db   map[string]*[]*http.Cookie
	stop chan struct{}
	lock sync.Locker
}

type cookieCacheOpt func(*defaultCache)

// 创建一个cookie缓存
// 注意Cache不会自动回收过期Cookie，只有在查询Cookie时尝试回收，可以通过调用Purge方法手工回收
// 可调用AutoPurge方法开启自动回收
// 注意：AutoPurge只可调用一次，该方法会开启一个协程，通过Close方法退出（也仅可调用一次）。
func NewCache(opts ...cookieCacheOpt) *defaultCache {
	ret := &defaultCache{
		purgeInterval: DefaultPurgeInterval,
		depth:         DefaultDepth,
		db:            map[string]*[]*http.Cookie{},
		stop:          make(chan struct{}),
		lock:          &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(ret)
	}

	return ret
}

func OptSetPurgeInterval(interval time.Duration) cookieCacheOpt {
	return func(cache *defaultCache) {
		cache.purgeInterval = interval
	}
}

func OptSetDepth(depth int) cookieCacheOpt {
	return func(cache *defaultCache) {
		cache.depth = depth
	}
}

func OptSetLocker(lock sync.Locker) cookieCacheOpt {
	return func(cache *defaultCache) {
		cache.lock = lock
	}
}

type pos struct {
	k string
	i int
}

func (dm *defaultCache) Purge() {
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

func (dm *defaultCache) delete(cookie2Del []pos) {
	for _, v := range cookie2Del {
		//fmt.Println(v, "deleted")
		cs := dm.db[v.k]
		*cs = append((*cs)[:v.i], (*cs)[v.i+1:]...)
		if len(*cs) == 0 {
			delete(dm.db, v.k)
		}
	}
}

//初始化并开启回收协程，必须调用
func (dm *defaultCache) AutoPurge() {
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
					dm.Purge()
				}
			}
		} else {
			for {
				select {
				case <-dm.stop:
					return
				default:
				}
				dm.Purge()

				runtime.Gosched()
			}
		}
	}()
}

//关闭
func (dm *defaultCache) Close() error {
	close(dm.stop)
	return nil
}

// 设置一个值，含过期时间
// 如果expireIn设置为-1，则永不过期
func (dm *defaultCache) Set(path string, cookie *http.Cookie) error {
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
func (dm *defaultCache) Get(path string) []*http.Cookie {
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

func (dm *defaultCache) foundAndPurge(key string, ret *[]*http.Cookie) bool {
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

func (dm *defaultCache) Filter(request *http.Request, fc filter.FilterChain) (*http.Response, error) {
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

func  Filter(cache Cache) filter.Filter {
	return func(request *http.Request, fc filter.FilterChain) (response *http.Response, e error) {
		path := request.URL.String()
		for _, v := range cache.Get(path) {
			request.AddCookie(v)
		}
		resp, err := fc.Filter(request)
		if resp != nil {
			for _, v := range resp.Cookies() {
				cache.Set(path, v)
			}
		}
		return resp, err
	}
}
