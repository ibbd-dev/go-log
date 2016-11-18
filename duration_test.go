package log

import (
	"os"
	"testing"
	"time"

	"github.com/ibbd-dev/go-rotate-file"
)

func TestLog(t *testing.T) {
	// 文件Flag
	fileFlag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile("/tmp/test-rate.log", fileFlag, 0666)
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	logger := NewDurationLogger(file, "", time.RFC3339)
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

func TestRotateLog(t *testing.T) {
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
