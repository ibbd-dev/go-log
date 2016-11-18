package errorLog

import (
	"fmt"
	"sync"
)

// 日志优先级
type Priority int

const (
	LevelAll Priority = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
)

type ILogger interface {
	Output(string) error
}

type ErrorLogger struct {
	logger ILogger

	mu    sync.Mutex
	level Priority
}

var (
	// 日志等级
	levelTitle = map[Priority]string{
		LevelDebug: "[DEBUG]",
		LevelInfo:  "[INFO]",
		LevelWarn:  "[WARN]",
		LevelError: "[ERROR]",
		LevelFatal: "[FATAL]",
	}
)

// NewLevelLog 写入等级日志
// 级别高于或者等于level才会被写入
func New(logger ILogger, level Priority) *ErrorLogger {
	return &ErrorLogger{
		logger: logger,
		level:  level,
	}
}

func (l *ErrorLogger) SetLevel(level Priority) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *ErrorLogger) directOutput(level Priority, v ...interface{}) {
	if level >= l.level {
		l.logger.Output(levelTitle[level] + " " + fmt.Sprintln(v...))
	}
}

func (l *ErrorLogger) Debug(v ...interface{}) {
	l.directOutput(LevelDebug, v...)
}

func (l *ErrorLogger) Info(v ...interface{}) {
	l.directOutput(LevelInfo, v...)
}

func (l *ErrorLogger) Warn(v ...interface{}) {
	l.directOutput(LevelWarn, v...)
}

func (l *ErrorLogger) Error(v ...interface{}) {
	l.directOutput(LevelError, v...)
}

// 这个级别的错误都要写日志
func (l *ErrorLogger) Fatal(v ...interface{}) {
	l.directOutput(LevelFatal, v...)
}
