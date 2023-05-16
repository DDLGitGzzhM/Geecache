// 用来表示缓存值
package geecache

type ByteView struct {
	//为了支持任意数据类型存储 所以是byte
	b []byte
}

// 实现Value接口
func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// 返回一个拷贝 防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
