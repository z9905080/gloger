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

type GLogger struct {
	RWLock             *sync.RWMutex
	File               *os.File
	defaultPrefix      string
	defaultCallerDepth int
	logger             *log.Logger
	logPrefix          string
	levelFlags         []string
	currentLevel       Level
	currentMode        OutputMode
}

func NewLogger() *GLogger {
	filePath := getLogFileFullPath()
	file := openLogFile(filePath)
	logger := log.New(file, "", log.LstdFlags)

	newLogger := &GLogger{
		RWLock:             new(sync.RWMutex),
		File:               file,
		defaultPrefix:      "",
		defaultCallerDepth: 2,
		logger:             logger,
		logPrefix:          "",
		levelFlags:         []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		currentLevel:       DEBUG,
		currentMode:        File,
	}
	return newLogger
}

// SetLogMode 設置Log模式
func (gLogger *GLogger) SetLogMode(mode OutputMode) {
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
func (gLogger *GLogger) SetCurrentLevel(level Level) {
	gLogger.currentLevel = level
}

// Debug 除錯層級
func (gLogger *GLogger) Debug(v ...interface{}) {
	if gLogger.currentLevel <= DEBUG {
		gLogger.writeLog(DEBUG, v...)
	}
}

// Info 一般層級
func (gLogger *GLogger) Info(v ...interface{}) {
	if gLogger.currentLevel <= INFO {
		gLogger.writeLog(INFO, v...)
	}
}

// Warn 警告層級
func (gLogger *GLogger) Warn(v ...interface{}) {
	if gLogger.currentLevel <= WARNING {
		gLogger.writeLog(WARNING, v...)
	}
}

// Error 錯誤層級
func (gLogger *GLogger) Error(v ...interface{}) {
	if gLogger.currentLevel <= ERROR {
		gLogger.writeLog(ERROR, v...)
	}
}

// Fatal 致命層級
func (gLogger *GLogger) Fatal(v ...interface{}) {
	gLogger.writeLog(FATAL, v...)
}


func (gLogger *GLogger) writeLog(level Level, v ...interface{}) {

	// 加上Lock是為了防止同時寫入時重複開檔
	gLogger.RWLock.Lock()
	defer gLogger.RWLock.Unlock()
	NowDate := time.Now().Format(TimeFormat)
	if CurrentDate != NowDate && gLogger.currentMode == File {
		CurrentDate = NowDate
		gLogger.File.Close()
		filePath := getLogFileFullPath()
		gLogger.File = openLogFile(filePath)
		gLogger.logger = log.New(gLogger.File, gLogger.defaultPrefix, log.LstdFlags)
	}

	_, file, line, ok := runtime.Caller(gLogger.defaultCallerDepth)
	if ok {
		t, _ := filepath.Abs(file)
		gLogger.logPrefix = fmt.Sprintf("[%s][%s:%d]", gLogger.levelFlags[level], t, line)
	} else {
		gLogger.logPrefix = fmt.Sprintf("[%s]", gLogger.levelFlags[level])
	}

	gLogger.logger.SetPrefix(gLogger.logPrefix)
	if level != FATAL {
		gLogger.logger.Println(v...)
	} else {
		gLogger.logger.Fatalln(v...)
	}
}

