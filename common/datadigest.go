package common

import "time"

// Datadigest 文件信息摘要
type Datadigest struct {
	//文件索引号
	Index int64 `json:"index"`
	//数据摘要
	Hash string `json:"hash"`
	//数据版本
	Version int64 `json:"version"`
	//数据大小
	Datasize int64 `json:"datasize"`
	//数据创建时间
	Created time.Time `json:"created"`
	//文件名称
	//Filename string `json:"filename"`
	//文件归属用户
	FileOwner string `json:"file_owner"`
}

// DataKeyValue 上传对象的元数据关联类，Key为对象名称，value是对象关联信息数组
// 因为对象上传存在多个版本，所以使用数据进行级联
type DataKeyValue struct {
	Data map[string][]Datadigest `json:"file_data"`
}

// SharedData 数据进行分片存储结构体类
type SharedData struct {
	//分片文件网络存储路径
	SharedFileUrlLocate string `json:"file_url_locate"`
	//分片文件的hash值
	SharedFileHash string `json:"file_hash"`
	//分片的索引号
	SharedIndex int `json:"shared_index"`
}

// FileCRC 文件去重功能结构体
type FileCRC struct {
	// key为文件的hash值，value对应的是分片的数据
	Data map[string][]SharedData `json:"data"`
}
