package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	. "com.mgface.disobj/metanode/mq/mgfacemq/nodeinfo"
	"com.mgface.disobj/metanode/mq/mgfacemq/ops"
	. "com.mgface.disobj/metanode/store"
)

type Server struct {
	Store
	//当前节点的信息
	Nodeinfo  *NodeInfo `json:"nodeinfo"`
	MutexServ sync.RWMutex
}

func (server *Server) Listen(nodeAddr string, ctx context.Context) {
	listen, e := net.Listen("tcp", nodeAddr)

	if e != nil {
		log.Fatal(fmt.Sprintf("tcp连接监听错误:%v", e.Error()))
		//没有办法提供服务,直接退出
		os.Exit(-1)
	}

	for {
		select {
		case <-ctx.Done():
			log.Warn("linsten接受终端发出退出信号")
			return
		default:
			conn, e := listen.Accept()
			if e != nil {
				panic(e)
			}
			go process(conn, server, ctx)
		}
	}
}

func process(conn net.Conn, server *Server, ctx context.Context) {
	var data interface{}
	var err error
	var op byte
	reader := bufio.NewReader(conn) //对客户端连接进行缓冲读取
	for {
		select {
		//从上层传递过来有中断的信息，关闭连接，退出当前协程
		case <-ctx.Done():
			log.Warn("process接受终端发出退出信号")
			conn.Close()
			return
		default:
			//在不发生错误的情况下，客户端可以复用这个连接不断发送命令到服务器，并且得到响应
			op, err = reader.ReadByte()
			if err != nil {
				if err != io.EOF {
					log.Debug("关闭连接发生错误:", err)
				}
				return
			}
			switch op {
			case 'S':
				err = ops.Set(server, reader)
			case 'G':
				data, err = ops.Get(server, reader)
			case 'D':
				err = ops.Del(server, reader)
			case 'F': //文件传输
				data, err = ops.Put(server, reader)
			case 'X': //同步请求
				data, err = ops.Syn(server, reader)
			default:
				err = errors.New(fmt.Sprintf("无效的操作符:%s", string(op)))
				return
			}
			sendRespone(data, err, conn)
		}

	}
}

func sendRespone(value interface{}, err error, conn net.Conn) error {
	if err != nil {
		errStr := err.Error()
		tmp := fmt.Sprintf("-%d ", len(errStr)) + errStr
		_, e := conn.Write(append([]byte(tmp)))
		return e
	}
	data, err := json.Marshal(value)
	vlen := fmt.Sprintf("%d ", len(data))
	_, e := conn.Write(append([]byte(vlen), data...))
	return e

}
