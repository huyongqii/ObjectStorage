package common

import (
	"encoding/json"
	"strings"
	"time"
)

type MetaValue struct {
	RealNodeValue string    `json:"real_node_value"` //node_ip值
	Created       time.Time `json:"created"`         //缓存创建时间
}

type WrapMetaValues []MetaValue

func (s WrapMetaValues) Len() int {
	return len(s)
}

func (s WrapMetaValues) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s WrapMetaValues) Less(i, j int) bool {
	nodeflagi := strings.Split(s[i].RealNodeValue, "-")
	var vari int
	if nodeflagi[1] == "master" {
		vari = 99
	} else {
		vari = 100
	}
	nodeflagj := strings.Split(s[j].RealNodeValue, "-")
	var varj int
	if nodeflagj[1] == "master" {
		varj = 99
	} else {
		varj = 100
	}
	return vari < varj
}

// 定义一个新的结构体用来序列化标识
type tempMetaValue struct {
	RealNodeValue string `json:"real_node_value"`
	Created       string `json:"created"`
}

//实现它的json序列化方法
func (this MetaValue) MarshalJSON() ([]byte, error) {
	// 定义一个新的结构体
	tempMetaValue := &tempMetaValue{
		RealNodeValue: this.RealNodeValue,
		Created:       this.Created.Format("2006-01-02 15:04:05"),
	}
	return json.Marshal(tempMetaValue)
}

//实现它的json反序列化方法
func (this *MetaValue) UnmarshalJSON(data []byte) error {
	var tmp tempMetaValue
	json.Unmarshal(data, &tmp)

	this.RealNodeValue = tmp.RealNodeValue
	local, err := time.ParseInLocation("2006-01-02 15:04:05", tmp.Created, time.Local)
	this.Created = local
	this = &MetaValue{tmp.RealNodeValue, local}
	return err
}
