package memory

import (
	"bufio"
	. "com.mgface.disobj/common"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
	"sync"
)

var nodes sync.Map

var sendBuffSyn sync.Once

//发送缓冲数据到客户端
func sendBuffMsg(ms *MemoryStore) {
	//添加客户端服务IP到缓冲发送里面
	for {
		if msg, ok := <-ms.BuffMsgs; ok {
			nodes.Range(func(key, value interface{}) bool {
				fmt.Println("k:", key, "v1:", value)
				servIp := strings.Split(key.(string), "-")[1]
				cli := NewReconTCPClient(servIp, 3)
				if cli == nil {
					nodes.Delete(key)
					return true
				}
				cmd := &Cmd{Name: "set", Key: msg.Key, Value: msg.Val.(string)}
				cmd.Run(cli)
				return true
			})
		}
	}
}

func reqSyncSnapshot(storepath string, client *TcpClient) {
	files := WalkDirectory(storepath)
	fds := FileDescs(files)
	sort.Stable(fds)
	//传输文件从前一天的开始，因为当天日志包含快照文件，所以下标从2开始
	if len(fds) < 2 {
		log.Warn("快照文件有问题,请检查.")
	}
	dd, _ := json.Marshal(fds)

	fmt.Println("fds:", string(dd))

	//发送sync文件
	syncfd := fds[0]
	if syncfd.FD.Size() > 0 {
		replyMsg := Sendfile(syncfd, client)
		log.Info("接受返回值:", replyMsg)
		//接受成功，删除该快照文件
		if strings.Contains(replyMsg, "success") {
			os.Remove(syncfd.Fpath)
		}
	} else {
		os.Remove(syncfd.Fpath)
	}
	//发送其他日志文件
	for _, fd := range fds[2:] {
		if fd.FD.Size() == 0 {
			continue
		}
		replyMsg := Sendfile(fd, client)
		log.Info("接受返回值:", replyMsg)
	}
}

func (cache *MemoryStore) Sync(synrlen int, reader *bufio.Reader) (recode interface{}, err error) {

	data := make([]byte, synrlen)
	reader.Read(data)

	decoded, _ := base64.StdEncoding.DecodeString(string(data))
	dx, _ := GzipDecode(decoded)

	var sync map[string]string

	json.Unmarshal(dx, &sync)

	//启动缓冲的数据信号
	cache.BuffSemaphore <- true

	nodeinfo := sync["nodeinfo"]

	servIp := strings.Split(nodeinfo, "-")[1]
	//添加客户端服务IP到缓冲发送里面
	client := NewReconTCPClient(servIp, 3)

	if sync["length"] == "all" {
		//快照文件生成完毕
		if _, ok := <-cache.FishedSnapshot; ok {
			//发送快照文件
			reqSyncSnapshot(cache.StorePath, client)
		}
	}

	nodes.Store(nodeinfo, nodeinfo)

	sendBuffSyn.Do(func() {
		go sendBuffMsg(cache)
	})
	fmt.Println("###sync end###")

	status := make(map[string]string)
	status["sync"] = "success"
	return status, nil
}
