package command

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/config"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/logger"
)

// Command includes the command handlers for each cni request.
type Command struct {
	LoggerFactory logger.Factory
}

// NewCommand returns a new instance of command
func NewCommand(lf logger.Factory) *Command {
	return &Command{
		LoggerFactory: lf,
	}
}

// CmdAdd is called for ADD requests
func (c *Command) CmdAdd(args *skel.CmdArgs) error {
	conf, err := config.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}

	logger, err := c.LoggerFactory.GetLogger(conf, args.ContainerID, args.IfName)
	if err != nil {
		return fmt.Errorf("Failed to create logger")
	}

	logger.Log(args, "add")

	// Pass through the result for the next plugin
	return types.PrintResult(conf.PrevResult, conf.CNIVersion)
}

// CmdDel is called for DELETE requests
func (c *Command) CmdDel(args *skel.CmdArgs) error {
	conf, err := config.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}
	_ = conf

	logger, err := c.LoggerFactory.GetLogger(conf, args.ContainerID, args.IfName)
	if err != nil {
		return fmt.Errorf("Failed to create logger")
	}
	logger.Log(args, "del")

	return nil
}

// CmdCheck is called for CHECK requests
// TODO: implement
func (c *Command) CmdCheck(args *skel.CmdArgs) error {
	return fmt.Errorf("not implemented")
}
