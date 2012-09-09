package logging

import (
	"fmt"
	"log"
)

const (
	NOLOGGING = -1
	FATAL     = 0
	ERROR     = 1
	WARN      = 2
	INFO      = 3
	DEBUG     = 4
)

var (
	LogLevel  int = ERROR
	logger    *log.Logger
	logPrefix string = "doozerd "
	LogColor         = map[int]string{FATAL: "\033[0m\033[37m",
		ERROR: "\033[0m\033[31m",
		WARN:  "\033[0m\033[33m",
		INFO:  "\033[0m\033[35m",
		DEBUG: "\033[0m\033[34m"}
	LogLevelWords map[string]int = map[string]int{"fatal": 0, "error": 1, "warn": 2, "info": 3, "debug": 4, "none": -1}
)

// you can set a logger, and log level,most common usage is:
//
//		SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")
//
//  Or you can set the log out file name as well:
//  	
//		SetLogger(log.New(os.Stderr, "", log.Lmicroseconds | log.Lshortfile), "debug")
// 
//  loglevls:   debug, info, warn, error, fatal
func SetLogger(l *log.Logger, logLevel string) {
	logger = l
	LogLevelSet(logLevel)
}
func GetLogger() *log.Logger {
	return logger
}
func SetLoggerPrefix(prefix string) {
	logPrefix = prefix
}

// sets the log level from a string
func LogLevelSet(levelWord string) {
	if lvl, ok := LogLevelWords[levelWord]; ok {
		LogLevel = lvl
	}
}

func LogPf(logLvl int, prefix string, format string, v ...interface{}) {
	if LogLevel >= logLvl && logger != nil {
		if len(format) > 0 {
			logger.Output(3, prefix+LogColor[logLvl]+fmt.Sprintf(format, v...)+"\033[0m")
		} else {
			logger.Output(3, prefix+LogColor[logLvl]+fmt.Sprint(v...)+"\033[0m")
		}
	}
}

func Debug(v ...interface{}) {
	LogPf(DEBUG, logPrefix, "", v...)
}
func Debugf(format string, v ...interface{}) {
	LogPf(DEBUG, logPrefix, format, v...)
}
func Log(lvl int, v ...interface{}) {
	LogPf(lvl, logPrefix, "", v...)
}
func Logf(lvl int, format string, v ...interface{}) {
	LogPf(lvl, logPrefix, format, v...)
}
