package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
	"com.mgface.disobj/metanode/mq/mgfacemq/memory"
)

// StoreSnapshotData 存储内存数据快照到文件系统中
func StoreSnapshotData(ms *memory.MemoryStore) {
	//清理log文件
	//go clearLogfile(storepath)

	filename := getFilename(ms.StorePath)
	duration := GetDuration()
	flushTicker := time.NewTicker(5 * time.Second)
	data := make([]common.RecMsg, 0)
	var activeBuff int32

	for {
		select {
		case msg := <-ms.Msgs:
			data = append(data, msg)
			if atomic.LoadInt32(&activeBuff) == 1 {
				//同步
				ms.BuffMsgs <- msg
			}
		case <-ms.BuffSemaphore: //如果是同步数据
			log.Info("当前还剩未处理size:", len(data), ",sync时间:", time.Now().Format("2006-01-02 15:04:05"))
			writeDataToFile(data, filename)
			syncFile(filename)

			//todo 快照文件和缓冲可以同时操作，不要互相影响
			ms.FishedSnapshot <- true
			//激活sync消息通道
			atomic.StoreInt32(&activeBuff, 1)

		case <-flushTicker.C: // 每隔5秒钟，将数据写入到文件中
			if len(data) > 0 {
				writeDataToFile(data, filename)
				data = make([]common.RecMsg, 0)
			}
		case <-duration.C: //已经到了凌晨第二天了
			log.Info("当前还剩未处理size:", len(data), ",已经翻滚到第二天日期:", time.Now().Format("2006-01-02 15:04:05"))
			//如果存在数据，那么还是回写到前一天的日志文件中
			if len(data) > 0 {
				writeDataToFile(data, filename)
			}
			// 写新一天的数据
			StoreSnapshotData(ms)
		}
	}
}

// 清理元数据log文件，只保留最新的五个文件
func clearLogfile(storepath string) {
	for {
		time.Sleep(time.Duration(3600) * time.Second)
		files := common.WalkDirectory(storepath)
		//移除老的文件
		if len(files) > 5 {
			tfs := common.FileDescs(files)
			sort.Stable(tfs)
			tfs = tfs[5:]
			for _, v := range tfs {
				os.Remove(v.Fpath)
			}
		}
	}

}

func getFilename(storePath string) string {
	dateFlag := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s%v%v%s", storePath, string(os.PathSeparator), dateFlag, ".log")
	return filename
}

func writeDataToFile(data []common.RecMsg, filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Error(err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)

	//如果存在数据，那么还是回写到前一天的日志文件中
	if len(data) > 0 {
		jsonData, _ := json.Marshal(data)
		byteData := make([]byte, 0)
		byteData = append(byteData, common.IntToBytes(len(jsonData))...)
		byteData = append(byteData, jsonData...)
		_, err = file.Write(byteData)
		if err != nil {
			log.Error(err)
			return
		}
		err = file.Sync()
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func syncFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Error("无法打开文件：", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)

	syncFilePath := fmt.Sprintf("%s.%s", filename, "sync")
	fileSync, err := os.OpenFile(syncFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Error("无法打开文件：", err)
		return
	}
	defer func(fileSync *os.File) {
		err = fileSync.Close()
		if err != nil {
			log.Error(err)
		}
	}(fileSync)

	buf := make([]byte, 4096)
	for {
		n, e := file.Read(buf)
		if e == io.EOF {
			break
		}
		_, err = fileSync.Write(buf[:n])
		if err != nil {
			log.Error(err)
			return
		}
		err = fileSync.Sync()
		if err != nil {
			return
		}
	}
}

func GetDuration() *time.Timer {
	//获取当前时间，放到now里面，要给next用
	now := time.Now()
	// 通过now偏移24小时
	nextDay := now.Add(time.Hour * 24)
	// 获取下一个凌晨的日期
	nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
	// 计算当前时间到凌晨的时间间隔，设置一个定时器
	duration := time.NewTimer(nextDay.Sub(now))
	return duration
}
