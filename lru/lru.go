package lru

import "container/list"

// LRU结构
type Cache struct {
	maxBytes  int64                         // 允许的最大内存
	nbytes    int64                         // 当前已使用的内存
	ll        *list.List                    // 双向链表	用来移动，删除和增加
	cache     map[string]*list.Element      // 字典存储	,, 存储链表类型
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为 nil。
}
// 记录的值	, 注意value类型，值类型
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {	// 根据Key获取值
		c.ll.MoveToFront(ele)	// 使用了，放在前面
		kv := ele.Value.(*entry)	// 类型断言，获取链表的值	,kv为entry类型
		return kv.value, true	// 再取出值
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()	// 获取链表最后的元素
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {	// 判断当前的值是否存储过
		c.ll.MoveToFront(ele)	// 如果存储过，移动到链表前面
		kv := ele.Value.(*entry)	//类型断言
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())	// 更新占用的内存,减去之前的加上现在的
		kv.value = value	// 将原来的值替换
	} else {
		ele := c.ll.PushFront(&entry{key, value})	// 没有存储过，直接放入前面；返回的是一个句柄 *Element
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())	// 添加至内存
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()	// 删除最前面的元素
	}
}