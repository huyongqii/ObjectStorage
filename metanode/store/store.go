package store

import (
	"bufio"
	. "com.mgface.disobj/metanode/mq/mgfacemq/nodeinfo"
)

type Store interface {
	//设置数据值
	Set(key string, value []byte) error
	//获取数据值
	Get(key string, value []byte) (interface{}, error)
	//删除数据值
	Del(key string, value []byte) error
	//保存文件数据
	Put(fname string, fszie int, reader *bufio.Reader) (interface{}, error)
	//同步请求
	Sync(synrlen int, reader *bufio.Reader) (interface{}, error)
	//显示数据值
	Show() string
	//过期后台数据
	ExpireData(info *NodeInfo, intervalCleanMem int)
	//把数据本身转成Json
	ThisToJson() string
	//对自身赋值操作
	AssignThis(data interface{})
}
