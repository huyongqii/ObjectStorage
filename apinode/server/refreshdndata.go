package server

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"disobj/common"
)

// RefreshDNData 刷新datanode值，每1秒更新一次
func RefreshDNData(startFlag chan bool) {
	log.Info("获得启动标识2:", <-startFlag)

	client, err := common.NewReCallFuncTCPClient(GetDynamicMetanodeAddr, 3)
	if err != nil {
		log.Warn("连接元数据节点失败，等待重连......")
	}

	var req *common.Request
	printCount := 0
	enabledCached := false

	for {
		if !enabledCached {
			req = common.NewRequest(client, "get", "dataNodes", "dataNodes")
			err = req.Run()
			if err != nil {
				// 如果连接不上元服务器端，那么直接取缓存的数据
				enabledCached = true
			} else {
				// 更新缓存的数据
				log.Debug("刷新dataNodes的数据,dynamicMetaNode节点为:", GetDynamicMetanodeAddr())
				SetDataNodes("dataNodes", req.GetValue())
				enabledCached = false
			}
		}

		if printCount > 5 {
			if enabledCached {
				dns := GetDataNodes("dataNodes")
				data, _ := json.MarshalIndent(dns, "", "\t")
				log.Info("获取【缓存】datanodes节点信息:", string(data))
			} else {
				data, _ := json.MarshalIndent(req.GetValue(), "", "\t")
				log.Info("获取【实时】datanodes服务节点信息:", string(data))
			}
			printCount = 0
		}

		if enabledCached {
			log.Info(fmt.Sprintf("[更新内存datanode]重连元数据服务端[%s]...", GetDynamicMetanodeAddr()))
			client, err = common.NewReCallFuncTCPClient(GetDynamicMetanodeAddr, 1)
			if err != nil {
				log.Warn("连接元数据节点失败，等待重连......")
			} else {
				enabledCached = false
			}
		}

		time.Sleep(1 * time.Second)
		printCount++
	}
}
