package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	debugLog *log.Logger
	infoLog *log.Logger
	errorLog *log.Logger
	displayLog *log.Logger
)

const logPath = "/.concert_manager/logs/"

func Initialize() error {
	fileName := fmt.Sprintf("log-%d.log", time.Now().UnixMilli())
	logFile, err := createLogFile(fileName)
	if err != nil {
		return err
	}

	debugLog = log.New(logFile, "debug   - ", log.LstdFlags)
	infoLog = log.New(logFile, "info    - ", log.LstdFlags)
	errorLog = log.New(logFile, "error  -  ", log.LstdFlags)
	displayLog = log.New(logFile, "display - ", log.LstdFlags)
	return nil
}

func createLogFile(fileName string) (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	logDir := homeDir + logPath
    if err := os.MkdirAll(logDir, 0750); err != nil {
		return nil, err
	}
	filePath := logDir + fileName
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, err
}

func Fatal(v ...any) {
	errorLog.Fatal(v...)
}

func Fatalf(format string, v ...any) {
    errorLog.Fatalf(format, v...)
}

func Panic(v ...any) {
	errorLog.Panic(v...)
}

func Panicf(format string, v ...any) {
    errorLog.Panicf(format, v...)
}

func Info(v ...any) {
	infoLog.Print(v...)
}

func Infof(format string, v ...any) {
	infoLog.Printf(format, v...)
}

func Debug(v ...any) {
	debugLog.Print(v...)
}

func Debugf(format string, v ...any) {
	debugLog.Printf(format, v...)
}

func Error(v ...any) {
	errorLog.Print(v...)
}

func Errorf(format string, v ...any) {
	errorLog.Printf(format, v...)
}

func Display(v ...any) {
	displayLog.Print(v...)
}

func Displayf(format string, v ...any) {
	displayLog.Printf(format, v...)
}
