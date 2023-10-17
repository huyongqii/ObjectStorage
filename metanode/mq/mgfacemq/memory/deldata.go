package memory

import . "com.mgface.disobj/common"

func (cache *MemoryStore) Del(key string, value []byte) error {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()
	mvalues, exit := cache.Datas[key]
	if key == "search" {
		if exit {
			forceTxData := mvalues.([]MetaValue)
			//假如数据长度为1，说明是单节点数据交互
			if len(forceTxData) == 1 {
				delete(cache.Datas, key)
				return nil
			}
			data := make([]MetaValue, 0)
			for _, val := range forceTxData {
				if val.RealNodeValue != string(value) {
					data = append(data, val)
				}
			}
			cache.Datas[key] = data
		}
	} else {
		delete(cache.Datas, key)
	}
	return nil
}
