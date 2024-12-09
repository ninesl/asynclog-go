# go-debug-logger

A lightweight, concurrent logging package for Go applications with support for debugging line information.

NOTE: Because of the overhead of the package, `fmt.Println()` on its own IS faster than using `gologger`.

However, in applications with various goroutines and asynchronous operations, this package shines as `fmt.Println()` can be blocking at scale.

I personally use this package to easily manage debugging my codebases with lots of concurrency, think webscrapers, etc.

## Configuration
- Number of workers (default: 15, configurable via `SetWorkers(w int)`)
- Output destination (default: os.Stdout, configurable via `SetOutput(w io.Writer)`)
- Buffer size (default: 100 messages, configurable via `SetBuffer(b int)`) 

In heavier logging workloads, increasing the worker count or message buffer size can be more performant.

## Features
- Non-blocking message queue
- Concurrent message processing
- Debug output with file/line information
- Custom output streams
- Help() quick logging function

## Planned
- more configuration
- improved batch consuming of messages
- Ability to turn on/off cache
- Cache for `Print()`
- `DebugArgs()`
- Improving performance of `Args` functions 

## Installation
```go
go get github.com/ninesl/go-debug-logger
```

## Quick Start
```go
package main

import (
    "github.com/ninesl/go-debug-logger"
)

func main() {
    // Start logger
    gologger.Start()
    defer gologger.Stop()

    // Basic logging
    gologger.Print("Application started")
    
    // Debug logging. Output is "file.go:line# Processing request"
    gologger.Debug("Processing request")

    // Easy concatenation. Tip: This is a slow function and should only be used to quickly
    gologger.PrintArgs("Processing request ", someVariable)

    // Prints "Here" to the console.
    gologger.Here()

    // Prints "file.go:line# Here" to the console.
    gologger.DebugHere()
}
```

## Concurrent Usage
The logger is designed for concurrent environments such as:
- web scraping
- data parsing
- http requests
- parallel tasks
```go
func worker(id int) {
    for i := 0; i < 1000; i++ {
        gologger.Print("Processing item " + strconv.Itoa(i))
        
        if err := fooBar(); err != nil {
            gologger.Debug("encountered an error during fooBar()")
        }
    }
}

func main() {
    gologger.Start()
    defer gologger.Stop()

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(id)
        }(i)
    }
    wg.Wait()
}
```

## Benchmarks

`make benchmark` or `make benchmark COUNT={num}` simulates a concurrent environment to test `fmt.Println(msg)` vs `gologger.Print(msg)` 

Results on my machine (AMD Ryzen 5 3600X 6-Core Processor) after running `make benchmark`:

```bash        
# work being 'done' is time.Sleep(time.Nanosecond) to keep operations consistent
# each print is "processed item X worker X"
fmt.Println(): baseline
fmt.Printf():  -0.08% slower
fmt.Fprintf(): +0.05% faster
Debug():       +4.60% faster
Print():       +2.90% faster
PrintArgs():   -4.82% slower
# each print is "Here"
fmt.Println(): baseline
fmt.Printf():  -3.26% slower
fmt.Fprintf(): -2.80% slower
Debug():       +6.89% faster
Print():       +4.12% faster
PrintArgs():   +1.23% faster
Here():        +5.59% faster
DebugHere():   +4.62% faster
```
- Concurrent workloads: Logger performs same to better due to worker pool design
- Line numbers: `Debug()` gives runtime filename and line number information with no slowdowns.

In conclusion, there can be slight improvement to the performance when using `Debug()` and `Print()`,
 but depending on the workload of your go routines this package may be slightly overkill.

`Debug()` is more performant than `Print()` as the information is cached during runtime. This is useful when working
 with lots of go routines that may need to print the same info over and over again.

However, being able to easily print to the console any message and the file/line number it came from with virtually 
 no cost is a big win when it comes to my style of debugging.

I believe that the overhead tradeoff is negligible from the information and speedup gained from `Debug()` and `Print()`.

## Quick Debugging with Here() and DebugHere()

The package provides two convenient debugging functions:

### `Here()`
`Here()` prints "Here" to the console without line information.

You can customize the message:
```go
gologger.SetHere("checkpoint") // Must be called before Start()
gologger.Here() // Output: checkpoint
```

### `DebugHere()`

Same as `Here()` but includes file and line information.

```go
package main
// we are in cmd/main.go
import (
    "github.com/ninesl/go-debug-logger"
)

func doSomeWork() {
    gologger.Start()
    defer gologger.Stop()
    
    gologger.DebugHere() // Output: main.go:11 Here
    foo()
    gologger.DebugHere() // Output: main.go:13 Here
    
    if err := bar(); err != nil {
        gologger.Debug(err.Error()) // Output: main.go:16 "error message"
    }
}
```

Like `Print()` and `Debug()` both functions are thread-safe and can be used in concurrent code with minimal overhead.

