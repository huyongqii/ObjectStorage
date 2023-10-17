package datarepair

import (
	. "com.mgface.disobj/datanode/api"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// 存储文件的描述信息
type fileDesc struct {
	filename string `json:"filename"`
	filehash string `json:"filehash"`
}

// Repair 主要通过hash值判断数据是否发生变化，如果发生变化，那么需要重新生成数据块
// 首先便利本地文件，获得本地文件的hash值，然后获取，metaNode中存放的文件的hash值，进行比较
// 如果发生变化，那么需要重新生成数据块
func Repair() {
	for {
		filedescs := make([]fileDesc, 0)

		//也可以直接使用
		//files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects")
		//for i := range files {
		//	hash := strings.Split(filepath.Base(files[i]), ".")[0]
		//	verify(hash)
		//}
		filepath.Walk(GetRDLocalStorePath(), func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			readdata, _ := ioutil.ReadFile(path)
			//算出文件的hash值
			hash := sha256.New()
			hash.Write(readdata)
			hashInBytes := hash.Sum(nil)
			sharedHashValue := hex.EncodeToString(hashInBytes)
			desc := fileDesc{
				filename: GetNodeAddr() + string(os.PathSeparator) + f.Name(),
				filehash: sharedHashValue,
			}
			filedescs = append(filedescs, desc)
			//log.Debug("file=>", path, "f=>", f.Name())
			return nil
		})

		//todo 准备查询filecrc数据，判断我们读取的文件是否hash码值变更了，如果变了或者不存在，说明数据已经丢失了，需要重塑出来
		//todo 先把批量数据发送到mgfaceMQ里面的filecrc进行校验，如果校验失败，那么在mgfaceMQ里面数据还原，并对需要还原的数据块
		//todo 重新插入
		time.Sleep(5 * time.Second)
	}

}
