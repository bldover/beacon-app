package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	debugLog   *log.Logger
	infoLog    *log.Logger
	errorLog   *log.Logger
	alertLog   *log.Logger
	displayLog *log.Logger
	alerter    *EmailAlerter
	logLevels  = [...]string{"DEBUG", "INFO", "ERROR"}
)

const logLevelEnv = "CM_LOG_LEVEL"

type loggerLevel int

const (
	debugLevel loggerLevel = iota
	infoLevel
	errorLevel
	defaultLevel
)

func Initialize() error {
	logFile, err := createLogFile()
	if err != nil {
		return err
	}

	displayLog = log.New(logFile, "display - ", log.LstdFlags)

	logLevel := getLogLevel()
	unsetLogLevel := logLevel == defaultLevel
	if unsetLogLevel {
		logLevel = infoLevel
	}

	if logLevel <= debugLevel {
		debugLog = log.New(logFile, "debug   - ", log.LstdFlags)
	}
	if logLevel <= infoLevel {
		infoLog = log.New(logFile, "info    - ", log.LstdFlags)
	}
	errorLog = log.New(logFile, "error  -  ", log.LstdFlags)
	alertLog = log.New(logFile, "alert  -  ", log.LstdFlags)
	if unsetLogLevel {
		Info("Unexpected or missing value for CM_LOG_LEVEL environment variable, defaulting to INFO level")
	}

	alerter, err = NewGmailAlerter()
	if err != nil {
		return err
	}

	Info("Successfully initialized logger with level", logLevels[logLevel])
	return nil
}

func getLogLevel() loggerLevel {
	in := os.Getenv(logLevelEnv)
	for level, levelName := range logLevels {
		if levelName == in {
			return loggerLevel(level)
		}
	}
	return defaultLevel
}

func createLogFile() (*os.File, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	filePath := executable + ".log"
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, err
}

func Fatal(v ...any) {
	errorLog.Fatalln(v...)
}

func Fatalf(format string, v ...any) {
	errorLog.Fatalf(format, v...)
}

func Panic(v ...any) {
	errorLog.Panicln(v...)
}

func Panicf(format string, v ...any) {
	errorLog.Panicf(format, v...)
}

func Info(v ...any) {
	if nil != infoLog {
		infoLog.Println(v...)
	}
}

func Infof(format string, v ...any) {
	if nil != infoLog {
		infoLog.Printf(format, v...)
	}
}

func IsDebug() bool {
	return debugLog != nil
}

func Debug(v ...any) {
	if nil != debugLog {
		debugLog.Println(v...)
	}
}

func Debugf(format string, v ...any) {
	if nil != debugLog {
		debugLog.Printf(format, v...)
	}
}

func Error(v ...any) {
	errorLog.Println(v...)
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

func Alert(v ...any) {
	message := fmt.Sprintln(v...)
	time := time.Now()
	header := fmt.Sprintf("Alert triggered in concert-manager!\nTime: %v\n", time)
	detail := fmt.Sprintf("Message: %s", message)
	body := fmt.Sprintf("%s\n%s", header, detail)
	if err := alerter.Alert(body); err != nil {
		Errorf("Failed to send alert, message: %s", body)
	}
	alertLog.Println(message)
}

func Alertf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	time := time.Now()
	header := fmt.Sprintf("Alert triggered in concert-manager!\nTime: %v\n", time)
	detail := fmt.Sprintf("Message: %s", message)
	body := fmt.Sprintf("%s\n%s", header, detail)
	if err := alerter.Alert(body); err != nil {
		Errorf("Failed to send alert, message: %s", body)
	}
	alertLog.Println(message)
}
