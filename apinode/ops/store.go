package ops

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"disobj/apinode/objstream"
	"disobj/apinode/server"
	"disobj/common"
	"github.com/klauspost/reedsolomon"
	log "github.com/sirupsen/logrus"
)

func storeObject(reqData io.Reader, objName string) (int, error) {
	readData, _ := ioutil.ReadAll(reqData)
	hashValue := caculateHash(readData)

	client, err := common.NewReCallFuncTCPClient(server.GetDynamicMetanodeAddr, 1)
	if err != nil {
		tips := fmt.Sprintf("storeObject执行获取client连接失败。")
		log.Warn(tips)
		return http.StatusInternalServerError, errors.New(tips)
	}

	request := common.NewRequest(client, "get", "filecrc", hashValue)
	err = request.Run()
	if err != nil {
		return http.StatusMethodNotAllowed, err
	}

	log.Debug("cmd.value=", request.GetValue(), ",type(cmd.value)=", reflect.TypeOf(request.GetValue()))
	// 查询文件是否在元数据存在，元数据不存在，那么写入元数据
	if request.GetValue() == "[]" {
		status, err := buildShared(readData, hashValue, objName)
		if status != http.StatusOK && err != nil {
			return status, err
		}
	}

	log.Debug("更新metadata......")
	resp := make(chan error)
	go buildDesc(int64(len(readData)), objName, hashValue, client, resp)
	err = <-resp
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func caculateHash(readData []byte) string {
	//算出请求数据的整体hash值
	hash := sha256.New()
	hash.Write(readData)
	hashInBytes := hash.Sum(nil)
	hashValue := hex.EncodeToString(hashInBytes)
	return hashValue
}

// 构建文件上传的描述
func buildDesc(datasize int64, objName, hashValue string, client *common.TcpClient, resp chan error) {
	digest := common.Datadigest{
		Index:     time.Now().Unix(),
		Hash:      hashValue,
		Version:   1,
		Datasize:  datasize,
		Created:   time.Now(),
		FileOwner: "admin",
	}

	dkv := &common.DataKeyValue{
		Data: make(map[string][]common.Datadigest),
	}
	dataDigests := make([]common.Datadigest, 0)
	dkv.Data[objName] = append(dataDigests, digest)
	data, _ := json.Marshal(dkv)
	request := common.NewRequest(client, "set", "metadata", string(data))
	err := request.Run()
	resp <- err
}

// 构建数据分片，然后将每个分片上传到数据节点
func buildShared(readData []byte, hashValue, objName string) (int, error) {
	// 分片数据，按照数据分片和校验分片的比例进行分片
	rsencoder, _ := reedsolomon.New(server.DataShards, server.ParityShards, reedsolomon.WithMaxGoroutines(runtime.NumCPU()))
	splitData, _ := rsencoder.Split(readData)

	expectIps := make([]string, 0)

	//todo 可以启用多个goroutine运行，并且需要考虑如果集中提交没有成功， 是否需要回退提交的部分数据
	for index, data := range splitData {
		nodeAddr, err := getRandomDataNode(&expectIps)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		putStream := objstream.NewPutStream(hashValue, nodeAddr, objName, index)

		gzipData, _ := common.GzipEncode(data)
		reader := bytes.NewReader(gzipData)
		// 把数据copy进内存管道
		_, err = io.Copy(putStream, reader)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		// 关闭写通道
		err = putStream.Close()
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

// 获取datanode节点数据
func getDataNodes() []string {
	dn := make([]string, 0)
	var metaValues []common.MetaValue
	err := json.Unmarshal([]byte(server.GetDataNodes("dataNodes")), &metaValues)
	if err != nil {
		return nil
	}
	for _, v := range metaValues {
		dn = append(dn, v.RealNodeValue)
	}
	return dn
}

// getRandomDataNode 选择请求的随机数据节点
func getRandomDataNode(expectIps *[]string) (string, error) {
	// 设置随机种子，确保每次随机产生的数都是需要的数据
	rand.Seed(time.Now().UnixNano())

	dn := getDataNodes()
	n := len(dn)
	if n == 0 {
		return "", errors.New("没有找到数据节点")
	}
	if n < 3 {
		data, _ := json.Marshal(&common.ReturnMsg{Msg: "数据节点必须>=3个", Flag: false})
		return "", errors.New(string(data))
	}

	for {
		if len(*expectIps) == n {
			return "", errors.New("数据节点太少，没有找到可用的数据节点")
		}
		// 随机从当前的数据节点集里面取一个节点
		ranNodeAddr := dn[rand.Intn(n)]
		found := false
		for _, v := range *expectIps {
			// 如果随机出来的数据节点已经存在了，重新随机
			if ranNodeAddr == v {
				found = true
				break
			}
		}
		if !found {
			*expectIps = append(*expectIps, ranNodeAddr)
			return ranNodeAddr, nil
		}
	}
}
