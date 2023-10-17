package command

import (
	. "com.mgface.disobj/common/k8s"
	"com.mgface.disobj/metanode/mq/mgfacemq"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/timest/env"
	"github.com/urfave/cli"
)

//todo 1.读取配置文件，覆盖默认值

//todo 2.根据配置文件来初始化文件存储还是缓存存储数据

//todo 3.判断如果是启用文件存储，启动的时候需要加载已有的文件数据

var startflag = []cli.Flag{
	cli.StringFlag{
		Name:     "na", //node address
		Usage:    "节点地址(node addres)",
		Required: false,
	},
	cli.StringFlag{
		Name:     "ca", //cluster address
		Usage:    "集群种子地址(cluster address).当节点只有自身时,把自己IP注册到gossip集群.",
		Required: false,
	},
	cli.StringFlag{
		Name:     "gna", //gossip node address
		Usage:    "集群节点goosip地址(gossip node address)",
		Required: false,
	},
	cli.StringFlag{
		Name:     "ms", //metadata store path
		Usage:    "元数据存储路径(metadata store path)",
		Required: false,
	},
	//todo 配置文件要单独处理
	cli.StringFlag{
		Name:     "file",
		Usage:    "配置文件",
		Required: false,
	},
}

var StartMNCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"st"},
	Usage: `示例指令:
		startmetanode start -na 127.0.0.1:3000 -gna 127.0.0.1:10000 -ca 127.0.0.1:10001,127.0.0.1:10000,127.0.0.1:10001,127.0.0.1:10002 -ms C:\\metadata
		`,
	Flags: startflag,
	Action: func(ctx *cli.Context) error {

		var nodeaddr, clusteraddr, gossipnodeaddr, metastorepath, podnamespace, svcname string
		//-------解析k8s env
		cfg := new(EnvConfig)
		env.IgnorePrefix()
		err := env.Fill(cfg)
		if err != nil {
			log.Info("未读取到env环境.", err.Error())
		} else {
			nodeaddr, clusteraddr, gossipnodeaddr, metastorepath, podnamespace, svcname = readEnv(cfg)
		}

		//假如pns为空，说明不是k8s环境
		if podnamespace == "" {
			nodeaddr, clusteraddr, gossipnodeaddr, metastorepath = readCli(ctx)
		}

		mgfacemq.StartEngine(nodeaddr, clusteraddr, gossipnodeaddr, metastorepath, podnamespace, svcname)
		return nil
	},
}

func readEnv(cfg *EnvConfig) (na, ca, gna, ms, pns, svcname string) {

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

// 读取命令行传入的参数(非k8s使用)
func readCli(ctx *cli.Context) (na, ca, gna, ms string) {
	//获得节点的地址
	na = ctx.String("na")
	log.Info("读取命令行参数[na]为:", na)

	//获得集群种子地址
	ca = ctx.String("ca")
	log.Info("读取命令行参数[ca]为:", ca)

	//获取集群节点goosip地址
	gna = ctx.String("gna")
	log.Info("读取命令行参数[gna]为:", gna)
	//获取元数据存储路径
	ms = ctx.String("ms")
	log.Info("读取命令行参数[ms]为:", ms)
	return
}
