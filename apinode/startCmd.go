package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/timest/env"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"disobj/apinode/server"
	"disobj/common"
	"disobj/common/k8s"
)

type Config struct {
	NodeAddr     string
	MetaNodeAddr string
	DataShards   int
	ParityShards int
	PodNameSpace string
}

func NewCliCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apiNode",
		Short: "apiNode服务",
	}
	cmd.AddCommand(
		NewStartCmd(),
		NewStopCmd(),
	)
	return cmd
}

func NewStartCmd() *cobra.Command {
	c := &Config{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "启动apiNode服务",
		Run: func(cmd *cobra.Command, args []string) {
			log.SetFormatter(&log.TextFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
				//PrettyPrint: true,
			})
			log.SetOutput(os.Stdout)
			log.SetLevel(log.InfoLevel)
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(false)

			cmdutil.CheckErr(c.Complete())
			cmdutil.CheckErr(c.Validate())
			cmdutil.CheckErr(c.StartRun())
		},
	}
	cmd.Flags().StringVar(&c.NodeAddr, "n", "", "节点地址(node addres)")
	cmd.Flags().StringVar(&c.MetaNodeAddr, "m", "", "元数据服务节点地址(metanode address).只要配置一个可以正常连接上的种子节点即可")
	cmd.Flags().IntVar(&c.DataShards, "ds", 2, "数据分片大小数量(datashards number)")
	cmd.Flags().IntVar(&c.ParityShards, "ps", 1, "奇偶校验数量(parityshards number)")

	return cmd
}

func NewStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "停止apiNode服务",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(StopRun())
		},
	}
	return cmd
}

func (c *Config) Complete() error {
	var nodeAddr, metaNodeAddr string
	var dataShards, parityShards int

	//-------解析k8s env
	cfg := new(k8s.EnvConfig)
	env.IgnorePrefix()
	err := env.Fill(cfg)
	if err != nil {
		log.Info("未读取到env环境.")
	} else {
		nodeAddr, metaNodeAddr, dataShards, parityShards, c.PodNameSpace = readEnv(cfg)
	}
	// k8s环境下，如果没有配置元数据服务地址，则从k8s的service中获取
	if c.PodNameSpace != "" {
		c.NodeAddr, c.MetaNodeAddr, c.DataShards, c.ParityShards = nodeAddr, metaNodeAddr, dataShards, parityShards
	}

	//初始化数据
	server.Initval(c.DataShards, c.ParityShards)
	return nil
}

func (c *Config) Validate() error {
	fmt.Println("nodeAddr: ", c.NodeAddr)
	fmt.Println("metaNodeAddr: ", c.MetaNodeAddr)
	fmt.Println("dataShards: ", c.DataShards)
	fmt.Println("parityShards: ", c.ParityShards)
	return nil
}

func (c *Config) StartRun() error {
	//启动服务
	log.Debug("启动API节点...")
	log.Debug(fmt.Sprintf("节点地址:%s", c.NodeAddr))
	log.Debug(fmt.Sprintf("元数据服务节点地址:%s", c.MetaNodeAddr))

	// 创建2个启动标志，一个用来启动发送心跳服务，一个用来更新数据节点数据
	startFlag := make(chan bool, 2)
	go server.RefreshDynamicMetaNode(c.MetaNodeAddr, c.PodNameSpace, startFlag)
	go server.StartApiHeartbeat(c.NodeAddr, startFlag)
	go server.RefreshDNData(startFlag)

	// 监听请求
	http.HandleFunc("/objects/", server.ApiHandler)
	http.HandleFunc("/locate/", server.LocateHandler)

	common.SupportServeAndGracefulExit(c.NodeAddr)
	return nil
}

func StopRun() error {
	return nil
}

func readEnv(cfg *k8s.EnvConfig) (na, mna string, ds, ps int, pns string) {

	//从env获取节点的地址
	na = cfg.Na
	log.Info("读取env[na]数据:", na)

	//从env获取到端口
	napt := cfg.Napt
	log.Info("读取env[napt]数据:", napt)
	if napt != "" {
		na = fmt.Sprintf("%s:%s", na, napt)
	}
	//从env获取到集群地址
	ca := cfg.Ca
	log.Info("读取env[ca]数据:", ca)
	//从env获取到集群地址的默认端口
	capt := cfg.Capt
	log.Info("读取env[capt]数据:", capt)

	mna = fmt.Sprintf("%s:%s", ca, capt)
	log.Info("metanode: ", mna)

	//获取数据分片大小数量
	ds, _ = strconv.Atoi(cfg.Ds)
	log.Info("读取env[ds]数据:", ds)

	//获取奇偶校验数量
	ps, _ = strconv.Atoi(cfg.Ps)
	log.Info("读取env[ps]数据:", ps)

	//获得节点的命名空检
	pns = cfg.Pns
	log.Info("读取env[pns]数据:", pns)

	return
}
