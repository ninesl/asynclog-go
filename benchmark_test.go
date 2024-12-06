package gologger

import (
	"fmt"
	"io"
	"sync"
	"testing"
)

func BenchmarkFmtPrintf(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "Processing item %d\n", i)
	}
}

func BenchmarkLoggerPrint(b *testing.B) {
	SetOutput(io.Discard)

	Start()
	defer Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Print(fmt.Sprintf("Processing item %d", i))
	}
}

func BenchmarkLoggerDebug(b *testing.B) {
	SetOutput(io.Discard)

	Start()
	defer Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Debug(fmt.Sprintf("Processing item %d", i))
	}
}

func BenchmarkConcurrentPrint(b *testing.B) {
	SetOutput(io.Discard)
	Start()
	defer Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numWorkers := 10

		// Launch workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]int, 50)
				for i := range matrix {
					matrix[i] = make([]int, 50)
					for j := range matrix[i] {
						matrix[i][j] = i * j
					}
				}

				// Mix logging with computation
				sum := 0
				for i := range matrix {
					for j := range matrix[i] {
						sum += matrix[i][j]
					}

					if i%10 == 0 {
						Print(fmt.Sprintf("Worker %d processed row %d: sum=%d",
							workerID, i, sum))
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentDebug(b *testing.B) {
	SetOutput(io.Discard)
	Start()
	defer Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numWorkers := 10

		// Launch workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]int, 50)
				for i := range matrix {
					matrix[i] = make([]int, 50)
					for j := range matrix[i] {
						matrix[i][j] = i * j
					}
				}

				// Mix logging with computation
				sum := 0
				for i := range matrix {
					for j := range matrix[i] {
						sum += matrix[i][j]
					}

					if i%10 == 0 {
						Debug(fmt.Sprintf("Worker %d processed row %d: sum=%d",
							workerID, i, sum))
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentFmtPrintf(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numWorkers := 10

		// Launch workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]int, 50)
				for i := range matrix {
					matrix[i] = make([]int, 50)
					for j := range matrix[i] {
						matrix[i][j] = i * j
					}
				}

				// Mix printf with computation
				sum := 0
				for i := range matrix {
					for j := range matrix[i] {
						sum += matrix[i][j]
					}

					if i%10 == 0 {
						fmt.Fprintf(io.Discard, "Worker %d processed row %d: sum=%d\n",
							workerID, i, sum)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
