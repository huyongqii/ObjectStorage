package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// TcpClient 请求客户端
type TcpClient struct {
	net.Conn
	reader *bufio.Reader
}

type Request struct {
	client    *TcpClient
	operation string
	key       string
	value     string
}

func NewRequest(client *TcpClient, operation, key, value string) *Request {
	return &Request{client, operation, key, value}
}

func (req *Request) Run() error {
	if req.operation == "get" {
		_, err := req.client.sendGet(req.key, req.value)
		if err != nil {
			return err
		}
		req.value, err = req.client.recvResponse()
		if err != nil {
			return err
		}
		return nil
	}
	if req.operation == "set" {
		_, err := req.client.sendSet(req.key, req.value)
		if err != nil {
			return err
		}
		_, err = req.client.recvResponse()
		if err != nil {
			return err
		}
		return nil
	}
	if req.operation == "del" {
		_, err := req.client.sendDel(req.key, req.value)
		if err != nil {
			return err
		}
		_, err = req.client.recvResponse()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(fmt.Sprintf("未知的命令: %s,仅支持set,get,del操作.", req.operation))
}

func (req *Request) GetValue() string {
	return req.value
}

// 缓存连接，不做无谓的频繁创建连接
var cacheClient sync.Map

func NewTCPClient(server string) *TcpClient {
	c, e := net.Dial("tcp", server)
	if e != nil {
		return nil
	}
	r := bufio.NewReader(c)
	return &TcpClient{c, r}
}

func NewReCallFuncTCPClient(server func() string, maxFailCount int) (*TcpClient, error) {
	failCount := 0
	for {
		serverIp := server()
		log.Debug("serverIp:::::", serverIp)
		c, e := net.DialTimeout("tcp", serverIp, 2*time.Second)

		if e != nil {
			failCount++
			if failCount > maxFailCount {
				return nil, nil
			}
			continue
		}
		if failCount > 0 {
			log.Info(fmt.Sprintf("重连成功，重连的IP地址为:%s", serverIp))
		}
		r := bufio.NewReader(c)
		return &TcpClient{c, r}, nil
	}
}

func NewReconTCPClient(server string, maxFailCount int) *TcpClient {
	failCount := 0

	for {
		c, e := net.DialTimeout("tcp", server, 2*time.Second)

		if e != nil {
			failCount++
			if failCount > maxFailCount {
				return nil
			}
			continue
		}

		r := bufio.NewReader(c)
		return &TcpClient{c, r}
	}
}

func (client *TcpClient) sendDel(key, value string) (n int, err error) {
	klen := len(key)
	vlen := len(value)
	return client.Write([]byte(fmt.Sprintf("D%d %d %s%s", klen, vlen, key, value)))
}

func (client *TcpClient) sendGet(key, value string) (n int, err error) {
	klen := len(key)
	vlen := len(value)
	return client.Write([]byte(fmt.Sprintf("G%d %d %s%s", klen, vlen, key, value)))
}

func (client *TcpClient) sendSet(key, value string) (n int, err error) {
	klen := len(key)
	vlen := len(value)
	return client.Write([]byte(fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)))
}

// 接受客户端响应数据
func (client *TcpClient) recvResponse() (string, error) {
	tmp, e := client.reader.ReadString(' ')
	if e != nil {
		return "", e
	}
	vlen, _ := strconv.Atoi(strings.TrimSpace(tmp))
	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(client.reader, err)
		if e != nil {
			return "", e
		}
		return "", errors.New(string(err))
	}
	value := make([]byte, vlen)
	_, e = io.ReadFull(client.reader, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}
