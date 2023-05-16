// 进行并发控制
package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

// 那这里是还要包装一个序列化吗
func (c *cache) add(key string, value ByteView) {
	//上锁
	c.mu.Lock()
	defer c.mu.Unlock()

	//创建缓存
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	//进行分配
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		//这里还进行了一个类型断言
		return v.(ByteView), ok
	}
	return
}
