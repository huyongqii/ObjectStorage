package k8s

type EnvConfig struct {
	Pns   string `env:"POD_NAMESPACE"`
	Na    string `env:"NODE_ADDR"`
	Napt  string `env:"NODE_PORT",default:"3000"`
	Ca    string `env:"CLUSTER_ADDRS"`
	Capt  string `env:"CLUSTER_ADDRS_PORT"` //`default:""`
	Sdp   string `env:"STORE_DATA_PATH"`
	Ms    string `env:"METADATA_STORE"`
	Svc   string `env:"SERVICE_NAME"`
	Gnapt string `env:"GOSSIP_NODE_PORT"` //`default:""`
	Ds    string `env:"DS"`
	Ps    string `env:"PS"`
	//Hosts   []string `slice_sep:","`
}
