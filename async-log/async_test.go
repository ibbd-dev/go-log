package asyncLog

import (
	"os"
	"testing"
	"time"

	//"github.com/ibbd-dev/go-log"
	"github.com/ibbd-dev/go-rotate-file"
)

func TestLog(t *testing.T) {
	// 文件Flag
	fileFlag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile("/tmp/test-async.log", fileFlag, 0666)
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	logger := New(file, "", time.RFC3339, "test-async")
	logger.SetDuration(time.Millisecond * 100)
	logger.SetPrefix(time.RFC822)
	logger.Println("hello world")
	logger.Println("hello world2")
	time.Sleep(time.Millisecond * 105)
	logger.Println("hello world3")
	logger.Println("hello world3")
	time.Sleep(time.Millisecond * 10)
	logger.Println("hello world4")
}

func TestLog2(t *testing.T) {
	file := rotateFile.Open("/tmp/test-async2.log")
	defer file.Close()

	logger := New(file, "", time.RFC3339, "test-async2")
	logger.SetDuration(time.Millisecond * 100)
	logger.Println("hello world")
	logger.Println("hello world2")
	time.Sleep(time.Millisecond * 105)
	logger.Println("hello world3")
	logger.Println("hello world3")
	time.Sleep(time.Millisecond * 10)
	logger.Println("hello world4")
}
