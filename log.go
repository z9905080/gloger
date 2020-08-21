package gloger

import (
	"log"
	"os"
	"sync"
)

// Level Debug的層級編號
type Level int

// 參數設定
//var (
//	RWLock             =new(sync.RWMutex)
//	File                  *os.File
//	DefaultPrefix      = ""
//	DefaultCallerDepth = 2
//	logger             *log.Logger
//	logPrefix          = ""
//	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
//	currentLevel       = DEBUG
//	currentMode        = File
//)

// 層級列表
const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type OutputMode int

// 層級列表
const (
	Stdout OutputMode = iota
	File
)

func init() {
	filePath := getLogFileFullPath()
	file := openLogFile(filePath)
	logger := log.New(file, "", log.LstdFlags)

	mLogger = GLogger{
		RWLock:             new(sync.RWMutex),
		File:               file,
		defaultPrefix:      "",
		defaultCallerDepth: 4,
		logger:             logger,
		logPrefix:          "",
		levelFlags:         []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		currentLevel:       DEBUG,
		currentMode:        File,
	}
}

var mLogger GLogger

func SetLogMode(mode OutputMode) {
	mLogger.currentMode = mode
	switch mode {
	case Stdout:
		mLogger.logger.SetOutput(os.Stdout)
	case File:
		mLogger.logger.SetOutput(mLogger.File)
	default:
		mLogger.logger.SetOutput(os.Stdout)
	}
}

// SetCurrentLevel 設置目前層級 DEBUG < INFO < WARNING < ERROR < FATAL
// EX : 設置 INFO，DEBUG就不會印出
func SetCurrentLevel(level Level) {
	mLogger.SetCurrentLevel(level)
}

// Debug 除錯層級
func Debug(v ...interface{}) {
	mLogger.Debug(v...)
}

// Info 一般層級
func Info(v ...interface{}) {
	mLogger.Info(v...)
}

// Warn 警告層級
func Warn(v ...interface{}) {
	mLogger.Warn(v...)
}

// Error 錯誤層級
func Error(v ...interface{}) {
	mLogger.Error(v...)
}

// Fatal 致命層級
func Fatal(v ...interface{}) {
	mLogger.Fatal(v...)
}

// DebugF 除錯層級
func DebugF(format string, v ...interface{}) {
	mLogger.DebugF(format, v...)
}

// InfoF 一般層級
func InfoF(format string, v ...interface{}) {
	mLogger.InfoF(format, v...)
}

// WarnF 警告層級
func WarnF(format string, v ...interface{}) {
	mLogger.WarnF(format, v...)
}

// ErrorF 錯誤層級
func ErrorF(format string, v ...interface{}) {
	mLogger.ErrorF(format, v...)
}

// Fatal 致命層級
func FatalF(format string, v ...interface{}) {
	mLogger.FatalF(format, v...)
}

// GetLogger 取得套件內的Logger
func GetLogger() *GLogger {
	return &mLogger
}
