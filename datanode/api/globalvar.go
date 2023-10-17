package api

import "sync"

var syncRun sync.Once

//全局锁
var dataNodeMutex sync.RWMutex

//因为dynamicMetanodeAddr是全局变量。为了保证可见性，使用时需要使用apiMutex全局锁
var dynamicMetanodeAddr string

//获取dynamicmetanode地址
func GetDNDynamicMetanodeAddr() string {
	dataNodeMutex.RLock()
	memdata := dynamicMetanodeAddr
	dataNodeMutex.RUnlock()
	return memdata

}

//设置dynamicmetanode地址
func SetDNDynamicMetanodeAddr(nodeAddr string) {
	dataNodeMutex.Lock()
	dynamicMetanodeAddr = nodeAddr
	dataNodeMutex.Unlock()
}
