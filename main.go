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
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"

	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

type PluginConf struct {
	types.NetConf
	RuntimeConfig *struct {
		SampleConfig map[string]interface{} `json:"sample"`
	} `json:"runtimeConfig"`

	RawPrevResult *map[string]interface{} `json:"prevResult"`
	PrevResult    *current.Result         `json:"-"`

	Debug    bool   `json:"debug"`
	DebugDir string `json:"debugDir"`

	LogGroupName string `json:"debugDir"`
}

type CNILogEntry struct {
	CommandType  string `json:"commandType"`
	NetNamespace string `json:"netNamespace"`
	Data         []byte `json:"data"`
}

func parseConfig(stdin []byte) (*PluginConf, error) {
	conf := PluginConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	// Parse previous result. Remove this if your plugin is not chained.
	if conf.RawPrevResult != nil {
		resultBytes, err := json.Marshal(conf.RawPrevResult)
		if err != nil {
			return nil, fmt.Errorf("could not serialize prevResult: %v", err)
		}
		res, err := version.NewResult(conf.CNIVersion, resultBytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse prevResult: %v", err)
		}
		conf.RawPrevResult = nil
		conf.PrevResult, err = current.NewResultFromResult(res)
		if err != nil {
			return nil, fmt.Errorf("could not convert result to current version: %v", err)
		}
	}
	// End previous result parsing

	if conf.Debug == true && conf.DebugDir == "" {
		return nil, fmt.Errorf("debugDir must be specified if debug is set to true")
	}

	if conf.LogGroupName == "" {
		return nil, fmt.Errorf("logGroupName must be specified")
	}

	return &conf, nil
}

// cmdAdd is called for ADD requests
func cmdAdd(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}

	if conf.PrevResult == nil {
		return fmt.Errorf("must be called as chained plugin")
	}

	debugError := cniDebug(conf.Debug, conf.DebugDir, args, "add")

	if debugError != nil {
		return debugError
	}

	// Pass through the result for the next plugin
	return types.PrintResult(conf.PrevResult, conf.CNIVersion)
}

// cmdDel is called for DELETE requests
func cmdDel(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}
	_ = conf

	debugError := cniDebug(conf.Debug, conf.DebugDir, args, "del")

	if debugError != nil {
		return debugError
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

func cniDebug(enabled bool, dir string, args *skel.CmdArgs, action string) error {
	if !enabled {
		return nil
	}
	dFilePath := filepath.Join(dir, args.ContainerID, args.IfName)
	if err := os.MkdirAll(dFilePath, 0770); err != nil {
		return fmt.Errorf("Failed to create log folder")
	}
	stdinFile := fmt.Sprintf("%s/%s_%v.log", dFilePath, action, time.Now().Unix())

	cniLog := CNILogEntry{
		NetNamespace: args.Netns,
		Data:         args.StdinData,
		CommandType:  action,
	}

	cniLogData, marshalError := json.MarshalIndent(cniLog, "", " ")

	if marshalError != nil {
		return fmt.Errorf("Failed to marshal log data")
	}

	if err := ioutil.WriteFile(stdinFile, cniLogData, 0770); err != nil {
		return fmt.Errorf("Failed to write log")
	}

	return nil
}
