package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FileLogger struct {
	DirPath string
}

func NewFileLogger(dir string, containerID string, ifName string) (*FileLogger, error) {
	dirPath := filepath.Join(dir, containerID, ifName)
	if err := os.MkdirAll(dirPath, 0770); err != nil {
		return nil, fmt.Errorf("Failed to create log folder")
	}

	return &FileLogger{
		DirPath: dirPath,
	}, nil
}

func (fl *FileLogger) Log(cniLogData []byte) error {
	filePath := fmt.Sprintf("%s/%v.log", fl.DirPath, time.Now().Unix())

	if err := ioutil.WriteFile(filePath, cniLogData, 0770); err != nil {
		return fmt.Errorf("Failed to write log")
	}

	return nil
}
