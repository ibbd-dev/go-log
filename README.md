# golang logger

实现log的基本操作，实现按照时间周期写入

## Install 

```sh
# log基本操作，并实现按时间周期写入log
# 按时间周期写入，保证一个周期内，只会写入一次。
# 对于很多写log的情况，我们都需要控制一定的输出频率，避免log文件被写爆掉。
go get -u github.com/ibbd-dev/go-log

# 异步写log, 支持写入概率，例如以30%的概率写日志
# 注意：该项目的接口是以ibbd-dev/go-log为基础的
go get -u github.com/ibbd-dev/go-log/async-log

# 错误日志，支持错误等级
go get -u github.com/ibbd-dev/go-log/error-log
```

注：可以和`github.com/ibbd-dev/go-rotate-file`组合使用，该项目可以支持按时间切割文件，例如每小时一个文件

## Example


### 随机的周期性写入log

例如每秒写入一条log，能有效过于频繁的写log。

```go
package main

import (
	"os"
	"time"

	"github.com/ibbd-dev/go-log"
	"github.com/ibbd-dev/go-rotate-file"
)

func main() {
	// 文件Flag
	fileFlag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile("/tmp/test-rate.log", fileFlag, 0666)
	if err != nil {
		// TODO
	}
	defer file.Close()

	logger := log.NewDurationLogger(file, "", time.RFC3339)
	logger.SetDuration(time.Millisecond * 100)
	logger.SetPrefix("=====")
	logger.Println("hello world")
	logger.Println("hello world2")
	time.Sleep(time.Millisecond * 105)
	logger.Println("hello world3")
	logger.Println("hello world3")
	time.Sleep(time.Millisecond * 10)
	logger.Println("hello world4")
}

func main2() {
	file := rotateFile.Open("/tmp/test-rotate-log.log")
	defer file.Close()

	logger := NewDurationLogger(file, "", time.RFC3339)
	logger.SetDuration(time.Millisecond * 100)
	logger.Println("hello world")
	logger.Println("hello world2")
	time.Sleep(time.Millisecond * 105)
	logger.Println("hello world3")
	logger.Println("hello world3")
	time.Sleep(time.Millisecond * 10)
	logger.Println("hello world4")
}
```

### 异步写入日志

一些不太重要的日志，可以采取异步写入的形式，大大减轻io压力。

注意：批量写入log的过程中，如果程序被异常中断，可能会丢失部分数据，而丢失的多少，而设置的写入周期和日志的频率有关。

```go
package main

import (
	"time"

	"github.com/ibbd-dev/go-rotate-file"
	"github.com/ibbd-dev/go-log/async-log"
)

func main() {
	file := rotateFile.Open("/tmp/test-async.log")
	defer file.Close()

	logger := asyncLog.New(file, "", time.RFC3339, "test-async")
	logger.SetDuration(time.Millisecond * 100)
	logger.Println("hello world")
	logger.Println("hello world2")
	time.Sleep(time.Millisecond * 105)
	logger.Println("hello world3")
	logger.Println("hello world3")
	time.Sleep(time.Millisecond * 10)
	logger.Println("hello world4")
}
```

### 按概率写入日志

在`github.com/ibbd-dev/go-log/async-log`项目中有接口`SetProbability`，该接口可以用来设置写入概率，例如

```go
// 设置按10%的概率写日志
logger.SetProbability(0.1)
```

### 错误日志

对于提供`Debug`, `Info`, `Warn`, `Error`, `Fatal`等接口，方便错误日志的写入。

该项目可以方便的和其他日志对象的项目进行组合使用，例如和`asyncLog`组合使用可以实现按概率写入，或者批量异步写入，也可以和`ibbd-dev/go-log` 结合使用，或者和官方的`log`结合使用等。

```go
package main

import (
	"time"

	"github.com/ibbd-dev/go-rotate-file"
	"github.com/ibbd-dev/go-log/async-log"
	"github.com/ibbd-dev/go-log/error-log"
)

func main() {
	file := rotateFile.Open("/tmp/test-error.log")
	defer file.Close()

	logger := asyncLog.New(file, "", time.RFC3339, "test-error")
	logger.SetDuration(time.Millisecond * 100)

    // logger 写日志的对象，该对象只要有接口Output(string)error即可。
    // LevelWarn 写入级别，只有大于或者等于该级别的日志才会被写入
	errorLog := errorLog.New(logger, errorLog.LevelWarn)
	errorLog.Debug(12, "Debug", 1.023)
	errorLog.Info(13, "Info", 1.023)
	errorLog.Warn(14, "Warn", 1.023)
	errorLog.Error(15, "Error", 1.023)
	errorLog.Fatal(16, "Fatal", 1.023)

	time.Sleep(time.Second)
}

