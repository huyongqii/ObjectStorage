package memory

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// AssignThis 设置自身数据
func (cache *MemoryStore) AssignThis(data interface{}) {

	//判断data是否是字符串类型
	var dx map[string]interface{}
	if err := json.Unmarshal(data.([]byte), &dx); err != nil {
		log.Warn("反序列化数据异常.")
	}

	if len(dx) > 0 {
		cache.Mutex.Lock()
		defer cache.Mutex.Unlock()
		for k, v := range dx {
			cache.Datas[k] = v
		}
	}
}
