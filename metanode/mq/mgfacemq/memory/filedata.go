package memory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/common"
)

func (cache *MemoryStore) Put(fname string, fsize int, reader *bufio.Reader) (recode interface{}, err error) {

	rfname := fmt.Sprintf("%s%s%s", cache.StorePath, string(os.PathSeparator), fname)

	if e := common.Savefile(rfname, fsize, reader); e != nil {
		return nil, e
	}

	//如果是包含sync文件，把内容同步到相同文件名的文件中，然后删除.sync文件
	if strings.Contains(rfname, ".sync") {

		fdsrc, _ := os.OpenFile(rfname, os.O_RDONLY, 0644)
		for {
			msgsize := make([]byte, 4)
			_, e := fdsrc.Read(msgsize)
			if e == io.EOF {
				break
			}
			size := common.BytesToInt(msgsize)
			ev := make([]byte, size)
			fdsrc.Read(ev)
			var msgs []common.RecMsg
			json.Unmarshal(ev, &msgs)
			for _, v := range msgs {
				fileSet(cache, v.Key, []byte(v.Val.(string)))
			}
		}

		//重头开始读取
		fdsrc.Seek(0, io.SeekStart)

		rs := []rune(fname)
		desname := string(rs[0:strings.Index(fname, ".sync")])
		desname = fmt.Sprintf("%s%s%s", cache.StorePath, string(os.PathSeparator), desname)
		fddes, _ := os.OpenFile(desname, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		io.Copy(fddes, fdsrc)

		fddes.Sync()
		fddes.Close()
		fdsrc.Close()
		os.Remove(rfname)
	} else {
		f, _ := os.OpenFile(rfname, os.O_RDONLY, 0644)
		for {
			msgsize := make([]byte, 4)
			_, e := f.Read(msgsize)
			if e == io.EOF {
				break
			}
			size := common.BytesToInt(msgsize)
			ev := make([]byte, size)
			f.Read(ev)
			var msgs []common.RecMsg
			json.Unmarshal(ev, &msgs)
			for _, v := range msgs {
				fileSet(cache, v.Key, []byte(v.Val.(string)))
			}
		}
		f.Close()
	}

	status := make(map[string]string)
	status[fname] = "success"
	return status, nil
}

func fileSet(cache *MemoryStore, key string, value []byte) {
	log.Debug("key::", key, ",value:", string(value))
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()
	if key == "filecrc" {
		sysvalues, exit := cache.Datas[key]
		var filecrc common.FileCRC
		json.Unmarshal(value, &filecrc)
		if exit {
			//加载内存数据
			var loadsysvalues map[string][]common.SharedData
			datas, _ := json.Marshal(sysvalues)
			json.Unmarshal(datas, &loadsysvalues)
			for k, v := range filecrc.Data {
				if vdata, ok := loadsysvalues[k]; ok { //假如存在，那么数据添加
					flag := false
					for _, dd := range vdata {
						if dd.SharedFileHash == v[0].SharedFileHash {
							flag = true
						}
					}
					if !flag {
						loadsysvalues[k] = append(vdata, v[0])
					}

				} else {
					loadsysvalues[k] = v
				}
			}
			cache.Datas[key] = loadsysvalues
		} else {
			cache.Datas[key] = filecrc.Data
		}
	}

	if key == "metadata" {
		sysvalues, exit := cache.Datas[key]
		//假如metadata不存在，那么进行初始化操作
		var valueDKV common.DataKeyValue
		json.Unmarshal(value, &valueDKV) //解析value的数据到dkv变量
		if exit {
			var loadsysvalues map[string][]common.Datadigest
			datas, _ := json.Marshal(sysvalues)
			json.Unmarshal(datas, &loadsysvalues)

			for k, v := range valueDKV.Data {
				if vdata, ok := loadsysvalues[k]; ok {
					flag := false
					for _, dx := range vdata {
						if dx.Version == v[0].Version {
							flag = true
						}
					}

					if !flag {
						loadsysvalues[k] = append(vdata, v[0])
					}
				} else {
					loadsysvalues[k] = v
				}
			}
			cache.Datas[key] = loadsysvalues
		} else {
			//假如数据不存在
			cache.Datas[key] = valueDKV.Data
		}

	}
}
