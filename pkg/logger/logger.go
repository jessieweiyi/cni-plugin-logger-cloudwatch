package logger

import (
	"github.com/containernetworking/cni/pkg/skel"
)

type Logger interface {
	Log(args *skel.CmdArgs, action string)
}
