package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/version"

	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/config"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/logger"

	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

type CNILogEntry struct {
	CommandType  string `json:"commandType"`
	NetNamespace string `json:"netNamespace"`
	StdinData    string `json:"stdinData"`
}

// cmdAdd is called for ADD requests
func cmdAdd(args *skel.CmdArgs) error {
	conf, err := config.ParseConfig(args.StdinData)
	if err != nil {
		return err
	}

	if conf.PrevResult == nil {
		return fmt.Errorf("must be called as chained plugin")
	}

	cniLog := createLogEntry(args, "add")
	debugError := cniDebug(conf.Debug, conf.DebugDir, args.ContainerID, args.IfName, cniLog)

	if debugError != nil {
		return debugError
	}

	logError := cniLogCloudWatch(conf.LogGroupName, args.ContainerID, args.IfName, cniLog)

	if logError != nil {
		return logError
	}

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

	cniLog := createLogEntry(args, "del")
	debugError := cniDebug(conf.Debug, conf.DebugDir, args.ContainerID, args.IfName, cniLog)

	if debugError != nil {
		return debugError
	}

	logError := cniLogCloudWatch(conf.LogGroupName, args.ContainerID, args.IfName, cniLog)

	if logError != nil {
		return logError
	}

	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("CNI Plugin Cloud Watch"))
}

func cmdCheck(args *skel.CmdArgs) error {
	// TODO: implement
	return fmt.Errorf("not implemented")
}

func createLogEntry(args *skel.CmdArgs, action string) CNILogEntry {
	return CNILogEntry{
		NetNamespace: args.Netns,
		StdinData:    string(args.StdinData),
		CommandType:  action,
	}
}

func cniDebug(enabled bool, dir string, containerID string, ifName string, cniLog CNILogEntry) error {
	if !enabled {
		return nil
	}
	dFilePath := filepath.Join(dir, containerID, ifName)
	if err := os.MkdirAll(dFilePath, 0770); err != nil {
		return fmt.Errorf("Failed to create log folder")
	}
	stdinFile := fmt.Sprintf("%s/%v.log", dFilePath, time.Now().Unix())

	cniLogData, marshalError := json.MarshalIndent(cniLog, "", " ")

	if marshalError != nil {
		return fmt.Errorf("Failed to marshal log data")
	}

	if err := ioutil.WriteFile(stdinFile, cniLogData, 0770); err != nil {
		return fmt.Errorf("Failed to write log")
	}

	return nil
}

func cniLogCloudWatch(logGroupName string, containerID string, ifName string, cniLog CNILogEntry) error {
	logStreamName := fmt.Sprintf("/%s/%s", containerID, ifName)

	cniLogData, marshalError := json.MarshalIndent(cniLog, "", " ")

	if marshalError != nil {
		return fmt.Errorf("Failed to marshal log data")
	}

	if err := logger.Log(logGroupName, logStreamName, string(cniLogData)); err != nil {
		return fmt.Errorf("Failed to publish log to aws cloudwatch")
	}

	return nil
}
