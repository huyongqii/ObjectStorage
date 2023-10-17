package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"com.mgface.disobj/common"
	"com.mgface.disobj/metanode/mq/mgfacemq/memory"

	log "github.com/sirupsen/logrus"
)

// LoadSnapshotData 加载内存数据快照到文件系统中
func LoadSnapshotData(cache *memory.MemoryStore) {
	storePath := cache.StorePath
	// 目录不存在直接跳过加载
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		return
	}
	filePaths := readMetafile(storePath)
	for _, filePath := range filePaths {
		// 跳过同步给其他节点的文件
		if strings.Contains(filePath, ".sync") {
			continue
		}
		log.Info(fmt.Sprintf("加载快照文件:%s", filePath))
		f, _ := os.OpenFile(filePath, os.O_RDONLY, 0644)
		for {
			msgSize := make([]byte, 4)
			_, e := f.Read(msgSize)
			if e == io.EOF {
				break
			}
			size := common.BytesToInt(msgSize)
			ev := make([]byte, size)
			f.Read(ev)
			var msgs []common.RecMsg
			err := json.Unmarshal(ev, &msgs)
			if err != nil {
				return
			}
			for _, v := range msgs {
				err = cache.Set(v.Key, []byte(v.Val.(string)))
				if err != nil {
					return
				}
			}
		}
		f.Close()
	}
	cache.LoadingSnapshot = false //加载快照文件完成
}

func readMetafile(storePath string) []string {
	filePaths := make([]string, 0)

	fds := common.WalkDirectory(storePath)

	for _, v := range fds {
		filePaths = append(filePaths, v.Fpath)
	}
	return filePaths
}
