package asyncLog

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ibbd-dev/go-log"
)

// log同步的状态
type syncStatus int

// 同步状态
const (
	statusInit  syncStatus = iota // 初始状态
	statusDoing                   // 同步中
	statusDone                    // 同步已经完成
)

const (
	// 换行符
	newlineStr  = "\n"
	newlineChar = '\n'

	// 缓存切片的初始容量
	cacheInitCap = 64
)

type AsyncLogger struct {
	log.Logger

	// 日志写入概率
	probability float32

	// 缓存配置
	cacheMu   sync.Mutex // 写cache时的互斥锁
	useCache  bool       // 是否使用缓存
	cacheData []string   // 缓存数据

	// 同步配置
	syncMu    sync.Mutex    // 保护下面两个属性
	duration  time.Duration // 同步数据到文件的周期，默认为1秒
	beginTime time.Time     // 开始同步的时间，判断同步的耗时
	status    syncStatus    // 同步状态
}

type tAsyncLogQueue struct {
	// 下标是文件名
	logs map[string]*AsyncLogger

	sync.RWMutex
}

var asyncLogQueue *tAsyncLogQueue

func init() {
	asyncLogQueue = &tAsyncLogQueue{
		logs: make(map[string]*AsyncLogger),
	}

	timer := time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			select {
			case <-timer.C:
				asyncLogQueue.RLock()
				for _, log := range asyncLogQueue.logs {
					if log.status != statusDoing {
						go log.flush()
					}
				}
				asyncLogQueue.RUnlock()
			}
		}

	}()
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
func NewAsyncLogger(out io.Writer, prefix string, flag string, filename string) *AsyncLogger {
	asyncLogQueue.Lock()
	defer asyncLogQueue.Unlock()

	if lf, ok := asyncLogQueue.logs[filename]; ok {
		return lf
	}

	lf := &AsyncLogger{
		Logger:      log.Logger{},
		probability: 1.1,
		useCache:    true,
		duration:    time.Second,
	}
	lf.Logger.SetOutput(out)
	lf.Logger.SetPrefix(prefix)
	lf.Logger.SetFlags(flag)
	asyncLogQueue.logs[filename] = lf
	return lf
}

func (l *AsyncLogger) Output(s string) error {
	now := time.Now() // get this early.
	l.syncMu.Lock()
	l.syncMu.Unlock()

	return l.Logger.Output(s)
}

func (l *AsyncLogger) SetUseCache(useCache bool) {
	l.useCache = useCache
}

func (l *AsyncLogger) SetProbability(probability float32) {
	l.probability = probability
}

func (l *AsyncLogger) SetDuration(duration time.Duration) {
	l.syncMu.Lock()
	defer l.syncMu.Unlock()
	l.duration = duration
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *AsyncLogger) Printf(format string, v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Output(fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *AsyncLogger) Print(v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Output(fmt.Sprint(v...))
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *AsyncLogger) Println(v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Output(fmt.Sprintln(v...))
}

func (l *AsyncLogger) PrintlnJson(data interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}

	s := json.Marshal(data) + newlineStr
	l.Output(s)
}

//***********************************

func (l *AsyncLogger) appendCache(s string) {
	l.cacheMu.Lock()
	l.cacheData = append(l.cacheData, s)
	l.cacheMu.Unlock()
}

func (l *AsyncLogger) flush() error {
	l.status = statusDoing
	defer func() {
		l.status = statusDone
	}()

	l.cacheMu.Lock()
	buffer := l.cacheData
	l.cacheData = make([]string, 0, cacheInitCap)
	l.cacheMu.Unlock()

	if len(buffer) == 0 {
		return nil
	}

	return l.Logger.Output(strings.Join(buffer, "") + newlineStr)
}
