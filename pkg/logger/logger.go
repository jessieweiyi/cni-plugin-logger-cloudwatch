package logger

import (
	"github.com/containernetworking/cni/pkg/skel"
)

// Logger interface for all loggers
type Logger interface {
	Log(args *skel.CmdArgs, action string)
}
