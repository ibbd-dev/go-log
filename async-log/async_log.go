package asyncLog

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
//newlineStr  = "\n"
//newlineChar = '\n'

// 缓存切片的初始容量
//cacheInitCap = 64
)

type AsyncLogger struct {
	log.Logger

	// 日志写入概率
	probability float32

	// 缓存配置
	useCache bool // 是否使用缓存

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
						go log.Flush()
					}
				}
				asyncLogQueue.RUnlock()
			}
		}

	}()
}

func New(out io.Writer, prefix string, flag string, filename string) *AsyncLogger {
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
	if l.useCache {
		l.Cache(s)
		return nil
	}

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
func (l *AsyncLogger) Printf(format string, v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Output(fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
func (l *AsyncLogger) Print(v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Output(fmt.Sprint(v...))
}

// Println calls l.Output to print to the logger.
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

	bts, err := json.Marshal(data)
	if err != nil {
		// TODO
	}
	l.Output(string(bts))
}
