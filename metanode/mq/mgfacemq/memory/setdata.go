package memory

import (
	"encoding/json"
	"time"

	"com.mgface.disobj/common"
)

// Set 将磁盘数据加载到内存中，只要关注一个问题就是：内存中是否已经存在了该数据
func (cache *MemoryStore) Set(key string, value []byte) error {
	// 假如为非心跳数据，那么则发送消息给消息队列
	if !(key == common.DataNodeKey || key == common.ApiNodeKey || key == common.MetaNodeKey) && !cache.LoadingSnapshot {
		// todo 传递string而不是[]byte，主要是json.Marshal对字节数组回进行base64，导致字符串足够长,不利于后续压缩算法
		cache.Msgs <- common.RecMsg{Key: key, Val: string(value)}
	}

	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()
	// 查看内存中是否存在数据
	cachedValues, exist := cache.Datas[key]
	if key == common.DataNodeKey || key == common.ApiNodeKey || key == common.MetaNodeKey {
		v := common.MetaValue{RealNodeValue: string(value), Created: time.Now()}
		data := make([]common.MetaValue, 0)
		data = append(data, v)
		if exist {
			var forceTxData []common.MetaValue
			cachedValuesByte, _ := json.Marshal(cachedValues)
			err := json.Unmarshal(cachedValuesByte, &forceTxData)
			if err != nil {
				return err
			}
			// 便利缓存数据，如果存在同样datanode IP的数据，那么移除老数据，添加新数据
			for _, val := range forceTxData {
				if val.RealNodeValue != v.RealNodeValue {
					data = append(data, val)
				}
			}
		}
		cache.Datas[key] = data
		return nil
	} else if key == common.MetaDataKey {
		// 先将磁盘存储的数据，解析出来
		var valueDKV common.DataKeyValue
		err := json.Unmarshal(value, &valueDKV)
		if err != nil {
			return err
		}

		if exist {
			// 解析内存数据
			var loadCachedValues map[string][]common.Datadigest
			loadCachedValuesByte, _ := json.Marshal(cachedValues)
			err = json.Unmarshal(loadCachedValuesByte, &loadCachedValues)
			if err != nil {
				return err
			}

			// 遍历磁盘数据，将磁盘数据添加到内存中
			for k, v := range valueDKV.Data {
				// dataDigests是历史版本的摘要数据， key是对象名
				// 如果内存中存在该对象的数据，则更新内存数据
				if dataDigests, ok := loadCachedValues[k]; ok {
					var max int64 = 0
					// 找到历史最大版本
					for _, v := range dataDigests {
						if v.Version > max {
							max = v.Version
						}
					}
					// 修改当前附加的版本号
					// 这啥意思？为什么要选择第一个版本，使版本号+1？
					v[0].Version = max + 1
					loadCachedValues[k] = append(dataDigests, v...)
				} else { // 对象数据还没加载到内存中，则直接添加
					loadCachedValues[k] = v
				}
				cache.Datas[key] = loadCachedValues
				return nil
			}
		}
		// 内存中不存在metadata
		cache.Datas[key] = valueDKV.Data
	} else if key == common.FileCrcKey {
		var fileCrc common.FileCRC
		err := json.Unmarshal(value, &fileCrc)
		if err != nil {
			return err
		}
		if exist {
			// 内存中存在filecrc，则加载内存数据
			var loadCachedValues map[string][]common.SharedData
			cachedValuesByte, _ := json.Marshal(cachedValues)
			err := json.Unmarshal(cachedValuesByte, &loadCachedValues)
			if err != nil {
				return err
			}
			// 比较磁盘数据和内存数据，如果磁盘数据版本号大于内存数据版本号，则更新内存数据
			// 遍历每个对象的分片meta数据
			for k, v := range fileCrc.Data {
				// 此对象存在内存中
				if vdata, ok := loadCachedValues[k]; ok { //假如存在，那么数据添加
					loadCachedValues[k] = append(vdata, v...)
				} else {
					loadCachedValues[k] = v
				}
			}
			cache.Datas[key] = loadCachedValues
			return nil
		}
		//假如数据不存在
		cache.Datas[key] = fileCrc.Data
	} else {
		v := common.MetaValue{RealNodeValue: string(value), Created: time.Now()}
		cache.Datas[key] = v
	}
	return nil
}
