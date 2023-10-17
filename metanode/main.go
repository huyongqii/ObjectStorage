package main

import (
	"k8s.io/component-base/cli"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {
	cmd := NewCliCmd()
	if err := cli.RunNoErrOutput(cmd); err != nil {
		cmdutil.CheckErr(err)
	}
}
