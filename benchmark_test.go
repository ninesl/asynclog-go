package asynclog_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	asynclog "github.com/ninesl/asynclog-go"
)

func BenchmarkFmtPrintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(io.Discard, "Processing item %d\n", i)
	}
}

func BenchmarkLoggerPrint(b *testing.B) {
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		asynclog.Print("Processing item " + strconv.Itoa(i))
	}
	b.StopTimer()
}

func BenchmarkLoggerDebug(b *testing.B) {
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		asynclog.Debug("Processing item " + strconv.Itoa(i))
	}
	b.StopTimer()
}

const (
	asynclogWorkers  = 15
	asynclogBuffer   = 500
	benchmarkWorkers = 50
)

// conncurrent logging benchmarks
func BenchmarkConcurrentFmtPrintln(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						fmt.Println("Processing item ", i, " worker ", workerID)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtPrintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						fmt.Printf("Processing item %d worker %d\n", i, workerID)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentGologPrintln(b *testing.B) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)
						go log.Println("Processing item", i, "worker", workerID)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtFprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						fmt.Fprintf(os.Stdout, "Processing item %d worker %d\n", i, workerID)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentDebug(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.Debug("Processing item " + strconv.Itoa(i) + " worker " + strconv.Itoa(workerID))
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentPrint(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.Print("Processing item " + strconv.Itoa(i) + " worker " + strconv.Itoa(workerID))
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentPrintArgs(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.PrintArgs("Processing item ", i, " worker ", workerID)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

// Single message benchmarks
func BenchmarkConcurrentFmtPrintlnSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						fmt.Println("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtPrintfSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						fmt.Printf("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
func BenchmarkConcurrentFmtFprintfSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						fmt.Fprintf(os.Stdout, "Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentDebugSingle(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.Debug("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentHere(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.Here()
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentDebugHere(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.Here()
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentGologPrintlnSingle(b *testing.B) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)
						go log.Println("Processing item", i, "worker", workerID)
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentPrintSingle(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.Print("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}

func BenchmarkConcurrentPrintArgsSingle(b *testing.B) {
	asynclog.SetBuffer(asynclogBuffer)
	asynclog.SetWorkers(asynclogWorkers)
	asynclog.Start()
	defer asynclog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for w := 0; w < benchmarkWorkers; w++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				// Simulate CPU work
				matrix := make([][]struct{}, i)
				for x := range matrix {
					for range matrix[x] {
						time.Sleep(time.Nanosecond)

						asynclog.PrintArgs("Here")
					}
				}
			}(w)
		}
		wg.Wait()
	}
}
