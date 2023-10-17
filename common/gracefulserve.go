package common

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//启动服务，并且支持优雅退出服务，可以响应终端
func SupportServeAndGracefulExit(nodeAddr string) {

	s := &http.Server{
		Addr:           nodeAddr,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Debug(s.ListenAndServe())
		log.Debug("server shutdown.")
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Debug(<-ch)

	//优雅的停止服务.
	s.Shutdown(context.Background())

	// 等待go routine停止完至后结束
	log.Debug("done.")
}
