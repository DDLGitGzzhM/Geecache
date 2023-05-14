package lru

import (
	"container/list"
)

type Cache struct {
	//允许使用的最大内存
	maxBytes int64
	//已经使用的内存
	nbytes int64
	ll     *list.List
	//值对应双向链表的指针
	cache map[string]*list.Element
	//某条记录被移除时候的回调函数
	OnEvicted func(key string, value Value)
}

// 双向链表的数据格式
type entry struct {
	key   string
	value Value
}
type Value interface {
	//返回占用内存的大小
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找功能
// 第一步从字典中找到对应的双向链表节点
// 将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 这里约定, Front是队尾
		c.ll.MoveToFront(ele)
		//这是一个类型断言
		kv := ele.Value.(*entry)
		//如果找到则返回这个值
		return kv.value, true
	}
	return
}

// 删除
func (c *Cache) ReMoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		//这是一个类型断言
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增/修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		//存在则修改亲爱
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.ReMoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
