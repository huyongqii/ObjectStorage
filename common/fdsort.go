package common

import "os"

type FileDesc struct {
	Filename string `json:"filename"` //文件名称
	Fpath    string `json:"fpath"`    //文件路径
	FD       os.FileInfo `json:"-"` //文件内容
}

type FileDescs []FileDesc

//数据分片切片长度
func (s FileDescs) Len() int {
	return len(s)
}

//数据分片的数据交换
func (s FileDescs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//数据分片的比较
func (s FileDescs) Less(i, j int) bool {
	//降序
	return s[i].Filename > s[j].Filename
}
