package common

import (
	"errors"
	"fmt"
)

//执行的命令
type Cmd struct {
	//操作命令
	Name string `json:"name"`
	//操作的key
	Key string `json:"key"`
	//操作的value
	Value string `json:"value"`
	//异常信息
	Error error `json:"error"`
}

//执行客户端请求命令
func (cmd *Cmd) Run(client *TcpClient) {
	if cmd.Name == "get" {
		_, err := client.sendGet(cmd.Key, cmd.Value)
		if err != nil {
			cmd.Error = err
		} else {
			cmd.Value, cmd.Error = client.recvResponse()
		}
		return
	}
	if cmd.Name == "set" {
		_, err := client.sendSet(cmd.Key, cmd.Value)
		if err != nil {
			cmd.Error = err
		} else {
			_, cmd.Error = client.recvResponse()
		}
		return
	}
	if cmd.Name == "del" {
		_, err := client.sendDel(cmd.Key, cmd.Value)
		if err != nil {
			cmd.Error = err
		} else {
			_, cmd.Error = client.recvResponse()
		}
		return
	}
	cmd.Error = errors.New(fmt.Sprintf("未知的命令: %s,仅支持set,get,del操作.", cmd.Name))
}
