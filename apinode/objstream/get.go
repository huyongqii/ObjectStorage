package objstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func directReq(url string) (*GetStream, error) {
	resp, err := http.Get(url)
	//调用失败resp会为nil
	if resp == nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("数据服务返回状态码为: %d", resp.StatusCode)
	}
	return &GetStream{resp.Body}, nil
}

func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("无效的节点 %s ,obj= %s", server, object)
	}
	return directReq(fmt.Sprintf("http://%s/objects/%s", server, object))
}

func (r *GetStream) Read(data []byte) (n int, err error) {
	return r.reader.Read(data)
}
