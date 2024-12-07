package gologger_test

// import (
// 	"fmt"
// 	"io"
// 	"math/rand"
// 	"sync"
// 	"testing"
// 	"time"

// 	gologger "github.com/ninesl/go-debug-logger"
// )

// func TestConcurrentLogging(t *testing.T) {
// 	gologger.SetOutput(io.Discard)
// 	gologger.Start()
// 	defer gologger.Stop()

// 	var wg sync.WaitGroup
// 	numRoutines := 10
// 	messagesPerRoutine := 100

// 	// Launch multiple goroutines to send messages
// 	for i := 0; i < numRoutines; i++ {
// 		wg.Add(1)
// 		go func(routineID int) {
// 			defer wg.Done()

// 			for j := 0; j < messagesPerRoutine; j++ {
// 				if j%2 == 0 {
// 					gologger.Print(fmt.Sprintf("Print from routine %d", routineID))
// 				} else {
// 					gologger.Debug(fmt.Sprintf("Debug from routine %d", routineID))
// 				}
// 				// Small sleep to simulate work
// 				time.Sleep(time.Millisecond)
// 			}
// 		}(i)
// 	}

// 	// Wait for all routines to complete
// 	wg.Wait()
// }
// func TestLoggerUnderLoad(t *testing.T) {
// 	gologger.SetOutput(io.Discard)
// 	gologger.Start()
// 	defer gologger.Stop()

// 	var wg sync.WaitGroup
// 	numRoutines := 20
// 	workItems := 50

// 	// Launch worker goroutines
// 	for i := 0; i < numRoutines; i++ {
// 		wg.Add(1)
// 		go func(routineID int) {
// 			defer wg.Done()

// 			// Simulate CPU-intensive work
// 			for j := 0; j < workItems; j++ {
// 				// Heavy computation
// 				result := 0
// 				for k := 0; k < 100000; k++ {
// 					result += k * k
// 				}

// 				if result%2 == 0 {
// 					gologger.Print(fmt.Sprintf("Routine %d completed heavy work item %d with result %d",
// 						routineID, j, result))
// 				} else {
// 					gologger.Debug(fmt.Sprintf("Routine %d heavy computation result: %d",
// 						routineID, result))
// 				}

// 				// Random sleep to simulate I/O
// 				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
// 			}
// 		}(i)
// 	}

// 	wg.Wait()
// }

// func TestLoggerUnderExtremeLoad(t *testing.T) {
// 	gologger.SetOutput(io.Discard)
// 	gologger.Start()
// 	defer gologger.Stop()

// 	var wg sync.WaitGroup
// 	numRoutines := 200
// 	// workItems := 1000
// 	matrixSize := 50

// 	// Launch massive number of goroutines
// 	for i := 0; i < numRoutines; i++ {
// 		wg.Add(1)
// 		go func(routineID int) {
// 			defer wg.Done()

// 			// Create large matrices for computation
// 			matrix1 := make([][]int, matrixSize)
// 			matrix2 := make([][]int, matrixSize)
// 			for i := range matrix1 {
// 				matrix1[i] = make([]int, matrixSize)
// 				matrix2[i] = make([]int, matrixSize)
// 				for j := range matrix1[i] {
// 					// Aggressive logging
// 					if j%3 == 0 {
// 						gologger.Print(fmt.Sprintf("Routine %d"))
// 					}
// 					gologger.Debug(fmt.Sprintf("Routine %d completed iteration %d", routineID, j))
// 					gologger.Print(fmt.Sprintf("Routine %d matrix sum: %d", routineID, result))

// 					if rand.Intn(100) < 10 {
// 						time.Sleep(time.Millisecond)
// 					}
// 			}
// 		}(i)
// 	}

// 	wg.Wait()
// }
