package logger

import (
	"encoding/json"
	"fmt"

	"github.com/containernetworking/cni/pkg/skel"
)

type Logger struct {
	Debug        bool
	DebugDir     string
	LogGroupName string
}

func NewLogger(debug bool, debugDir string, logGroupName string) *Logger {
	return &Logger{
		Debug:        debug,
		DebugDir:     debugDir,
		LogGroupName: logGroupName,
	}
}

func (l *Logger) Log(args *skel.CmdArgs, action string) {
	cniLog := NewCNILogEntry(args, action)

	cniLogData, marshalError := json.MarshalIndent(cniLog, "", " ")

	if marshalError != nil {
		fmt.Println("Error: Failed to marshal log data", marshalError)
		return
	}

	if l.Debug {
		l.logToFile(args.ContainerID, args.IfName, cniLogData)
	}

	l.logToCloudWatch(args.ContainerID, args.IfName, cniLogData)
}

func (l *Logger) logToFile(containerID string, ifName string, cniLogData []byte) {
	fileLogger, flError := NewFileLogger(l.DebugDir, containerID, ifName)
	if flError != nil {
		fmt.Println("Error: Failed to create file loader", flError)
		return
	}

	flLogError := fileLogger.Log(cniLogData)

	if flLogError != nil {
		fmt.Println("Error: Failed to log to file", flLogError)
	}
}

func (l *Logger) logToCloudWatch(containerID string, ifName string, cniLogData []byte) {
	cloudWatchLogger, cwlError := NewCloudWatchLogger(l.LogGroupName, containerID, ifName)
	if cwlError != nil {
		fmt.Println("Error: Failed to create cloudwatch loader", cwlError)
		return
	}

	cwlLogError := cloudWatchLogger.Log(cniLogData)

	if cwlLogError != nil {
		fmt.Println("Error: Failed to log to cloudwatch", cwlLogError)
	}
}
