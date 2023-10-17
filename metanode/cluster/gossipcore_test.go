package cluster

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
)

var (
	mtx        sync.RWMutex
	members    = flag.String("members", "127.0.0.1:14000", "comma seperated list of members")
	port       = flag.Int("port", 5001, "http port")
	items      = map[string]string{}
	broadcasts *memberlist.TransmitLimitedQueue
)

func init() {
	flag.Parse()
}

type broadcastTemp struct {
	msg    []byte
	notify chan<- struct{}
}

type delegateTemp struct{}

type updateTemp struct {
	Action string // add, del
	Data   map[string]string
}

func (b *broadcastTemp) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcastTemp) Message() []byte {
	return b.msg
}

func (b *broadcastTemp) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func (d *delegateTemp) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegateTemp) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}

	switch b[0] {
	case 'd': // data
		var updates []*updateTemp
		if err := json.Unmarshal(b[1:], &updates); err != nil {
			return
		}
		mtx.Lock()
		for _, u := range updates {
			for k, v := range u.Data {
				switch u.Action {
				case "add":
					items[k] = v
				case "del":
					delete(items, k)
				}
			}
		}
		mtx.Unlock()
	}
}

func (d *delegateTemp) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegateTemp) LocalState(join bool) []byte {
	mtx.RLock()
	m := items
	mtx.RUnlock()
	b, _ := json.Marshal(m)
	return b
}

func (d *delegateTemp) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
	var m map[string]string
	if err := json.Unmarshal(buf, &m); err != nil {
		return
	}
	mtx.Lock()
	for k, v := range m {
		items[k] = v
	}
	mtx.Unlock()
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	val := r.Form.Get("val")
	mtx.Lock()
	items[key] = val
	mtx.Unlock()

	b, err := json.Marshal([]*updateTemp{
		&updateTemp{
			Action: "add",
			Data: map[string]string{
				key: val,
			},
		},
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//广播数据
	broadcasts.QueueBroadcast(&broadcastTemp{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	mtx.Lock()
	delete(items, key)
	mtx.Unlock()

	b, err := json.Marshal([]*updateTemp{
		&updateTemp{
			Action: "del",
			Data: map[string]string{
				key: "",
			},
		},
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	broadcasts.QueueBroadcast(&broadcastTemp{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	mtx.RLock()
	val := items[key]
	mtx.RUnlock()
	w.Write([]byte(val))
}

func start() error {
	hostname, _ := os.Hostname()
	c := memberlist.DefaultLANConfig()
	c.Delegate = &delegateTemp{}
	c.BindPort = 24000
	c.BindAddr = "127.0.0.1"
	c.Name = hostname + "-" + uuid.NewUUID().String()
	//创建gossip网络
	m, err := memberlist.Create(c)
	if err != nil {
		return err
	}
	//第一个节点没有member，但从第二个开始就有member了
	if len(*members) > 0 {
		parts := strings.Split(*members, ",")
		_, err := m.Join(parts)
		if err != nil {
			return err
		}
	}
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}
	fmt.Printf("Local member %s:%d\n", c.BindAddr, c.BindPort)
	return nil
}

func TestNew(t *testing.T) {
	if err := start(); err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/del", delHandler)
	http.HandleFunc("/get", getHandler)
	fmt.Printf("Listening on :%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Println(err)
	}
}
