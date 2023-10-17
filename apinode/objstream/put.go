package objstream

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type PutStream struct {
	//获取写数据的管道
	writer *io.PipeWriter `json:"writer"`
	//传输过程收集的异常信息
	err chan error `json:"err"`
}

func (stream *PutStream) Write(data []byte) (n int, err error) {
	return stream.writer.Write(data)
}

func (stream *PutStream) Close() error {
	//关闭pipe写管道，让管道的reader读到EOF
	stream.writer.Close()
	return <-stream.err
}

// NewPutStream 创建HTTP客户端，发送数据给数据中心
func NewPutStream(hashValue, nodeAddr, objName string, index int) *PutStream {
	log.Debug(fmt.Sprintf("入参:%s,%s,%s,%d", hashValue, nodeAddr, objName, index))
	// 获得一个同步内存管道，同时获得一个输入和一个输出
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		url := fmt.Sprintf("http://%s/objects/%s", nodeAddr, objName)
		// 没有数据时会阻塞在这步
		request, err := http.NewRequest("PUT", url, reader)
		if err != nil {
			log.Warn(fmt.Sprintf("NewRequest失败:%s", err.Error()))
			c <- err
			return
		}
		request.Header.Set("index", strconv.Itoa(index))
		request.Header.Set("hash", hashValue)
		client := http.Client{}
		resp, err := client.Do(request)
		if err == nil && resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("datanode返回的HTTP状态码:%d", resp.StatusCode)
		}
		c <- err
	}()
	return &PutStream{writer, c}
}
