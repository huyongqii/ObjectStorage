package memory

import (
	. "com.mgface.disobj/common"
	. "com.mgface.disobj/metanode/mq/mgfacemq/nodeinfo"
	"time"
)

//过期内存数据
func (cache *MemoryStore) ExpireData(nodeinfo *NodeInfo, intervalCleanMem int) {
	if cache.EnableClean {
		for {
			time.Sleep(time.Duration(intervalCleanMem) * time.Second)
			//判断是否是master
			master := nodeinfo.DecideMaster()
			cache.Mutex.Lock()
			//假如当前节点不是master，把自身变成slave，不主动启动心跳检测过期处理
			//主要是为了让datanaode只向master节点汇报心跳，减少tcp交互，然后由master把内存存储的meta数据同步给其他slave
			//APINode服务可以接入slave，因为slave有完整从master同步过来的内存数据，可以达到提高多个matanoade的使用效率
			if master {
				for key, mvalues := range cache.Datas {
					//对心跳数据进行过期处理
					if key == "dataNodes" || key == "apiNodes" || key == "metaNodes" {
						legal := make([]MetaValue, 0)
						if forceTxData, ok := mvalues.([]MetaValue); ok {
							for _, value := range forceTxData {
								if time.Now().Before(value.Created.Add(time.Duration(intervalCleanMem) * time.Second)) {
									legal = append(legal, value)
								}
							}
						}
						cache.Datas[key] = legal
					}
				}
			}
			cache.Mutex.Unlock()
		}
	}
}
