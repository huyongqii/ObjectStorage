package main

import (
	"com.mgface.disobj/common"
	"com.mgface.disobj/common/k8s"
	"com.mgface.disobj/metanode/mq/mgfacemq"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/timest/env"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
)

type Config struct {
	NodeAddr      string
	GossipNode    string
	ClusterNode   string
	StoreDataPath string
	PodNamespace  string
	ServiceName   string
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
		Short: "启动metaNode服务",
		Long: `示例指令:
		startmetanode start -na 127.0.0.1:3000 -gna 127.0.0.1:10000 -ca 127.0.0.1:10001,127.0.0.1:10000,127.0.0.1:10001,127.0.0.1:10002 -ms C:\\metadata
		`,
		Run: func(cmd *cobra.Command, args []string) {
			log.SetOutput(os.Stdout)
			log.SetLevel(log.InfoLevel)
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(false)
			
			log.Info("启动datanode服务")
			cmdutil.CheckErr(c.Complete())
			cmdutil.CheckErr(c.Validate())
			cmdutil.CheckErr(c.StartRun())
		},
	}
	cmd.Flags().StringVarP(&c.NodeAddr, "node", "n", "", "节点地址(node addres)")
	cmd.Flags().StringVarP(&c.GossipNode, "gossip_node", "g", "", "Gossip节点地址(gossip node address).只要配置一个可以正常连接上的种子节点即可")
	cmd.Flags().StringVarP(&c.ClusterNode, "cluster_nodes", "c", "", "集群节点地址(cluster node address)")
	cmd.Flags().StringVarP(&c.StoreDataPath, "store_path", "s", "", "数据存储路径(store data path)")
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
			c.NodeAddr, c.GossipNode, c.ClusterNode, c.StoreDataPath, c.PodNamespace, c.ServiceName = readEnv(cfg)
		}
	}
	return nil
}

func (c *Config) StartRun() error {
	mgfacemq.StartEngine(c.NodeAddr, c.ClusterNode, c.GossipNode, c.StoreDataPath, c.PodNamespace, c.ServiceName)
	return nil
}

func readEnv(cfg *k8s.EnvConfig) (na, gna, ca, ms, pns, svcname string) {

	//获得节点的地址
	pns = cfg.Pns
	log.Info("读取env[pns]数据:", pns)

	//获得节点的地址
	na = cfg.Na
	log.Info("读取env[na]数据:", na)

	//从env获取到端口
	napt := cfg.Napt
	log.Info("读取env[napt]数据:", napt)
	if napt != "" {
		na = fmt.Sprintf("%s:%s", na, napt)
	}

	ca = cfg.Ca
	log.Info("读取env[ca]数据:", ca)

	gnapt := cfg.Gnapt
	log.Info("读取env[gnapt]数据:", gnapt)

	//把端口直接赋值给集群，后面去生成集群
	gna = gnapt
	//获取元数据存储路径
	ms = cfg.Ms
	log.Info("读取env[ms]数据:", ms)

	//获取SERVICE服务名称
	svcname = cfg.Svc
	log.Info("读取env[svc]数据:", svcname)
	return
}
