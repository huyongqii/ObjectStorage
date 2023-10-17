package api

import (
	"time"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
)

// StartDNHeartbeat 心跳统一3秒发送一次
func StartDNHeartbeat(nodeAddr string, startflag chan bool) {
	log.Info("获得启动标识:", <-startflag)

	client, err := common.NewReCallFuncTCPClient(GetDNDynamicMetanodeAddr, 3)
	if err != nil {
		log.Warn("datanode心跳包服务连接元数据节点失败，等待重连......")
	}

	for {
		if client == nil {
			client, err = common.NewReCallFuncTCPClient(GetDNDynamicMetanodeAddr, 3)
			if err != nil {
				log.Warn("datanode心跳包服务连接元数据节点失败，等待重连......")
				continue
			}
		}

		// 发送心跳包操作
		log.Debugf("当前执行的master节点为: %s", GetDNDynamicMetanodeAddr())
		req := common.NewRequest(client, "set", "dataNodes", nodeAddr)
		err = req.Run()
		if err != nil {
			log.Warnf("%s, datanode心跳包服务发送心跳失败.", time.Now().Format("2006-01-02 15:04:05"))
		}

		timer := time.NewTimer(3 * time.Second)
		<-timer.C
	}
}
