package server

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"disobj/common"
)

// StartApiHeartbeat 心跳统一3秒发送一次,也相当于元数据服务注册
func StartApiHeartbeat(nodeAddr string, startflag chan bool) {
	log.Info("获得启动标识1:", <-startflag)

	for {
		client, err := common.NewReCallFuncTCPClient(GetDynamicMetanodeAddr, 3)
		if err != nil {
			log.Warn("API心跳包服务连接元数据节点失败，等待重连......")
			continue
		}

		log.Debug("当前执行的master节点为:", GetDynamicMetanodeAddr())
		req := common.NewRequest(client, "set", "apiNodes", nodeAddr)
		err = req.Run()
		if err != nil {
			log.Warn(fmt.Sprintf("%s,apinode心跳包服务发送心跳失败.", time.Now().Format("2006-01-02 15:04:05")))
		}
		time.Sleep(3 * time.Second)
	}
}
