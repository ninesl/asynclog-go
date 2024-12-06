package gologger

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentLogging(t *testing.T) {
	Start()
	defer Stop()

	var wg sync.WaitGroup
	numRoutines := 10
	messagesPerRoutine := 100

	// Launch multiple goroutines to send messages
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()

			for j := 0; j < messagesPerRoutine; j++ {
				if j%2 == 0 {
					Print(fmt.Sprintf("Print from routine %d", routineID))
				} else {
					Debug(fmt.Sprintf("Debug from routine %d", routineID))
				}
				// Small sleep to simulate work
				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	// Wait for all routines to complete
	wg.Wait()

	// Verify channel is not blocked
	if len(messages) > 0 {
		t.Errorf("Messages channel not empty, %d messages remaining", len(messages))
	}
}
func TestLoggerUnderLoad(t *testing.T) {
	Start()
	defer Stop()

	var wg sync.WaitGroup
	numRoutines := 20
	workItems := 50

	// Launch worker goroutines
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()

			// Simulate CPU-intensive work
			for j := 0; j < workItems; j++ {
				// Heavy computation
				result := 0
				for k := 0; k < 100000; k++ {
					result += k * k
				}

				// Mix in some logging
				if result%2 == 0 {
					Print(fmt.Sprintf("Routine %d completed heavy work item %d with result %d",
						routineID, j, result))
				} else {
					Debug(fmt.Sprintf("Routine %d heavy computation result: %d",
						routineID, result))
				}

				// Random sleep to simulate I/O
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	// Verify logger handled the load
	if len(messages) > 0 {
		t.Errorf("Logger backed up: %d messages remaining", len(messages))
	}
}

func TestLoggerUnderExtremeLoad(t *testing.T) {
	Start()
	defer Stop()

	var wg sync.WaitGroup
	numRoutines := 200
	workItems := 1000
	matrixSize := 50

	// Create shared data structure to increase contention
	sharedCounter := atomic.Int64{}

	// Launch massive number of goroutines
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()

			// Create large matrices for computation
			matrix1 := make([][]int, matrixSize)
			matrix2 := make([][]int, matrixSize)
			for i := range matrix1 {
				matrix1[i] = make([]int, matrixSize)
				matrix2[i] = make([]int, matrixSize)
				for j := range matrix1[i] {
					matrix1[i][j] = rand.Intn(100)
					matrix2[i][j] = rand.Intn(100)
				}
			}

			for j := 0; j < workItems; j++ {
				// Heavy matrix multiplication
				result := 0
				for x := 0; x < matrixSize; x++ {
					for y := 0; y < matrixSize; y++ {
						for k := 0; k < matrixSize; k++ {
							result += matrix1[x][k] * matrix2[k][y]
						}
					}
				}

				// Increment shared counter
				count := sharedCounter.Add(1)

				// Aggressive logging
				if j%3 == 0 {
					Print(fmt.Sprintf("Routine %d matrix calc %d: %d (total ops: %d)",
						routineID, j, result, count))
				}
				Debug(fmt.Sprintf("Routine %d completed iteration %d", routineID, j))
				Print(fmt.Sprintf("Routine %d matrix sum: %d", routineID, result))

				// Small random pauses
				if rand.Intn(100) < 10 {
					time.Sleep(time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify logger handled extreme load
	if len(messages) > 0 {
		t.Errorf("Logger backed up under load: %d messages remaining", len(messages))
	}
}
