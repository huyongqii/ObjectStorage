package ops

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
	"com.mgface.disobj/datanode/api"
)

func Put(writer http.ResponseWriter, req *http.Request) {
	filename, realURL := getFilenameAndRealURL(req)

	readData, _ := ioutil.ReadAll(req.Body)
	sharedHashValue := caculateHash(readData)

	file, err := os.Create(filename)
	if err != nil {
		log.Debug(err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Debug(err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
	}(file)

	_, err = file.Write(readData)
	if err != nil {
		log.Debug(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	//强制把缓存页刷到硬盘
	err = file.Sync()
	if err != nil {
		log.Debug(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	index, _ := strconv.Atoi(req.Header.Get("index"))
	hashValue := req.Header.Get("hash")

	//记录文件的CRC
	status, e := buildFileCRC(hashValue, sharedHashValue, realURL, index)
	writer.WriteHeader(status)
	if e != nil {
		log.Debug(e)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(e.Error()))
	}
}

func buildFileCRC(hashValue, sharedHashValue, realURL string, index int) (int, error) {
	// 做文件去重插入数据
	crc := &common.FileCRC{
		Data: make(map[string][]common.SharedData),
	}

	sharedData := common.SharedData{
		SharedFileUrlLocate: realURL,
		SharedFileHash:      sharedHashValue,
		SharedIndex:         index,
	}
	sd := make([]common.SharedData, 0)
	sd = append(sd, sharedData)
	crc.Data[hashValue] = sd
	data, _ := json.Marshal(crc)
	cmd := &common.Cmd{Name: "set", Key: "filecrc", Value: string(data)}

	client, err := common.NewReCallFuncTCPClient(api.GetDNDynamicMetanodeAddr, 3)
	if err != nil {
		tips := fmt.Sprintf("获取动态metanode值失败.")
		log.Warn(tips)
		return http.StatusInternalServerError, errors.New(tips)
	}
	req := common.NewRequest(client, "set", "filecrc", string(data))
	err = req.Run()
	if err != nil {
		return http.StatusInternalServerError, cmd.Error
	}
	
	return http.StatusOK, nil
}

func getFilenameAndRealURL(req *http.Request) (string, string) {
	url := req.URL.EscapedPath()
	objName := strings.Split(url, "/")[2]
	// 修改底层数据名称
	objRealName := objName + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	//例如："10.1.2.207:5000/sdp/20210207-10.1.2.207-5000/abc17787-1612644719934396484"
	filename := api.GetRDLocalStorePath() + string(os.PathSeparator) + objRealName

	var realURL string
	if runtime.GOOS == "windows" {
		realURL = api.GetNodeAddr() + string(os.PathSeparator) + filename
	} else {
		realURL = api.GetNodeAddr() + filename
	}
	return filename, realURL
}

func caculateHash(readData []byte) string {
	//算出请求数据的整体hash值
	hash := sha256.New()
	hash.Write(readData)
	hashInBytes := hash.Sum(nil)
	hashValue := hex.EncodeToString(hashInBytes)
	return hashValue
}
