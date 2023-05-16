package geecache

import (
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom": "630",
	"Jac": "589",
	"Sam": "567",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")

	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
