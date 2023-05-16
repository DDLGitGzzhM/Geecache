// 负责与外部交互，控制缓存存储和获取的主流程
package geecache

import (
	. "geecache/byteview"
	"sync"
)

/*
如果缓存不在,我们应该从数据源(文件、数据库等)
并添加到缓存中
*/

type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义回调函数接口
type GetterFunc func(key string) ([]byte, error)

// Get的回调 ,这里为什么不反悔的是byteview不懂
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/*
下面我们讲
处理缓存
处理回调
的两个部分封装成一共Group由Group统一控制
*/

// 缓存的命名空间
type Group struct {
	name   string
	getter Getter
	//缓存为命中获取的数据回调
	mainCache cache
}

var (
	//因为map并发不安全,统一固定
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	//规定必然要有getter
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	//修改进行读写控制
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err  := g.getter.Get(key)

	if err != nil {
		return ByteView{} , err
	}

	value := ByteView{b: }
}
