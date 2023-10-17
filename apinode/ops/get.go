package ops

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"disobj/common"
)

func Get(respo http.ResponseWriter, r *http.Request) {
	// nodeAddr/objects/objName
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	streams, errors, dataSize := getStream(object)

	for _, v := range errors {
		if v != nil {
			log.Warn("###发生错误:", v)
			respo.WriteHeader(http.StatusInternalServerError)
			respo.Write([]byte(v.Error()))
			return
		}
	}
	resp := make([]byte, 0)
	for _, data := range streams {
		data, e := ioutil.ReadAll(data)
		if e != nil {
			log.Warn("读取数据发生失败:", e.Error())
			respo.WriteHeader(http.StatusInternalServerError)
			respo.Write([]byte(e.Error()))
		}
		//进行gzip解压缩
		decodeData, _ := common.GzipDecode(data)
		resp = append(resp, decodeData...)
	}

	reader := bytes.NewReader(resp[:dataSize])
	//todo 查看其它需要进行GZIP压缩的操作
	log.Debug("数据压缩识别.")
	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	if acceptGzip {
		respo.Header().Set("content-encoding", "gzip")
		gzipWriter := gzip.NewWriter(respo)
		io.Copy(gzipWriter, reader)
		gzipWriter.Close()
	} else {
		io.Copy(respo, reader)
	}
}
