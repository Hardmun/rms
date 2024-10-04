package logs

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type LogStruct struct {
	AbsPath string
	logType string
	LogFile *os.File
	Mtx     sync.RWMutex
}

var ErrLog = GetNewLog("errors.log", "error")
var InfoLog = GetNewLog("info.log", "info")

func GetLogDir(absPath string) (string, error) {
	loggingDir := filepath.Join(absPath, "logs")
	if info, errDir := os.Stat(loggingDir); errDir != nil || !info.IsDir() {
		if errDir = os.Mkdir(loggingDir, os.ModePerm); errDir != nil {
			return "", errDir
		}
	}
	return loggingDir, nil
}

func (l *LogStruct) initialize(filename string) error {
	var (
		loggingDir string
		lFile      *os.File
	)

	absPath, err := os.Getwd()
	if err != nil {
		return err
	}
	l.AbsPath = absPath

	loggingDir, err = GetLogDir(l.AbsPath)
	if err != nil {
		return err
	}
	lFile, err = os.OpenFile(filepath.Join(loggingDir, filename), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	l.LogFile = lFile
	l.Mtx = sync.RWMutex{}

	return nil
}

func (l *LogStruct) Write(errMsg ...any) {
	l.Mtx.Lock()
	defer l.Mtx.Unlock()
	newLine := log.New(l.LogFile, fmt.Sprintf("[%s]", l.logType), log.LstdFlags)
	newLine.Println(errMsg...)
}

func (l *LogStruct) Fatal(errMsg ...any) {
	l.Write(errMsg...)
	log.Fatal(errMsg...)
}

func (l *LogStruct) ErrorHTTP(w http.ResponseWriter, err string, s int) {
	l.Write(err)
	http.Error(w, err, s)
}

func (l *LogStruct) CloseLog() error {
	err := l.LogFile.Close()
	if err != nil {
		l.Write(err)
		return err
	}
	return nil
}

func GetNewLog(filename, logType string) *LogStruct {
	lg := &LogStruct{}
	lg.logType = logType
	if err := lg.initialize(filename); err != nil {
		log.Fatal(err)
	}

	return lg
}
