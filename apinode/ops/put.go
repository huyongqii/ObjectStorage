package ops

import (
	"net/http"
	"strings"
	
	log "github.com/sirupsen/logrus"
)

func Put(w http.ResponseWriter, req *http.Request) {
	objName := strings.Split(req.URL.EscapedPath(), "/")[2]
	// 往http://nodeAddr/objects/objName存储数据
	status, err := storeObject(req.Body, objName)
	w.WriteHeader(status)
	if err != nil {
		log.Info(err)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
	}
}
