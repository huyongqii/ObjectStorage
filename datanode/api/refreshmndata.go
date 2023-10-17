package api

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
	"com.mgface.disobj/common/k8s"
)

func RefreshDNMetaNode(metaNode, podNamespace string, startFlag chan bool) {
	if podNamespace != "" {
		val := strings.Split(metaNode, ":")
		metaNodes := k8s.DetectMetaService(val[1], podNamespace, val[0])
		if len(metaNodes) > 0 {
			//随机选一个metaNode
			metaNode = metaNodes[0]
		}
		log.Debug("随机选取的metaNode:", metaNode)
	}
	//该metaNodeAddr作为种子节点取元数据
	client := common.NewReconTCPClient(metaNode, 3)
	if client == nil {
		log.Fatal(fmt.Sprintf("元数据服务连接失败,请提供一个正确的元数据服务IP地址"))
		os.Exit(-1)
	}

	// 第一次获取metaNodes
	metaNodes := getMetaNodes(metaNode)

	for _, metaNode := range metaNodes {
		log.Debug("当前查询的节点::", metaNode.RealNodeValue)
		data := strings.Split(metaNode.RealNodeValue, "-")
		if data[1] == "master" {
			SetDNDynamicMetanodeAddr(data[0])
			syncRun.Do(func() {
				startFlag <- true
				close(startFlag)
			})
		}
	}

	// 之后循环获取metaNodes
	loopRefresh(metaNodes, startFlag)
}

func getMetaNodes(metaNode string) []common.MetaValue {
	client := common.NewReconTCPClient(metaNode, 3)
	if client == nil {
		log.Fatal(fmt.Sprintf("元数据服务连接失败,请提供一个正确的元数据服务IP地址"))
		os.Exit(-1)
	}
	req := common.NewRequest(client, "get", "metaNodes", "")
	err := req.Run()
	if err != nil {
		return nil
	}
	var metaNodes []common.MetaValue
	err = json.Unmarshal([]byte(req.GetValue()), &metaNodes)
	if err != nil {
		return nil
	}
	return metaNodes
}

func loopRefresh(metaNodes []common.MetaValue, startFlag chan bool) {
	printCount := 0
	for {
		sort.Stable(common.WrapMetaValues(metaNodes))
		displayMetaNodes, _ := json.MarshalIndent(metaNodes, "", "\t")
		if printCount > 20 {
			fmt.Println("更新时间:", time.Now().Format("2006-01-02 15:04:05"), ",metaNodes数据::", string(displayMetaNodes))
			printCount = 0
		}

		// 查询出来的数据
		searchData := make([]common.MetaValue, 0)
		for _, metaVale := range metaNodes {
			log.Debug("当前查询的节点::", metaVale.RealNodeValue)
			data := strings.Split(metaVale.RealNodeValue, "-")

			// client = NewReconTCPClient(livenodeAddr, 3)
			client := common.NewTCPClient(data[0])
			if client == nil {
				continue
			}
			req := common.NewRequest(client, "get", "metaNodes", "metaNodes")
			err := req.Run()
			if err != nil {
				continue
			}

			searchData = searchMaster(req.GetValue(), startFlag)
			// 如果返回的数据为空，说明当前节点还没有选择出master(处于不一致状态)，可以暂停100ms，重新让下一个节点取获取
			if GetDNDynamicMetanodeAddr() == "" {
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		if searchData != nil {
			metaNodes = searchData
		} else {
			log.Warn("没有更新到dynamicMetaNode.")
		}

		time.Sleep(500 * time.Millisecond)
		printCount++
	}
}

func searchMaster(value string, startFlag chan bool) []common.MetaValue {
	var metaNodes []common.MetaValue
	err := json.Unmarshal([]byte(value), &metaNodes)
	if err != nil {
		return metaNodes
	}

	var searchData []common.MetaValue
	for _, metaNode := range metaNodes {
		data := strings.Split(metaNode.RealNodeValue, "-")
		if data[1] == "master" {
			SetDNDynamicMetanodeAddr(data[0])
			syncRun.Do(func() {
				startFlag <- true
				close(startFlag)
			})
			// 如果是master的话直接赋值给results，让其遍历最新获取的数据
			searchData = metaNodes
			break
		}
	}
	return searchData
}
