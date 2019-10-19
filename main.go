package main

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/version"

	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/config"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/logger"

	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

// cmdAdd is called for ADD requests
func cmdAdd(args *skel.CmdArgs) error {
	conf, err := config.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}

	if conf.PrevResult == nil {
		return fmt.Errorf("must be called as chained plugin")
	}

	logger := logger.NewLogger(conf.Debug, conf.DebugDir, conf.LogGroupName)
	logger.Log(args, "add")

	// Pass through the result for the next plugin
	return types.PrintResult(conf.PrevResult, conf.CNIVersion)
}

// cmdDel is called for DELETE requests
func cmdDel(args *skel.CmdArgs) error {
	conf, err := config.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}
	_ = conf

	logger := logger.NewLogger(conf.Debug, conf.DebugDir, conf.LogGroupName)
	logger.Log(args, "del")

	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("CNI Plugin Cloud Watch"))
}

func cmdCheck(args *skel.CmdArgs) error {
	// TODO: implement
	return fmt.Errorf("not implemented")
}
