package logger

import "github.com/containernetworking/cni/pkg/skel"

type CNILogEntry struct {
	CommandType  string `json:"commandType"`
	NetNamespace string `json:"netNamespace"`
	StdinData    string `json:"stdinData"`
}

func NewCNILogEntry(args *skel.CmdArgs, action string) CNILogEntry {
	return CNILogEntry{
		NetNamespace: args.Netns,
		StdinData:    string(args.StdinData),
		CommandType:  action,
	}
}
