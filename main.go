package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/command"

	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

func main() {
	skel.PluginMain(command.CmdAdd, command.CmdCheck, command.CmdDel, version.All, bv.BuildString("CNI Plugin Cloud Watch"))
}
