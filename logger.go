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
		defaultCallerDepth: 3,
		logger:             logger,
		logPrefix:          "",
		currentLevel:       DEBUG,
		currentMode:        File,
	}
	return newLogger
}

// SetLogMode 設置Log模式
func (gLogger *GLogger) SetLogMode(mode OutputMode) {
	gLogger.currentMode = mode
	switch mode {
	case Stdout:
		gLogger.logger.SetOutput(os.Stdout)
	case File:
		gLogger.logger.SetOutput(gLogger.File)
	default:
		gLogger.logger.SetOutput(os.Stdout)
	}
}

// SetCurrentLevel 設置目前層級 DEBUG < INFO < WARNING < ERROR < FATAL
// EX : 設置 INFO，DEBUG就不會印出
func (gLogger *GLogger) SetCurrentLevel(level Level) {
	gLogger.currentLevel = level
}

// FatalF 致命層級
func (gLogger *GLogger) SetCallerDepth(depth int) {
	gLogger.defaultCallerDepth = depth
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

// Force 不管層級一定要印出
func (gLogger *GLogger) Force(v ...interface{}) {
	gLogger.writeLog(FORCE, v...)
}

// DebugF 除錯層級 printf
func (gLogger *GLogger) DebugF(format string, v ...interface{}) {
	if gLogger.currentLevel <= DEBUG {
		gLogger.writeFormatLog(DEBUG, format, v...)
	}
}

// InfoF 一般層級 printf
func (gLogger *GLogger) InfoF(format string, v ...interface{}) {
	if gLogger.currentLevel <= INFO {
		gLogger.writeFormatLog(INFO, format, v...)
	}
}

// WarnF 警告層級
func (gLogger *GLogger) WarnF(format string, v ...interface{}) {
	if gLogger.currentLevel <= WARNING {
		gLogger.writeFormatLog(WARNING, format, v...)
	}
}

// ErrorF 錯誤層級
func (gLogger *GLogger) ErrorF(format string, v ...interface{}) {
	if gLogger.currentLevel <= ERROR {
		gLogger.writeFormatLog(ERROR, format, v...)
	}
}

// FatalF 致命層級
func (gLogger *GLogger) FatalF(format string, v ...interface{}) {
	gLogger.writeFormatLog(FATAL, format, v...)
}

// ForceF 致命層級
func (gLogger *GLogger) ForceF(format string, v ...interface{}) {
	gLogger.writeFormatLog(FORCE, format, v...)
}

func (gLogger *GLogger) writeLog(level Level, v ...interface{}) {
	gLogger.write(level, "", v...)
}

func (gLogger *GLogger) writeFormatLog(level Level, format string, v ...interface{}) {
	gLogger.write(level, format, v...)
}

func (gLogger *GLogger) write(level Level, format string, v ...interface{}) {
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
		gLogger.logPrefix = fmt.Sprintf("[%s][%s:%d]", level.String(), t, line)
	} else {
		gLogger.logPrefix = fmt.Sprintf("[%s]", level.String())
	}

	gLogger.logger.SetPrefix(gLogger.logPrefix)

	if format == "" {
		if level != FATAL {
			gLogger.logger.Println(v...)
		} else {
			gLogger.logger.Fatalln(v...)
		}
	} else {
		if level != FATAL {
			gLogger.logger.Printf(format, v...)
		} else {
			gLogger.logger.Printf(format, v...)
		}
	}
}
