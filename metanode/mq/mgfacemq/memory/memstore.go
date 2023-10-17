package memory

import (
	"sync"

	"com.mgface.disobj/common"
)

// MemoryStore 实现Store接口
type MemoryStore struct {
	//存储的路径
	StorePath string `json:"store_path"`
	//互斥锁
	Mutex sync.RWMutex
	//存储数据的map
	Datas map[string]interface{} `json:"datas"`
	//是否开启后台定期清理数据
	EnableClean bool `json:"enable_clean"`
	//接收到的消息
	Msgs chan common.RecMsg `json:"msgs"`
	//快照生成完成的标志
	FishedSnapshot chan bool `json:"fished_snapshot"`
	//加载快照文件标识
	LoadingSnapshot bool `json:"loading_snapshot"`
	//缓冲的消息
	BuffMsgs chan common.RecMsg `json:"buff_msgs"`
	//缓冲的信号量
	BuffSemaphore chan bool `json:"buff_semaphore"`
}
