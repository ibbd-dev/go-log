package errorLog

import (
	"os"
	"testing"
	"time"

	"github.com/ibbd-dev/go-log/async-log"
	"github.com/ibbd-dev/go-rotate-file"
)

func TestErrorLog1(t *testing.T) {
	// 文件Flag
	fileFlag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile("/tmp/test-error.log", fileFlag, 0666)
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	logger := asyncLog.New(file, "", time.RFC3339, "test-error")
	logger.SetDuration(time.Millisecond * 100)

	errorLog := New(logger, LevelWarn)
	errorLog.Debug(12, "Debug", 1.023)
	errorLog.Info(13, "Info", 1.023)
	errorLog.Warn(14, "Warn", 1.023)
	errorLog.Error(15, "Error", 1.023)
	errorLog.Fatal(16, "Fatal", 1.023)

	time.Sleep(time.Second)
}

func TestErrorLog2(t *testing.T) {
	file := rotateFile.Open("/tmp/test-error2.log")
	defer file.Close()

	logger := asyncLog.New(file, "", time.RFC3339, "test-error2")
	logger.SetDuration(time.Millisecond * 100)

	errorLog := New(logger, LevelWarn)
	errorLog.Debug(12, "Debug", 1.023)
	errorLog.Info(13, "Info", 1.023)
	errorLog.Warn(14, "Warn", 1.023)
	errorLog.Error(15, "Error", 1.023)
	errorLog.Fatal(16, "Fatal", 1.023)

	time.Sleep(time.Second)
}
