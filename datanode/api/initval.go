package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

//根文件存储路径
var rootLocalStorePath string

func GetRootLocalStorePath() string {
	return rootLocalStorePath
}

func setRootLocalStorePath(rlsp string) {
	rootLocalStorePath = rlsp
}

//每天滚动生成的文件存储路径
var rollingDateLocalStorePath string

//获取滚动日期存储的本地存储路径
func GetRDLocalStorePath() string {
	return rollingDateLocalStorePath
}

func setRDLocalStorePath(rdlsp string) {
	rollingDateLocalStorePath = rdlsp
}

//当前节点的地址
var nodeAddr string

func GetNodeAddr() string {
	return nodeAddr
}

func setNodeAddr(ndAddr string) {
	nodeAddr = ndAddr
}

func Initval(rootLocalStorePath, nodeAddr string) {
	setRootLocalStorePath(rootLocalStorePath)

	datePath := time.Now().Format("20060102")
	nodePath := strings.Replace(nodeAddr, ":", "-", -1)
	//每天生成文件存储的目录
	//二级目录不能随机大小字母，要按照日期来生成
	rdLocalStorePath := fmt.Sprintf("%s%s%s-%s", rootLocalStorePath, string(os.PathSeparator), datePath, nodePath)
	setRDLocalStorePath(rdLocalStorePath)
	log.Debug("文件存储的目录:", rdLocalStorePath)
	//如果目录不存在，这新建该目录，授予目录权限为777

	//目录不存在则创建
	if _, err := os.Stat(rdLocalStorePath); os.IsNotExist(err) {
		os.MkdirAll(rdLocalStorePath, os.ModePerm)
	}
	setNodeAddr(nodeAddr)
}
