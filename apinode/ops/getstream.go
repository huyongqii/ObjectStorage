package ops

import (
	"fmt"
	"io"

	"disobj/apinode/objstream"
	"disobj/apinode/server"
)

// 获得数据流
func getStream(object string) ([]io.Reader, []error, int64) {
	// todo 如果定位不到，说明可能metadatanode挂了，那么直接根据缓存的数据连接datanode
	nodes, objRealNames, dataSize := server.Locate(object, 3)
	errors := make([]error, server.DataShards)
	// 返回的节点数组应该是是数据分片+奇偶校验分片数量
	totalSizes := server.DataShards + server.ParityShards
	if len(nodes) == 0 || len(nodes) != totalSizes {
		errors := append(errors, fmt.Errorf("obj:%s 定位失败", object))
		return nil, errors, 0
	}
	// 我们只需数据分片大小的数据,不需要奇偶校验分片数据
	reader := make([]io.Reader, server.DataShards)
	//todo 这里需要考虑到查询数据的时候，数据丢失，需要使用SR纠删码进行数据修复
	//todo 并且考虑到get操作可能没有那么及时，所以我们需要在后台进行定期数据轮训做数据修复
	for index := range nodes[:server.DataShards] {
		reader[index], errors[index] = objstream.NewGetStream(nodes[index], objRealNames[index])
	}
	return reader, errors, dataSize
}
