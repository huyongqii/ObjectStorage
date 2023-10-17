package server

// DataShards 数据分片大小数量
var DataShards = 2

// ParityShards 奇偶校验数量
var ParityShards = 1

func Initval(dataShards, parityShards int) {
	DataShards = dataShards
	ParityShards = parityShards
}
