package gloger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Level Debug的層級編號
type Level int

// 參數設定
var (
	RWLock             sync.RWMutex
	F                  *os.File
	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logger             *log.Logger
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	currentLevel       = DEBUG
	currentMode        = File
)

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
	F = openLogFile(filePath)
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func SetLogMode(mode OutputMode) {
	currentMode = mode
	switch mode {
	case Stdout:
		logger.SetOutput(os.Stdout)
	case File:
		logger.SetOutput(F)
	default:
		logger.SetOutput(os.Stdout)
	}

}

// SetCurrentLevel 設置目前層級 DEBUG < INFO < WARNING < ERROR < FATAL
// EX : 設置 INFO，DEBUG就不會印出
func SetCurrentLevel(level Level) {
	currentLevel = level
}

// Debug 除錯層級
func Debug(v ...interface{}) {
	if currentLevel <= DEBUG {
		writeLog(DEBUG, v...)
	}
}

// Info 一般層級
func Info(v ...interface{}) {
	if currentLevel <= INFO {
		writeLog(INFO, v...)
	}
}

// Warn 警告層級
func Warn(v ...interface{}) {
	if currentLevel <= WARNING {
		writeLog(WARNING, v...)
	}
}

// Error 錯誤層級
func Error(v ...interface{}) {
	if currentLevel <= ERROR {
		writeLog(ERROR, v...)
	}
}

// Fatal 致命層級
func Fatal(v ...interface{}) {
	writeLog(FATAL, v...)
}

func writeLog(level Level, v ...interface{}) {

	// 加上Lock是為了防止同時寫入時重複開檔
	RWLock.Lock()
	defer RWLock.Unlock()
	NowDate := time.Now().Format(TimeFormat)
	if CurrentDate != NowDate {
		CurrentDate = NowDate
		F.Close()
		filePath := getLogFileFullPath()
		F = openLogFile(filePath)
		if currentMode == File {
			logger = log.New(F, DefaultPrefix, log.LstdFlags)
		}
	}

	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		t, _ := filepath.Abs(file)
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], t, line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
	if level != FATAL {
		logger.Println(v...)
	} else {
		logger.Fatalln(v...)
	}
}