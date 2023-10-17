package server

import "sync"

// 用来协调心跳服务和刷新服务
var SyncRun sync.Once

// 全局锁
var apiNodeMutex sync.RWMutex

// 全局数据节点map集合，key为节点的addr，value是节点的更新时间
var dataNodes = make(map[string]string)

func SetDataNodes(key, value string) {
	apiNodeMutex.Lock()
	dataNodes[key] = value
	apiNodeMutex.Unlock()
}

func GetDataNodes(key string) string {
	apiNodeMutex.RLock()
	memdata := dataNodes[key]
	apiNodeMutex.RUnlock()
	return memdata
}

// 因为dynamicMetanodeAddr是全局变量。为了保证可见性，使用时需要使用apiMutex全局锁
var dynamicMetanodeAddr string

// 获取dynamicmetanode地址
func GetDynamicMetanodeAddr() string {
	apiNodeMutex.RLock()
	memdata := dynamicMetanodeAddr
	apiNodeMutex.RUnlock()
	return memdata

}

// 设置dynamicmetanode地址
func SetDynamicMetanodeAddr(nodeAddr string) {
	apiNodeMutex.Lock()
	dynamicMetanodeAddr = nodeAddr
	apiNodeMutex.Unlock()
}
