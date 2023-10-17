package cluster

import (
	"strconv"
	"strings"

	"github.com/hashicorp/memberlist"
)

// WrapMemberlistNodes 节点进行排序
type WrapMemberlistNodes []*memberlist.Node

func (s WrapMemberlistNodes) Len() int {
	return len(s)
}

func (s WrapMemberlistNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s WrapMemberlistNodes) Less(i, j int) bool {
	nodenamei := strings.Split(s[i].Name, "-")
	nodei, _ := strconv.ParseInt(nodenamei[0], 10, 64)
	nodenamej := strings.Split(s[j].Name, "-")
	nodej, _ := strconv.ParseInt(nodenamej[0], 10, 64)
	return nodei < nodej
}
