package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/timest/env"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"com.mgface.disobj/common"
	"com.mgface.disobj/common/k8s"
	"com.mgface.disobj/datanode/api"
	"com.mgface.disobj/datanode/server"
)

//todo 1.读取配置文件，覆盖默认值

//todo 2.根据配置文件来初始化文件存储还是缓存存储数据

//todo 3.判断如果是启用文件存储，启动的时候需要加载已有的文件数据

type Config struct {
	NodeAddr      string
	MetaNodeAddr  string
	StoreDataPath string
	PodNamespace  string
}

func NewCliCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datanode",
		Short: "datanode服务",
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
		Short: "启动datanode服务",
		Long: `示例指令:
		startmetanode start -na 127.0.0.1:3000 -gna 127.0.0.1:10000 -ca 127.0.0.1:10001,127.0.0.1:10000,127.0.0.1:10001,127.0.0.1:10002 -ms C:\\metadata
		`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("启动datanode服务")
			cmdutil.CheckErr(c.Complete())
			cmdutil.CheckErr(c.Validate())
			cmdutil.CheckErr(c.StartRun())
		},
	}
	cmd.Flags().StringVar(&c.NodeAddr, "n", "", "节点地址(node addres)")
	cmd.Flags().StringVar(&c.MetaNodeAddr, "m", "", "元数据服务节点地址(metanode address).只要配置一个可以正常连接上的种子节点即可")
	cmd.Flags().StringVar(&c.StoreDataPath, "s", "", "数据存储路径(store data path)")
	return cmd
}

func NewStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "停止datanode服务",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("停止datanode服务")
		},
	}
	return cmd
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) Complete() error {
	if common.IsRunningInKubernetes() {
		//-------解析k8s env
		cfg := new(k8s.EnvConfig)
		env.IgnorePrefix()
		err := env.Fill(cfg)
		if err != nil {
			log.Info("未读取到env环境.")
		} else {
			c.NodeAddr, c.MetaNodeAddr, c.StoreDataPath, c.PodNamespace = readEnv(cfg)
		}
	}

	//初始化数据
	api.Initval(c.StoreDataPath, c.NodeAddr)
	return nil
}

func (c *Config) StartRun() error {
	server.StartServer(c.NodeAddr, c.MetaNodeAddr, c.PodNamespace)
	return nil
}

func readEnv(cfg *k8s.EnvConfig) (na, mna, sdp, pns string) {

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

	//获取数据存储路径
	sdp = cfg.Sdp
	log.Info("读取env[sdp]数据:", sdp)

	//获得节点的命名空检
	pns = cfg.Pns
	log.Info("读取env[pns]数据:", pns)

	return
}
