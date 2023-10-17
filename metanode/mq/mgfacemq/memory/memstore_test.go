package memory

import (
	"sync"
	"testing"
)

func TestMemoryCache_Set(t *testing.T) {
	var m = &MemoryStore{Mutex: sync.RWMutex{}, Datas: make(map[string]interface{})}
	m.Set("123", []byte("admin1"))
	m.Set("123", []byte("admin2"))
	v, _ := m.Get("123", nil)
	t.Log(v)
}
