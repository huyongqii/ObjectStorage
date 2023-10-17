package common

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"runtime"
)

// 执行命令的终端函数
func RunFn(usage string, cmds []cli.Command) {
	app := cli.NewApp()
	app.Name = "distributedObjectStorage"
	app.Version = "1.0.0"
	app.Description = "分布式对象存储系统"
	app.Author = "Yuxiang Wan"
	app.Copyright = "mgface@2021-∞"
	app.Usage = usage
	app.Email = "wanyuxiang000@163.com"
	//默认都添加stop命令
	app.Commands = append(
		cmds,
		StopCommand,
	)
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			//PrettyPrint: true,
		})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.InfoLevel)
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(false)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(fmt.Sprintf("发生fatal错误: %v", err))
	}
}

var stopflag = []cli.Flag{
	cli.IntFlag{
		Name:     "pid",
		Usage:    "暂停进程PID",
		Required: true,
	},
}
var StopCommand = cli.Command{
	Name:    "stop",
	Aliases: []string{"sp"},
	Usage: `示例指令:
		mgface.exe stop -pid 777
		`,
	Flags: stopflag,
	Action: func(ctx *cli.Context) error {
		pid := ctx.Int("pid")
		switch runtime.GOOS {
		case "windows":
			//获得暂停PID进程号
			cmd := exec.Command("taskkill", "/f", "/t", "/pid", fmt.Sprint(pid))
			err := cmd.Run()
			return err
		//因为本地开发环境是windows，所以这个代码在本地会显示错误，需要goland->settings->go->buildtags->os->linux
		case "linux":
		//err := syscall.Kill(pid, syscall.SIGSTOP)
		//return err
		default:
			return errors.New("为止的操作系统")
		}
		return nil
	},
}
