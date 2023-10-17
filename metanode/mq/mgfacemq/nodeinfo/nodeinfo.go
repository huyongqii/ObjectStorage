package nodeinfo

import (
	"sync"
)

//节点的状态情况
type NodeInfo struct {
	//节点的标识:master,slave
	NodeFlag      string `json:"node_flag"`
	MutexNodeInfo sync.RWMutex
}

//判断是否是master
func (nodeinfo *NodeInfo) DecideMaster() bool {
	flag := ""
	nodeinfo.MutexNodeInfo.RLock()
	flag = nodeinfo.NodeFlag
	nodeinfo.MutexNodeInfo.RUnlock()
	return flag == "master"
}

//设置节点为master
func (nodeinfo *NodeInfo) SetMaster() {
	nodeinfo.MutexNodeInfo.Lock()
	nodeinfo.NodeFlag = "master"
	nodeinfo.MutexNodeInfo.Unlock()
}

//获得节点的状态
func (nodeinfo *NodeInfo) GetNodeInfo() string {
	flag := ""
	nodeinfo.MutexNodeInfo.RLock()
	flag = nodeinfo.NodeFlag
	nodeinfo.MutexNodeInfo.RUnlock()
	return flag
}
