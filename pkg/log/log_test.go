package log

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

var g_waitGroup sync.WaitGroup

/*
*测试Log函数
 */
func Test_Log(t *testing.T) {
	SetLogLevel(LEVEL_DEBUG)
	SetLogModel(MODEL_DEV)
	t1 := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	for n := 0; n <= 1; n++ {
		g_waitGroup.Add(1)
		go func() {
			UserInfoLog(10, "test")
			Debug("test")
			Info("test")
			Warn("test")
			Error("test")
			g_waitGroup.Done()
		}()
	}
	g_waitGroup.Wait()
	fmt.Println(time.Since(t1))
	select {}
}
