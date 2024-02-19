package out

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
	display *log.Logger
)

const logPath = "~/.concert_manager/logs/"

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
	display = log.New(os.Stdout, "", 0)
	return nil
}

func createLogFile(fileName string) (*os.File, error) {
    if err := os.MkdirAll(logPath, 0750); err != nil {
		return nil, err
	}
	if err := os.Chdir(logPath); err != nil {
		return nil, err
	}
	file, err := os.Create(fileName)
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

func Fatalln(v ...any) {
    errorLog.Fatalln(v...)
}

func Panic(v ...any) {
	errorLog.Panic(v...)
}

func Panicf(format string, v ...any) {
    errorLog.Panicf(format, v...)
}

func Panicln(v ...any) {
    errorLog.Panicln(v...)
}

func Info(v ...any) {
	infoLog.Print(v...)
}

func Infof(format string, v ...any) {
	infoLog.Printf(format, v...)
}

func Infoln(v ...any) {
	infoLog.Println(v...)
}

func Debug(v ...any) {
	debugLog.Print(v...)
}

func Debugf(format string, v ...any) {
	debugLog.Printf(format, v...)
}

func Debugln(v ...any) {
	debugLog.Println(v...)
}

func Error(v ...any) {
	errorLog.Print(v...)
}

func Errorf(format string, v ...any) {
	errorLog.Printf(format, v...)
}

func Errorln(v ...any) {
	errorLog.Println(v...)
}

func Display(v ...any) {
	displayLog.Print(v...)
	display.Print(v...)
}

func Displayf(format string, v ...any) {
	displayLog.Printf(format, v...)
	display.Printf(format, v...)
}

func Displayln(v ...any) {
	displayLog.Println(v...)
	display.Println(v...)
}
