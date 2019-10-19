package logger

import (
	"encoding/json"
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
)

type CNILogger struct {
	Publishers []Publisher
}

func NewCNILogger(publishers []Publisher) Logger {
	return &CNILogger{
		Publishers: publishers,
	}
}

func (l *CNILogger) Log(args *skel.CmdArgs, action string) {
	cniLog := NewCNILogEntry(args, action)

	cniLogData, marshalError := json.MarshalIndent(cniLog, "", " ")

	if marshalError != nil {
		fmt.Println("Error: Failed to marshal log data", marshalError)
		return
	}

	for _, publisher := range l.Publishers {
		if err := publisher.Publish(cniLogData); err != nil {
			fmt.Printf("Error: Failed to publisher log data to %T\n", publisher)
		}
	}
}
