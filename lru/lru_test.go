package lru

import (
	"reflect"
	"testing"
)

// 这里还需要这样才能用？
type String string

func (d String) Len() int {
	return len(d)
}

// 测试用的参数只有一个 所以是tessting,T
// 基准测试的参数是 testing.B
// TestMain的参数是 testing.M
func TestGet(t *testing.T) {
	//这里设置为0 我们用来测试是否当内存超过了定值
	//是否会触发移除操作
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	//测试删除的情况
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache != 1234 failed")
	}
	//测试是否能查不到的情况
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// 测试是否会删除无用的节点
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

// 测试回调函数
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
