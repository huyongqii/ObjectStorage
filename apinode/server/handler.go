package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"disobj/apinode/ops"
)

func ApiHandler(writer http.ResponseWriter, req *http.Request) {
	method := req.Method
	if method == http.MethodPut {
		ops.Put(writer, req)
		return
	}
	if method == http.MethodGet {
		ops.Get(writer, req)
		return
	}
	writer.WriteHeader(http.StatusMethodNotAllowed)
}

// LocateHandler 定位数据存储在哪个客户端
func LocateHandler(w http.ResponseWriter, req *http.Request) {
	m := req.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	nodeAddr, _, _ := Locate(strings.Split(req.URL.EscapedPath(), "/")[2], 3)
	if len(nodeAddr) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type", "application/json")
	b, _ := json.Marshal(nodeAddr)
	w.Write(b)
}
