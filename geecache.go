// 负责与外部交互，控制缓存存储和获取的主流程
package geecache

import (
	"fmt"
	"log"
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
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		//必然报错
		return ByteView{}, fmt.Errorf("key is required")
	}
	//如果存在 则直接返回缓存中的值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

// 通过load 调用getLocally
// 从其他节点获取,去找数据源
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	//通过回调函数从数据源中获取数据
	bytes, err := g.getter.Get(key)

	if err != nil {
		return ByteView{}, err
	}

	//为了不对源数据做出影响
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 添加到缓存中
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
