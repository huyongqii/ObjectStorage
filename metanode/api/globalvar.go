package api

import "sync"

//全局锁
var metaNodeMutex sync.RWMutex

//因为dynamicMNAddr是全局变量。为了保证可见性，使用时需要使用apiMutex全局锁
var dynamicMNAddr string

//获取dynamicmetanode地址
func GetDynamicMNAddr() string {
	metaNodeMutex.RLock()
	memdata := dynamicMNAddr
	metaNodeMutex.RUnlock()
	return memdata

}

//设置dynamicmetanode地址
func SetDynamicMNAddr(nodeAddr string) {
	metaNodeMutex.Lock()
	dynamicMNAddr = nodeAddr
	metaNodeMutex.Unlock()
}
