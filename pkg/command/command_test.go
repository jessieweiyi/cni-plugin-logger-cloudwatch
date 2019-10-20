package command_test

import (
	"testing"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/command"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/config"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var testLoggerInstance = &testLogger{}

var stdinData = `
{
	"cniVersion": "0.3.1",
	"debug": true,
	"debugDir": "/var/vcap/data/cni-configs/net-debug",
	"logGroupName": "cni-plugin-logger-cloudwatch",
	"name": "mynet",
	"prevResult": {
		"cniVersion": "0.3.1",
		"interfaces": [{
			"name": "br0",
			"mac": "76:46:52:b7:75:b0"
		}, {
			"name": "vetheab91e52",
			"mac": "ce:6d:f9:bb:64:f3"
		}, {
			"name": "eth0",
			"mac": "9a:16:3a:18:42:a5",
			"sandbox": "/var/run/netns/cake"
		}],
		"ips": [{
			"version": "4",
			"interface": 2,
			"address": "10.11.0.7/24",
			"gateway": "10.11.0.1"
		}],
		"routes": [{
			"dst": "0.0.0.0/0"
		}, {
			"dst": "0.0.0.0/0",
			"gw": "10.11.0.1"
		}],
		"dns": {
			"nameservers": ["8.8.8.8"]
		}
	},
	"type": "cni-plugin-logger-cloudwatch"
}
`

var testArgs = &skel.CmdArgs{
	Netns:       "testNetNamespace",
	ContainerID: "testContainerID",
	IfName:      "testIfName",
	StdinData:   []byte(stdinData),
}

type testLoggerFactory struct {
}

type testLogger struct {
	isCalled bool
	args     *skel.CmdArgs
	action   string
}

func (tl *testLogger) Log(args *skel.CmdArgs, action string) {
	tl.isCalled = true
	tl.args = args
	tl.action = action
}

func (tlf *testLoggerFactory) GetLogger(conf *config.PluginConf, containerID string, ifName string) (logger.Logger, error) {
	return testLoggerInstance, nil
}

func TestCommand_CmdAdd_Success(t *testing.T) {
	resetTestLogger()
	command := command.NewCommand(&testLoggerFactory{})
	err := command.CmdAdd(testArgs)
	assert.NoError(t, err)
	assert.Equal(t, true, testLoggerInstance.isCalled)
	assert.Equal(t, testArgs, testLoggerInstance.args)
	assert.Equal(t, "add", testLoggerInstance.action)
}

func TestCommand_CmdDel_Success(t *testing.T) {
	resetTestLogger()
	command := command.NewCommand(&testLoggerFactory{})
	err := command.CmdDel(testArgs)
	assert.NoError(t, err)
	assert.Equal(t, true, testLoggerInstance.isCalled)
	assert.Equal(t, testArgs, testLoggerInstance.args)
	assert.Equal(t, "del", testLoggerInstance.action)
}

func resetTestLogger() {
	testLoggerInstance.isCalled = false
	testLoggerInstance.args = nil
	testLoggerInstance.action = ""
}
