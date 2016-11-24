package asyncLog

import (
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
	// 默认的异步写入周期
	defaultDuration = time.Millisecond * 100
)

type AsyncLogger struct {
	log.Logger

	// 日志写入概率
	probability float32

	// 同步配置
	syncMu    sync.Mutex    // 保护下面两个属性
	duration  time.Duration // 同步数据到文件的周期，默认为1秒
	beginTime time.Time     // 开始同步的时间，判断同步的耗时
	status    syncStatus    // 同步状态
}

type tAsyncLogQueue struct {
	// 下标是文件名
	logs []*AsyncLogger

	sync.RWMutex
}

var asyncLogQueue *tAsyncLogQueue

func init() {
	asyncLogQueue = &tAsyncLogQueue{
		logs: make([]*AsyncLogger, 0, 8),
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

// New 获取日志对象
func New(out io.Writer, prefix string, flag string) *AsyncLogger {
	asyncLogQueue.Lock()
	defer asyncLogQueue.Unlock()

	l := &AsyncLogger{
		Logger:      log.Logger{},
		probability: 1.1, // 默认全部写入
		duration:    defaultDuration,
	}
	l.Logger.SetOutput(out)
	l.Logger.SetPrefix(prefix)
	l.Logger.SetFlags(flag)

	asyncLogQueue.logs = append(asyncLogQueue.logs, l)
	return l
}

func (l *AsyncLogger) Output(s string) error {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return nil
	}

	l.Cache(s)
	return nil
}

func (l *AsyncLogger) OutputBytes(s []byte) error {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return nil
	}

	l.CacheBytes(s)
	return nil
}

func (l *AsyncLogger) SetProbability(probability float32) {
	l.probability = probability
}

func (l *AsyncLogger) SetDuration(duration time.Duration) {
	l.syncMu.Lock()
	defer l.syncMu.Unlock()
	l.duration = duration
}

// Printf calls l.Cache to print to the logger.
func (l *AsyncLogger) Printf(format string, v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Cache(fmt.Sprintf(format, v...))
}

// Print calls l.Cache to print to the logger.
func (l *AsyncLogger) Print(v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Cache(fmt.Sprint(v...))
}

// Println calls l.Cache to print to the logger.
func (l *AsyncLogger) Println(v ...interface{}) {
	if l.probability < 1.0 && rand.Float32() > l.probability {
		return
	}
	l.Cache(fmt.Sprintln(v...))
}
