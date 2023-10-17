module com.mgface.disobj/datanode

go 1.15

require (
	com.mgface.disobj/common v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.7.0
	github.com/timest/env v0.0.0-20180717050204-5fce78d35255
	github.com/urfave/cli v1.22.5
	k8s.io/kubectl v0.27.3
)

replace com.mgface.disobj/common => ../common
