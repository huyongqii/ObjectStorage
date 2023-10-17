package file

import (
	. "com.mgface.disobj/common"
	"com.mgface.disobj/metanode/mq/mgfacemq/memory"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

func TestStoresnapshotData(t *testing.T) {
	memoryStore := &memory.MemoryStore{
		StorePath:     "c:\\metadata",
		Mutex:         sync.RWMutex{},
		Datas:         make(map[string]interface{}),
		EnableClean:   true,
		Msgs:          make(chan RecMsg, 1_000),
		BuffMsgs:      make(chan RecMsg, 5_000),
		BuffSemaphore: make(chan bool),
	}
	Msg := make(chan RecMsg, 1_000)
	go func() {
		for i := 0; i < 100; i++ {
			Msg <- RecMsg{
				Key: fmt.Sprintf("%s%d", "EEE", i),
				Val: fmt.Sprintf("%s%d", "ABC", i),
			}
		}
	}()
	StoresnapshotData(memoryStore)
	time.Sleep(100 * time.Second)
}

func TestBytesToInt(t *testing.T) {
	f, _ := os.OpenFile("C:\\metadata\\20210214.log", os.O_RDONLY, 0644)
	for {
		dd := make([]byte, 4)
		n, e := f.Read(dd)
		fmt.Println("n==>", n)
		if e == io.EOF {
			fmt.Println("read the file finished")
			break
		}
		max := BytesToInt(dd)
		t.Log("max=", max)
		ev := make([]byte, max)
		n, e = f.Read(ev)
		fmt.Println(n, e)
		t.Log(string(ev))
	}
	fmt.Println("end...")
}
