package command

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/config"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/logger"
)

// CmdAdd is called for ADD requests
func CmdAdd(args *skel.CmdArgs) error {
	conf, err := config.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}

	if conf.PrevResult == nil {
		return fmt.Errorf("must be called as chained plugin")
	}

	logger, err := logger.GetLogger(conf, args.ContainerID, args.IfName)
	if err != nil {
		return fmt.Errorf("Failed to create logger")
	}

	logger.Log(args, "add")

	// Pass through the result for the next plugin
	return types.PrintResult(conf.PrevResult, conf.CNIVersion)
}

// CmdDel is called for DELETE requests
func CmdDel(args *skel.CmdArgs) error {
	conf, err := config.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}
	_ = conf

	logger, err := logger.GetLogger(conf, args.ContainerID, args.IfName)
	if err != nil {
		return fmt.Errorf("Failed to create logger")
	}
	logger.Log(args, "del")

	return nil
}

func CmdCheck(args *skel.CmdArgs) error {
	// TODO: implement
	return fmt.Errorf("not implemented")
}
