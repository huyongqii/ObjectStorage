package api

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
	"com.mgface.disobj/metanode/mq/mgfacemq/server"
)

// StartMDHeartbeat metadata节点也需要注册，供ApiNode和datanode使用
// 向master节点发送心跳数据，数据包含当前节点的地址和节点状态
func StartMDHeartbeat(nodeAddr string, serv *server.Server, startflag chan bool) {
	log.Info("获得启动标识:", <-startflag)
restart:
	for {
		client, _ := common.NewReCallFuncTCPClient(GetDynamicMNAddr, 3)
		if client == nil {
			log.Warn("metaNode心跳包服务连接master节点失败，等待重连......")
			goto restart
		}
		log.Debug("当前执行的master节点为:", GetDynamicMNAddr())
		nodeStatus := fmt.Sprintf("%s-%s", nodeAddr, serv.Nodeinfo.GetNodeInfo())
		cmd := &common.Cmd{Name: "set", Key: "metaNodes", Value: nodeStatus}
		cmd.Run(client)
		if cmd.Error != nil {
			log.Warn(fmt.Sprintf("%s,metanode心跳包服务发送心跳失败.", time.Now().Format("2006-01-02 15:04:05")))
		}
		time.Sleep(3 * time.Second)
	}
}
