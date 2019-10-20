package logger_test

import (
	"testing"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type testPublisher struct {
	isCalled bool
}

func (p *testPublisher) Publish(cniLogData []byte) error {
	p.isCalled = true
	return nil
}

func TestCNILogger_Log_Success(t *testing.T) {
	testPublisher1 := testPublisher{}
	testPublisher2 := testPublisher{}
	var publishers []logger.Publisher
	publishers = append(publishers, &testPublisher1)
	publishers = append(publishers, &testPublisher2)
	logger := logger.NewCNILogger(publishers)

	args := &skel.CmdArgs{
		Netns:       "testNetNamespace",
		ContainerID: "testContainerID",
		IfName:      "testIfName",
		StdinData:   []byte("testData"),
	}

	logger.Log(args, "test")
	assert.Equal(t, true, testPublisher1.isCalled)
	assert.Equal(t, true, testPublisher2.isCalled)
}
