package ops

//
//import (
//	"com.mgface.disobj/apinode/server"
//	"com.mgface.disobj/apinode/objstream"
//	"math/rand"
//	"time"
//)
//
//// objName 对象名称
//// index 数据分片
//func putStream(hashValue, objName string, index int, expectIps []string) (*objstream.PutStream, []string, error) {
//	nodeAddr, expectIps, err := server.ChooseRandomDataNode(index, expectIps)
//	if err != nil {
//		return nil, nil, err
//	}
//	putstream := objstream.NewPutStream(hashValue, nodeAddr, objName, index)
//	return putstream, expectIps, nil
//}
//
//func chooseRandomNode(index int, nodesIndex []int) string {
//	dn := getDataNodes()
//	return dn[nodesIndex[index]]
//}
//
//func GetNodesIndex(n int) []int {
//	numbers := make([]int, n)
//	for i := 0; i < n; i++ {
//		numbers[i] = i
//	}
//	randomNumbers := RandomSelectNumbers(numbers)
//	return randomNumbers
//}
//
//func RandomSelectNumbers(numbers []int) []int {
//	rand.Seed(time.Now().UnixNano())
//	n := len(numbers)
//	result := make([]int, n)
//	copy(result, numbers)
//
//	for i := n - 1; i > 0; i-- {
//		j := rand.Intn(i + 1)
//		result[i], result[j] = result[j], result[i]
//	}
//	return result
//}
