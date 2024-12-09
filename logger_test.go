// logger_test.go
package gologger_test

import (
	"sync"
	"testing"
	"time"

	gologger "github.com/ninesl/go-debug-logger"
)

func TestLoggerWorkerPattern(t *testing.T) {
	gologger.SetBuffer(10)
	gologger.SetWorkers(5)
	gologger.Start()
	defer gologger.Stop()

	var wg sync.WaitGroup
	workerCount := 10

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Simulate work
			time.Sleep(time.Millisecond * 10)
			gologger.Print("finished work")
			gologger.PrintArgs("finished worker ", id)
			gologger.Debug("Debugging")
			gologger.Here()
			gologger.DebugHere()
		}(i)
	}

	wg.Wait()
}
