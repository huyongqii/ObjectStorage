package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
	"com.mgface.disobj/datanode/api"
	"com.mgface.disobj/datanode/datarepair"
	"com.mgface.disobj/datanode/ops"
)

func StartServer(na, mna, podnamespace string) {
	log.Info("启动数据节点...")
	log.Info(fmt.Sprintf("节点地址:%s", na))
	log.Info(fmt.Sprintf("元数据服务节点地址:%s", mna))

	//后台数据修复
	go datarepair.Repair()

	startFlag := make(chan bool)
	go api.RefreshDNMetaNode(mna, podnamespace, startFlag)
	go api.StartDNHeartbeat(na, startFlag)

	http.HandleFunc("/objects/", ApiHandler)

	common.SupportServeAndGracefulExit(na)
}

func ApiHandler(writer http.ResponseWriter, req *http.Request) {
	method := req.Method
	if method == http.MethodPut {
		ops.Put(writer, req)
		return
	}
	if method == http.MethodGet {
		ops.Get(writer, req)
		return
	}
	writer.WriteHeader(http.StatusMethodNotAllowed)
}
