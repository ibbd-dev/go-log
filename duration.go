/*
按照时间周期写入日志，例如每秒写入1条
*/

package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type DurationLogger struct {
	Logger

	timeMu   sync.Mutex    // 保护下面两个属性
	duration time.Duration // 写log的周期，例如1秒
	lastTime time.Time     // 最后写log的时间
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties. 如time.RFC3339
func NewDurationLogger(out io.Writer, prefix string, flag string) *DurationLogger {
	return &DurationLogger{Logger: Logger{out: out, prefix: prefix, flag: flag}, duration: time.Second}
}

func (l *DurationLogger) Output(s string) error {
	if l.check() {
		return l.Logger.Output(s)
	}
	return nil
}

func (l *DurationLogger) SetDuration(duration time.Duration) {
	l.timeMu.Lock()
	defer l.timeMu.Unlock()
	l.duration = duration
}

// Printf calls l.Output to print to the logger.
func (l *DurationLogger) Printf(format string, v ...interface{}) {
	if l.check() {
		l.output(fmt.Sprintf(format, v...))
	}
}

// Print calls l.Output to print to the logger.
func (l *DurationLogger) Print(v ...interface{}) {
	if l.check() {
		l.output(fmt.Sprint(v...))
	}
}

// Println calls l.Output to print to the logger.
func (l *DurationLogger) Println(v ...interface{}) {
	if l.check() {
		l.output(fmt.Sprintln(v...))
	}
}

func (l *DurationLogger) output(s string) error {
	return l.Logger.Output(s)
}

func (l *DurationLogger) check() bool {
	now := time.Now()
	l.timeMu.Lock()
	if l.lastTime.Add(l.duration).After(now) {
		l.timeMu.Unlock()
		return false
	}
	l.lastTime = now
	l.timeMu.Unlock()

	return true
}
