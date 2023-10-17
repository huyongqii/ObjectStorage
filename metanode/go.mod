module com.mgface.disobj/metanode

go 1.15

require (
	com.mgface.disobj/common v0.0.0-00010101000000-000000000000
	github.com/hashicorp/memberlist v0.2.2
	github.com/pborman/uuid v1.2.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.6.0
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/timest/env v0.0.0-20180717050204-5fce78d35255
	github.com/urfave/cli v1.22.5
	k8s.io/component-base v0.27.3
	k8s.io/kubectl v0.27.3
)

replace com.mgface.disobj/common => ../common
