package geecache

import (
	"lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex	// 添加互斥锁
	lru        *lru.Cache	// LRU结构
	cacheBytes int64
}

// add 添加
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	/*
	，如果等于 nil 再创建实例。
	这种方法称之为延迟初始化，一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。
	主要用于提高性能，并减少程序内存要求。
	*/
	if c.lru == nil {	// 如果没有缓存列表
		c.lru = lru.New(c.cacheBytes, nil)	// 新建并初始化一个
	}
	c.lru.Add(key, value)		// 赋值
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {	// 如果没有缓存列表
		return
	}

	if v, ok := c.lru.Get(key); ok {	// 如果有缓存列表并取到值
		return v.(ByteView), ok	// 返回
	}

	return
}
