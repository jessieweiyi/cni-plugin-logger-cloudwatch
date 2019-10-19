package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FilePublisher struct {
	DirPath string
}

func NewFilePublisher(dir string, containerID string, ifName string) (*FilePublisher, error) {
	dirPath := filepath.Join(dir, containerID, ifName)
	if err := os.MkdirAll(dirPath, 0770); err != nil {
		return nil, fmt.Errorf("Failed to create log folder")
	}

	return &FilePublisher{
		DirPath: dirPath,
	}, nil
}

func (fl *FilePublisher) Publish(cniLogData []byte) error {
	filePath := fmt.Sprintf("%s/%v.log", fl.DirPath, time.Now().Unix())

	if err := ioutil.WriteFile(filePath, cniLogData, 0770); err != nil {
		return fmt.Errorf("Failed to write log")
	}

	return nil
}
