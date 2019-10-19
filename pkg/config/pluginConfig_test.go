package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/config"
)

func TestConfig_NewConfig_Success(t *testing.T) {
	stdinData := "{\"cniVersion\":\"0.3.1\",\"debug\":true,\"debugDir\":\"/var/vcap/data/cni-configs/net-debug\",\"logGroupName\":\"cni-plugin-logger-cloudwatch\",\"name\":\"mynet\",\"prevResult\":{\"cniVersion\":\"0.3.1\",\"interfaces\":[{\"name\":\"br0\",\"mac\":\"76:46:52:b7:75:b0\"},{\"name\":\"vetheab91e52\",\"mac\":\"ce:6d:f9:bb:64:f3\"},{\"name\":\"eth0\",\"mac\":\"9a:16:3a:18:42:a5\",\"sandbox\":\"/var/run/netns/cake\"}],\"ips\":[{\"version\":\"4\",\"interface\":2,\"address\":\"10.11.0.7/24\",\"gateway\":\"10.11.0.1\"}],\"routes\":[{\"dst\":\"0.0.0.0/0\"},{\"dst\":\"0.0.0.0/0\",\"gw\":\"10.11.0.1\"}],\"dns\":{\"nameservers\":[\"8.8.8.8\"]}},\"type\":\"cni-plugin-logger-cloudwatch\"}"
	_, err := config.ParseConfig([]byte(stdinData))
	assert.NoError(t, err)
}
