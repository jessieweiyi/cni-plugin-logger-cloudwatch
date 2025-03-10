package logger

import (
	"fmt"

	"github.com/jessieweiyi/cni-plugin-logger-cloudwatch/pkg/config"
)

// Factory interfact for all logger factories
type Factory interface {
	GetLogger(conf *config.PluginConf, containerID string, ifName string) (Logger, error)
}

// CNILoggerFactory returns a instance of CNILogger
type CNILoggerFactory struct {
}

// GetLogger returns a instance of CNILogger
func (clf *CNILoggerFactory) GetLogger(conf *config.PluginConf, containerID string, ifName string) (Logger, error) {
	var publishers []Publisher
	if conf.Debug {
		filePublisher, flError := NewFilePublisher(conf.DebugDir, containerID, ifName)
		if flError != nil {
			return nil, fmt.Errorf("Error: Failed to create file loader")
		}
		publishers = append(publishers, filePublisher)
	}

	cloudWatchLogPublisher, cwlError := NewCloudWatchLogsPublisher(conf.LogGroupName, containerID, ifName)
	if cwlError != nil {
		return nil, fmt.Errorf("Failed to create cloudwatch loader")
	}

	publishers = append(publishers, cloudWatchLogPublisher)

	return NewCNILogger(publishers), nil
}
