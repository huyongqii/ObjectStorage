package ops

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"com.mgface.disobj/datanode/api"
)

func Get(writer http.ResponseWriter, req *http.Request) {
	url := req.URL.EscapedPath()
	rollingDateLSPath := strings.Split(url, "/")[2]
	objName := strings.Split(url, "/")[3]
	realObjName := fmt.Sprintf("%s/%s", rollingDateLSPath, objName)
	filename := api.GetRootLocalStorePath() + string(os.PathSeparator) + realObjName
	file, err := os.Open(filename)
	if err != nil {
		log.Debug(err)
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	io.Copy(writer, file)
}
