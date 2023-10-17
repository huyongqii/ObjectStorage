package memory

import (
	"encoding/json"
	"fmt"
	"time"
)

// Show 显示内存数据信息到控制台
func (cache *MemoryStore) Show() string {
	for {
		cache.Mutex.RLock()
		data, err := json.MarshalIndent(cache.Datas, "", "\t")
		//data, err := json.Marshal(cache.Datas)
		cache.Mutex.RUnlock()
		if err != nil {
			return "发生错误信息:" + err.Error()
		}
		fmt.Println("内存数据:", string(data))
		time.Sleep(5 * time.Second)
	}
}
