// 负责与外部交互，控制缓存存储和获取的主流程
package geecache

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
