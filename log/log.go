package log

import (
	"log"
	"os"
)

var Logger *log.Logger

func New() *log.Logger {
	return log.New(os.Stdout, "oreshell", log.LstdFlags)
}

func NewForFile(filePath string) (*log.Logger, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return log.New(f, "oreshell", log.LstdFlags), nil
}
