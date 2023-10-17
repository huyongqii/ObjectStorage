package common

const (
	ApiNodeKey  = "apiNodes"
	DataNodeKey = "dataNodes"
	MetaNodeKey = "metaNodes"
	FileCrcKey  = "filecrc"
	MetaDataKey = "metadata"
)

// ReturnMsg 返回客户端信息
type ReturnMsg struct {
	Msg  string `json:"msg"`  //返回信息
	Flag bool   `json:"flag"` //返回的错误标识
}

// RecMsg 服务器段接收到的消息
type RecMsg struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}
