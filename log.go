package services

import (
	"fmt"
	"os"
	"time"
)

var (
	filename string
)

const (
	RuntimePath string = "runtime"
	LogPath     string = RuntimePath + "/logs"
)

type LogService struct {
	pidFile string
	logFile string
}

func NewLogService(scriptName string) *LogService {
	return &LogService{
		pidFile: scriptName + ".pid",
		logFile: scriptName + ".log",
	}
}

func (l *LogService) Init() error {
	if _, err := os.Stat(RuntimePath); os.IsNotExist(err) {
		if err := os.Mkdir(RuntimePath, os.ModePerm); err != nil {
			return err
		}
	}

	if _, err := os.Stat(LogPath); os.IsNotExist(err) {
		if err := os.Mkdir(LogPath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (l *LogService) Log(a ...interface{}) {
	now := time.Now()
	date := now.Format(YYYYMMDDHHMMSS)
	fmt.Print("[" + date + "] ")
	fmt.Println(a...)
}

func (t *LogService) GetPidFile() string {
	return RuntimePath + "/" + t.pidFile
}

func (t *LogService) GetLogFile() string {
	return LogPath + "/" + t.logFile
}
