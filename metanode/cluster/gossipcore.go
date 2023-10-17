package cluster

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
	"com.mgface.disobj/common/k8s"
	"com.mgface.disobj/metanode/api"
	"com.mgface.disobj/metanode/mq/mgfacemq/server"
)

// 这里有一点不懂的就是，为什么在发送消息给其他集群时，当集群master变更，要立即同步一次信息呢？这个信息是nodename

func joinGossipCluster(nodeAddr, cluster, gossipAddr, podNamespace, serviceName string, serv *server.Server) (string, *memberlist.TransmitLimitedQueue, *memberlist.Memberlist) {
	splitInfo := strings.Split(gossipAddr, ":")
	var gossipPort string
	if len(splitInfo) == 1 {
		gossipPort = gossipAddr
		gossipAddr = fmt.Sprintf("%s:%s", strings.Split(nodeAddr, ":")[0], gossipAddr)
	}

	nodeName := getNodeName(nodeAddr)
	// 用于指定集群中的代理对象。在Gossip协议中，代理对象负责处理集群中的事件和状态变化
	proxy := &delegate{
		mtx:   sync.RWMutex{},
		items: make(map[string]interface{}),
		serv:  serv,
	}
	conf := buildConf(nodeName, gossipAddr, proxy)
	list, err := memberlist.Create(conf)
	if err != nil {
		panic("错误创建集群信息: " + err.Error())
	}
	broadcasts := &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return list.NumMembers()
		},
		RetransmitMult: 1, //最大传输次数
	}
	proxy.broadcasts = broadcasts

	// 加入存在的集群节点，最少需要指定一个已知的节点信息
	allNodes := strings.Split(cluster, ",")
	// 假如长度为1，说明是从k8s yaml env环境过来的值
	if len(allNodes) == 1 {
		allNodes = k8s.DetectMetaService(gossipPort, podNamespace, serviceName)
	}

	_, err = list.Join(allNodes)
	if err != nil {
		panic("错误加入集群: " + err.Error())
		return "", nil, nil
	}

	return nodeName, broadcasts, list
}

// 以节点启动的纳秒数作为节点名称,节点名称为：节点创建的纳秒数-节点的addr值
func getNodeName(nodeAddr string) string {
	currentNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	nodeName := fmt.Sprintf("%s-%s", currentNano, nodeAddr)
	log.Info("当前节点名称:", nodeName)

	return nodeName
}

func buildConf(nodeName, gossipAddr string, proxy *delegate) *memberlist.Config {
	// 获取默认的局域网配置(DefaultLANConfig)对象
	conf := memberlist.DefaultLANConfig()

	conf.Delegate = proxy
	// conf.Events字段也被赋值为一个事件代理对象，用于处理集群中的事件，比如节点加入、离开等
	conf.Events = &mgfaceEventDelegate{}
	conf.UDPBufferSize = 50_000 //gossip包传输的最大长度

	conf.Name = nodeName

	splitInfo := strings.Split(gossipAddr, ":")
	if len(splitInfo) != 2 {
		log.Fatal("元数据IP地址格式为IP:Port")
		return nil
	}
	conf.BindAddr = splitInfo[0]
	conf.BindPort, _ = strconv.Atoi(splitInfo[1])
	conf.LogOutput = ioutil.Discard
	return conf
}

// 显示当前集群信息
func showMemberList(list *memberlist.Memberlist) {
	go func() {
		for {
			log.Debug("@@@@@@@@@当前集群@@@@@@@@@")
			for i, v := range list.Members() {
				log.Debug(fmt.Sprintf("节点(%d)-%s", i, v.Name))
			}
			log.Debug("#######10S~15S刷新#######")
			time.Sleep(time.Duration(10+rand.Intn(5)) * time.Second)
		}
	}()
}

// 用来启动心跳服务
var startHBRun sync.Once

// 给gossip集群发送消息
func sendMsg2Cluster(nodeName string, serv *server.Server, broadcasts *memberlist.TransmitLimitedQueue,
	list *memberlist.Memberlist, startFlag chan bool) {
	//todo 1.如果存在网络分区，会存在多个master,这个需要重新调整代码,要考虑合并多个master数据存储元数据（心跳数据不需要）
	//todo 2.如果重新连接上之后发现自己状态从master变成slave之后，那么需要把数据同步给当前的master

	printCount := 0
	for {
		wrapNode := WrapMemberlistNodes(list.Members())
		sort.Stable(wrapNode)

		// 不是master，那么就判断当前节点是否创建的最早，那么就把当前节点设置为master
		if !serv.Nodeinfo.DecideMaster() && wrapNode[0].Name == nodeName {
			serv.Nodeinfo.SetMaster()
			log.Info(fmt.Sprintf("把当前节点[%s]设置为master.", nodeName))
		}

		// 当前节点是master，那么就广播数据
		if serv.Nodeinfo.DecideMaster() {
			broadcasts.QueueBroadcast(&broadcast{
				msg:    serialization(serv),
				notify: nil,
			})
		}

		// Name数据格式为:节点创建的纳秒数-节点的addr值
		masterNode := strings.Split(wrapNode[0].Name, "-")[1]
		if printCount > 10 {
			log.Info("当前节点状态:", serv.Nodeinfo.GetNodeInfo())
			log.Debug("masterMetadata::", masterNode)
			printCount = 0
		}

		// 假如master有变化，那么直接重新同步一次，为什么要重新同步，因为master变化之后，可能会导致数据不一致
		if !serv.Nodeinfo.DecideMaster() && masterNode != api.GetDynamicMNAddr() {
			log.Info(fmt.Sprintf("当前master:%s,上一个master:%s", masterNode, api.GetDynamicMNAddr()))
			// 发送sync同步请求master的快照文件，告知接收快照的服务端口是哪一个
			info := make(map[string]string)
			info["length"] = "all"
			info["nodeinfo"] = nodeName
			dx, _ := json.Marshal(info)
			zipByte, _ := common.GzipEncode(dx)
			byteData := base64.StdEncoding.EncodeToString(zipByte)
			// 添加客户端服务IP到缓冲发送里面
			client := common.NewReconTCPClient(masterNode, 3)
			data := []byte(fmt.Sprintf("X%d %s", len(byteData), byteData))
			client.Conn.Write(data)

			recode := make([]byte, 4096)
			n, _ := client.Conn.Read(recode)
			replyMsg := string(recode[:n])
			log.Debug("replyMsg::", replyMsg)
		}

		api.SetDynamicMNAddr(masterNode)
		//启动心跳服务
		startHBRun.Do(func() {
			startFlag <- true
			close(startFlag)
		})

		time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)
		printCount++
	}
}

func serialization(serv *server.Server) []byte {
	// 只写心跳数据给其他metaNode节点
	strBytes := []byte(serv.ThisToJson())
	zipByte, _ := common.GzipEncode(strBytes)
	encoded := base64.StdEncoding.EncodeToString(zipByte)
	gossipInfo := buildGossipInfo("heartbeat", encoded)
	data, _ := json.Marshal(gossipInfo)
	return data
}
