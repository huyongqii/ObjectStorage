package cluster

import (
	. "github.com/hashicorp/memberlist"
	log "github.com/sirupsen/logrus"
)

type mgfaceEventDelegate struct {
}

func (c *mgfaceEventDelegate) NotifyJoin(n *Node) {
	log.Info("通知节点加入......", n.Name)
}

func (c *mgfaceEventDelegate) NotifyLeave(n *Node) {
	log.Info("通知节点离开......", n.Name)
}

func (c *mgfaceEventDelegate) NotifyUpdate(n *Node) {
	log.Info("通知更新......", n.Name)
}
