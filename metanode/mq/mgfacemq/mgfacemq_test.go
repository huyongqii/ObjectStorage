package mgfacemq

import (
	. "com.mgface.disobj/common"
	"flag"
	"fmt"
	"testing"
	"time"
)

func TestTCP(t *testing.T) {
	//元数据服务节点
	DefaultMatanodeAddr := "127.0.0.1:30000"
	server := flag.String("h", DefaultMatanodeAddr, "缓存服务器IP地址")
	op := flag.String("c", "get", "命令行操作,应该为get/set其中一种")
	key := flag.String("k", "test", "key")
	value := flag.String("v", "admin123", "value")
	flag.Parse()
	client := NewReconTCPClient(*server, 3)
	cmd := &Cmd{Name: *op, Key: *key, Value: *value}
	cmd.Run(client)
	if cmd.Error != nil {
		fmt.Println("error:", cmd.Error)
	} else {
		fmt.Println(cmd.Value)
		fmt.Println("等待2S准备退出连接")
		time.Sleep(time.Second * 2)
	}
}
