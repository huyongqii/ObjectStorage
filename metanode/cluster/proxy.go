package cluster

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	"github.com/hashicorp/memberlist"

	"com.mgface.disobj/common"
	"com.mgface.disobj/metanode/mq/mgfacemq/server"
)

// 委托对象，使用该对象处理接收到的消息，委托对象可以处理消息的接收、处理和传递等操作。
type delegate struct {
	mtx        sync.RWMutex
	items      map[string]interface{}
	broadcasts *memberlist.TransmitLimitedQueue
	serv       *server.Server
}

type gossipInfo struct {
	Action string // add, del
	Data   interface{}
}

func buildGossipInfo(Action string, Data interface{}) gossipInfo {
	return gossipInfo{
		Action: Action,
		Data:   Data,
	}
}

// NodeMeta 返回本地节点的元数据。你可以在这个函数中定义本地节点的附加信息，这些信息会在远程节点的NodeEventDelegate.NotifyJoin()方法中被访问到。
func (proxy *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

// NotifyMsg 获取gossip服务传递过来的数据
func (proxy *delegate) NotifyMsg(data []byte) {
	dst := make([]byte, len(data))
	copy(dst, data)
	if len(dst) == 0 {
		return
	}

	var info gossipInfo
	err := json.Unmarshal(dst, &info)
	if err != nil {
		return
	}

	decoded, _ := base64.StdEncoding.DecodeString(info.Data.(string))
	dx, _ := common.GzipDecode(decoded)
	if info.Action == "heartbeat" { //心跳数据
		proxy.serv.AssignThis(dx)
	}
}

// GetBroadcasts 获取要广播的消息。你可以在这个函数中指定要广播的消息内容。
func (proxy *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return proxy.broadcasts.GetBroadcasts(overhead, limit)
}

// LocalState 返回本地节点的状态。你可以在这个函数中定义本地节点的状态信息。
func (proxy *delegate) LocalState(join bool) []byte {
	//proxy.mtx.RLock()
	//m := proxy.items
	//proxy.mtx.RUnlock()
	//data, _ := json.Marshal(m)
	//return data
	return nil
}

// MergeRemoteState 合并远程节点的状态。你可以在这个函数中定义如何合并远程节点的状态信息。
func (proxy *delegate) MergeRemoteState(buf []byte, join bool) {
	//if len(buf) == 0 {
	//	return
	//}
	//if !join {
	//	return
	//}
	//var m map[string]string
	//if err := json.Unmarshal(buf, &m); err != nil {
	//	return
	//}
	//proxy.mtx.Lock()
	//for k, v := range m {
	//	proxy.items[k] = v
	//}
	//proxy.mtx.Unlock()
}
