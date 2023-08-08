package receiver

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileReceiver struct {
	file *os.File
}

func (f FileReceiver) Notify(message string) error {
	lineBytes := []byte(fmt.Sprintf("%s\n", message))
	_, err := f.file.Write(lineBytes)
	return err
}

func (f FileReceiver) Type() string {
	return "file"
}

func NewFileReceiver(settings map[string]interface{}) (*FileReceiver, error) {
	name, ok := settings["name"]
	if !ok {
		return nil, fmt.Errorf("file must have a name")
	}

	nameStr, ok := name.(string)
	if !ok {
		return nil, fmt.Errorf("file name must be a string, found: %+v", name)
	}

	targetDir, ok := settings["directory"]
	if !ok {
		return nil, fmt.Errorf("target directory must be specified")
	}

	targetDirStr, ok := targetDir.(string)
	if !ok {
		return nil, fmt.Errorf("directory must be a string found %+v", targetDir)
	}

	absDir, err := filepath.Abs(targetDirStr)
	if err != nil {
		return nil, err
	}

	targetFilePath := filepath.Join(absDir, nameStr)
	file, err := os.OpenFile(targetFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileReceiver{file: file}, nil
}
