package memory

import (
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
)

func (cache *MemoryStore) Get(key string, value []byte) (interface{}, error) {
	cache.Mutex.RLock()
	defer cache.Mutex.RUnlock()
	data := make([]interface{}, 0)
	if key == common.DataNodeKey || key == common.ApiNodeKey || key == common.MetaNodeKey {
		nodeValues, exist := cache.Datas[key]
		if exist {
			log.Debug("操作:", key, ",value:", string(value), ",当前类型：", reflect.TypeOf(nodeValues))

			var forceTxData []common.MetaValue
			nodeValuesByte, _ := json.Marshal(nodeValues)
			err := json.Unmarshal(nodeValuesByte, &forceTxData)
			if err != nil {
				return nil, err
			}

			for _, val := range forceTxData {
				data = append(data, val)
			}
			return data, nil
		}
	} else if key == common.MetaDataKey {
		metaDataValues, exit := cache.Datas[key]
		if exit {
			log.Debug("操作:", key, ",value:", string(value), ",当前类型：", reflect.TypeOf(metaDataValues))

			var forceTxData map[string][]common.Datadigest
			metaDataValuesByte, _ := json.Marshal(metaDataValues)
			err := json.Unmarshal(metaDataValuesByte, &forceTxData)
			if err != nil {
				return nil, err
			}

			if v, ok := forceTxData[string(value)]; ok {
				return v, nil
			}
		}
	} else if key == common.FileCrcKey {
		fileCrcValues, exit := cache.Datas[key]
		if exit {
			log.Debug("操作:", key, ",value:", string(value), ",当前类型：", reflect.TypeOf(fileCrcValues))

			var forceTxData map[string][]common.SharedData
			fileCrcValuesByte, _ := json.Marshal(fileCrcValues)
			err := json.Unmarshal(fileCrcValuesByte, &forceTxData)
			if err != nil {
				return nil, err
			}

			if v, ok := forceTxData[string(value)]; ok {
				return v, nil
			}
		}
	} else {
		return cache.Datas[key], nil
	}
	return data, nil
}
