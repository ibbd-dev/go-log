package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type RateLogger struct {
	Logger

	timeMu   sync.Mutex    // 保护下面两个属性
	duration time.Duration // 写log的周期，例如1秒
	lastTime time.Time     // 最后写log的时间
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties. 如time.RFC3339
/*
   ANSIC       = "Mon Jan _2 15:04:05 2006"
   UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
   RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
   RFC822      = "02 Jan 06 15:04 MST"
   RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
   RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
   RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
   RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
   RFC3339     = "2006-01-02T15:04:05Z07:00"
   RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
   Kitchen     = "3:04PM"
   // Handy time stamps.
   Stamp      = "Jan _2 15:04:05"
   StampMilli = "Jan _2 15:04:05.000"
   StampMicro = "Jan _2 15:04:05.000000"
   StampNano  = "Jan _2 15:04:05.000000000"
*/
func NewRateLogger(out io.Writer, prefix string, flag string) *RateLogger {
	return &RateLogger{Logger: Logger{out: out, prefix: prefix, flag: flag}, duration: time.Second}
}

func (l *RateLogger) Output(s string) error {
	now := time.Now() // get this early.
	l.timeMu.Lock()
	if l.lastTime.Add(l.duration).After(now) {
		l.timeMu.Unlock()
		return nil
	}
	l.lastTime = now
	l.timeMu.Unlock()

	return l.Logger.Output(s)
}

func (l *RateLogger) SetDuration(duration time.Duration) {
	l.timeMu.Lock()
	defer l.timeMu.Unlock()
	l.duration = duration
}

// Printf calls l.Output to print to the logger.
func (l *RateLogger) Printf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
func (l *RateLogger) Print(v ...interface{}) { l.Output(fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
func (l *RateLogger) Println(v ...interface{}) { l.Output(fmt.Sprintln(v...)) }
