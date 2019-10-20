package logger

import "github.com/containernetworking/cni/pkg/skel"

// CNILogEntry lists the log content
type CNILogEntry struct {
	CommandType  string `json:"commandType"`
	NetNamespace string `json:"netNamespace"`
	StdinData    string `json:"stdinData"`
}

// NewCNILogEntry returns a instance of CNILogEntry
func NewCNILogEntry(args *skel.CmdArgs, action string) CNILogEntry {
	return CNILogEntry{
		NetNamespace: args.Netns,
		StdinData:    string(args.StdinData),
		CommandType:  action,
	}
}
