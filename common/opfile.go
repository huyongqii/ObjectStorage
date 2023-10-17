package common

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//保存文件
func Savefile(fname string, fsize int, reader *bufio.Reader) (err error) {

	if fsize ==0 {
		return nil
	}
	
	dstFile, e := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if e != nil {
		log.Warn("打开文件:", fname, "失败,异常:", e.Error())
		return e
	}

	writer := bufio.NewWriter(dstFile)

	buffer := make([]byte, 4096)
	total := 0
	for {
		//接受客户端上传的文件
		n, _ := reader.Read(buffer)
		total += n
		//写入服务端本地文件
		writer.Write(buffer[:n])
		writer.Flush()

		//如果实际总接受字节数与客户端给的要传输字节数相等，说明传输完毕
		if total == fsize {
			log.Info(fmt.Sprintf("文件:%s 接受成功,共%d字节", fname, total))
			break
		}
	}

	dstFile.Close()

	return nil
}

//发送文件
func Sendfile(fd FileDesc, client *TcpClient) string {
	//1.先写文件的FD到客户端
	data := []byte(fmt.Sprintf("F%s %d ", fd.Filename, fd.FD.Size()))
	client.Conn.Write(data)
	//2.再写实际文件数据过去
	srcFile, _ := os.OpenFile(fd.Fpath, os.O_RDONLY, 0644)
	defer srcFile.Close()
	reader := bufio.NewReader(srcFile)
	buffer := make([]byte, 4096)
	total := 0
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			log.Info(fmt.Sprintf("%s文件发送完毕,大小:%d字节", fd.Filename, total))
			break
		} else {
			client.Conn.Write(buffer[:n])
			total += n
		}
	}
	recode := make([]byte, 4096)
	n, _ := client.Conn.Read(recode)
	replyMsg := string(recode[:n])

	return replyMsg
}

//遍历文件夹
func WalkDirectory(storepath string) []FileDesc {
	fds := make([]FileDesc, 0)
	filepath.Walk(storepath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			//log.Debug("dir:", path)
			return nil
		}
		if strings.Contains(f.Name(), ".prof") {
			return nil
		}
		fd := FileDesc{
			Filename: f.Name(),
			Fpath:    path,
			FD:       f,
		}
		fds = append(fds, fd)
		return nil
	})
	return fds
}
