package k8s

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func client() *http.Client {

	pemData, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pemData)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	return &http.Client{Transport: tr}
}

var k8sToken = "/var/run/secrets/kubernetes.io/serviceaccount/token"

// DetectMetaService 从k8s中获取metanode的地址
// port: metanode的端口号
// namespace: k8s的namespace
// servicename: metanode的地址
func DetectMetaService(port, namespace, servicename string) []string {
	//------------------k8s operations----------------
	log.Debug("暂停30s,等待k8s进行更新数据......")
	time.Sleep(time.Duration(30+rand.Intn(5)) * time.Second)
	log.Debug("service name:", servicename, ",namespace:", namespace, ",port:", port)
	var request *http.Request

	search := fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/endpoints", namespace)

	request, _ = http.NewRequest("GET", search, nil)
	var token string

	fileInfo, _ := os.Lstat(k8sToken)
	log.Debug("k8sToken:", k8sToken, ",fileInfo.Mode():", fileInfo.Mode(), ",os.ModeSymlink:", os.ModeSymlink)

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		realPath, _ := filepath.EvalSymlinks(k8sToken)
		tokenbyte, _ := ioutil.ReadFile(realPath)
		token = string(tokenbyte)
		log.Debug("token:", token)
	}
	log.Debug("草泥马的token")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, _ := client().Do(request)
	log.Debug("resp:", resp)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("获取service值为:", string(body))
	//解析断点地址
	var endpoints v1.EndpointsList
	json.Unmarshal(body, &endpoints)

	metanodes := make([]string, 0)
	for _, v := range endpoints.Items {
		if strings.HasPrefix(servicename, v.Name) {
			for _, v1 := range v.Subsets {
				for _, v2 := range v1.Addresses {
					metanodes = append(metanodes, fmt.Sprintf("%s:%s", v2.IP, port))
					log.Debug("metanode:", fmt.Sprintf("%s:%s", v2.IP, port))
				}
			}
		}
	}
	return metanodes
}
