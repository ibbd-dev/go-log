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
	logs map[string]*AsyncLogger

	// 需要保护的key值
	// 在使用的时候，可以允许调用接口设置一次
	// 设置完之后，如果外部使用了相同的key，则会报错
	protectKeys map[string]bool

	sync.RWMutex
}

var asyncLogQueue *tAsyncLogQueue

func init() {
	asyncLogQueue = &tAsyncLogQueue{
		logs:        make(map[string]*AsyncLogger),
		protectKeys: make(map[string]bool),
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

// AddProtectKey 增加一个保护的key
func AddProtectKey(key string) {
	asyncLogQueue.Lock()
	defer asyncLogQueue.Unlock()

	if _, ok := asyncLogQueue.protectKeys[key]; !ok {
		asyncLogQueue.protectKeys[key] = true
	}
}

// New 获取日志对象
// key 日志对象的唯一标识，如果需要设置为保护key，则需要先调用AddProtectKey方法
func New(out io.Writer, prefix string, flag string, key string) *AsyncLogger {
	if key[0] == '_' {
		panic("Error key start with '_' NOT Allowed!")
	}

	var ok bool
	asyncLogQueue.Lock()
	defer asyncLogQueue.Unlock()

	// 判断是否为受保护key
	if _, ok = asyncLogQueue.protectKeys[key]; ok {
		// 受保护key加上特殊的前缀
		key = "_" + key
	}

	if l, ok := asyncLogQueue.logs[key]; ok {
		return l
	}

	l := &AsyncLogger{
		Logger:      log.Logger{},
		probability: 1.1, // 默认全部写入
		duration:    defaultDuration,
	}
	l.Logger.SetOutput(out)
	l.Logger.SetPrefix(prefix)
	l.Logger.SetFlags(flag)
	asyncLogQueue.logs[key] = l
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
