package gologger_test

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"

	gologger "github.com/ninesl/go-debug-logger"
)

func BenchmarkFmtPrintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "Processing item %d\n", i)
	}
}

func BenchmarkLoggerPrint(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gologger.Print("Processing item " + strconv.Itoa(i))
	}
	b.StopTimer()
}

func BenchmarkLoggerDebug(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gologger.Debug("Processing item " + strconv.Itoa(i))
	}
	b.StopTimer()
}

const (
	gologgerWorkers  = 15
	gologgerBuffer   = 500
	benchmarkWorkers = 50
)

// conncurrent logging benchmarks
func BenchmarkConcurrentFmtPrintln(b *testing.B) {
	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						fmt.Println("Processing item ", i, " sum: ", sum)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtPrintf(b *testing.B) {
	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						fmt.Printf("Processing item %d sum %d\n", i, sum)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtFprintf(b *testing.B) {
	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						fmt.Fprintf(io.Discard, "Processing item %d sum %d\n", i, sum)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentDebug(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.Debug("Processing item " + strconv.Itoa(i) + " sum " + strconv.Itoa(sum))
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentPrint(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.Print("Processing item " + strconv.Itoa(i) + " sum " + strconv.Itoa(sum))
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentPrintArgs(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.PrintArgs("Processing item ", i, " sum: ", sum)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

// Single message benchmarks
func BenchmarkConcurrentFmtPrintlnSingle(b *testing.B) {
	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						fmt.Println("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtPrintfSingle(b *testing.B) {
	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						fmt.Printf("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtFprintfSingle(b *testing.B) {
	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						fmt.Fprintf(os.Stdout, "Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentDebugSingle(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.Debug("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentHere(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.Here()
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentDebugHere(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.Here()
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentPrintSingle(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.Print("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentPrintArgsSingle(b *testing.B) {
	// gologger.SetOutput(io.Discard)
	gologger.SetBuffer(gologgerBuffer)
	gologger.SetWorkers(gologgerWorkers)
	gologger.Start()
	defer gologger.Stop()

	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for j := range matrix[x] {
						sum += workerID * x
						sum *= j

						gologger.PrintArgs("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
