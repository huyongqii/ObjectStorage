package cluster

import "com.mgface.disobj/metanode/mq/mgfacemq/server"

// StartGossipCluster 启动gossip cluster集群
// 这个函数的作用是启动gossip cluster集群，包括加入集群、显示集群状态、发送集群消息
// 首先加入加入到gossip集群中，gossip集群内可以广播信息，是的每个节点的最终状态是一样的，即拥有相同的信息
// 显示集群状态么，可以看到目前集群中有哪些节点
// 之后同步集群消息，这个消息是master节点发送的，包含了集群中所有节点的信息，这个信息会被所有节点接收到
// 以及包含master的选取，master选取的算法很简单，就是按照节点创建时间排序，取第一个节点作为master
func StartGossipCluster(nodeAddr, cluster, gossipAddr, podNamespace, serviceName string, serv *server.Server, startflag chan bool) {
	// 1.加入集群
	nodeName, broadcasts, list := joinGossipCluster(nodeAddr, cluster, gossipAddr, podNamespace, serviceName, serv)

	// 2.显示集群状态
	showMemberList(list)

	// 3.master发送集群消息
	go sendMsg2Cluster(nodeName, serv, broadcasts, list, startflag)
}
