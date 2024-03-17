package log

import (
	"log"
	"os"
)

var (
	debugLog *log.Logger
	infoLog *log.Logger
	errorLog *log.Logger
	displayLog *log.Logger
	logLevels = [...]string{"DEBUG", "INFO", "ERROR"}
)

const (
	logPath = "/.concert_manager/logs/"
	fileName = "concert_manager.log"
	logLevelEnv = "LOG_LEVEL"
)

type loggerLevel int

const (
	debugLevel loggerLevel = iota
	infoLevel
	errorLevel
	defaultLevel
)

func Initialize() error {
	logFile, err := createLogFile(fileName)
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
	if logLevel <= errorLevel {
		errorLog = log.New(logFile, "error  -  ", log.LstdFlags)
	}

	if unsetLogLevel {
		Info("Unexpected or missing value for LOG_LEVEL environment variable, defaulting to INFO level")
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
