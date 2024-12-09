# asynclog

NOTE: `go log.Println()` with flags set is generally more
performant than this package. The worker pool design for consuming
messages could be helpful, but any meaningful performance
improvements would be negligible.

A lightweight, concurrent logging package for Go applications with support for debugging line information.

NOTE: Because of the overhead of the package, `fmt.Println()` on its own IS faster than using `asynclog` in most single-threaded Go programs.

However, in applications with various goroutines and asynchronous operations, this package shines as `fmt.Println()` can be blocking at scale.

I personally use this package to easily manage debugging my codebases with lots of concurrency, think webscrapers, etc.

## Configuration
- Number of workers, default: 15, configurable via `SetWorkers(w int)`
- Output destination, default: os.Stdout, configurable via `SetOutput(w io.Writer)`
- Buffer size, default: 100 messages, configurable via `SetBuffer(b int)`

In heavier logging workloads, increasing the worker count or message buffer size can be more performant.

## Features
- Non-blocking message queue
- Concurrent message processing
- Debug output with file/line information
- Custom output streams
- `Help()` quick logging function

## Planned
- more configuration
- improved batch consuming of messages
- Ability to turn on/off cache
- Cache for `Print()`
- `DebugArgs()`
- Improving performance of `Args` functions 
- Timestamps

## Installation
```go
go get github.com/ninesl/asynclog-go
```

## Quick Start
```go
package main

import (
    "github.com/ninesl/asynclog-go"
)

func main() {
    // Start logger
    asynclog.Start()
    defer asynclog.Stop()

    // Basic logging
    asynclog.Print("Application started")
    
    // Debug logging. Output is "file.go:line# Processing request"
    asynclog.Debug("Processing request")

    // Easy concatenation. Tip: This is a slow function and should only be used to quickly
    asynclog.PrintArgs("Processing request ", someVariable)

    // Prints "Here" to the console.
    asynclog.Here()

    // Prints "file.go:line# Here" to the console.
    asynclog.DebugHere()
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
        asynclog.Print("Processing item " + strconv.Itoa(i))
        
        if err := fooBar(); err != nil {
            asynclog.Debug("encountered an error during fooBar()")
        }
    }
}

func main() {
    asynclog.Start()
    defer asynclog.Stop()

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


You can run benchmarks using the following commands:
- `make benchmark` 
- `make benchmark COUNT={num}` (to specify the number of iterations)

These benchmarks simulate a concurrent environment to test `fmt.Println(msg)` vs `asynclog.Print(msg)`.


```go       
// work being 'done' is time.Sleep(time.Nanosecond) to keep operations consistent
// each print is "processed item X worker X"
fmt.Println():          baseline
fmt.Printf():           -0.08% slower
fmt.Fprintf():          +0.05% faster
asynclog.Debug():       +4.60% faster
asynclog.Print():       +2.90% faster
asynclog.PrintArgs():   -4.82% slower
// each print is "Here"
fmt.Println():          baseline
fmt.Printf():           -3.26% slower
fmt.Fprintf():          -2.80% slower
asynclog.Debug():       +6.89% faster
asynclog.Print():       +4.12% faster
go log.Println():       +9.74% faster // log.SetFlags(log.LstdFlags | log.Lshortfile)
go log.Print():         +3.84% faster // log.SetFlags(log.LstdFlags | log.Lshortfile)
go log.Print():         // log.SetFlags(log.Lshortfile)
asynclog.PrintArgs():   +1.23% faster
asynclog.Here():        +5.59% faster
asynclog.DebugHere():   +4.62% faster
```
- Concurrent workloads: Logger performs same to better due to worker pool design
- `Debug()` gives runtime filename and line number information with no slowdowns.

In conclusion, there can be slight improvement to the performance when using `Debug()` and `Print()`,
 but depending on the workload of your go routines this package may be slightly overkill.

`Debug()` is more performant than `Print()` as the information is cached during runtime. This is useful when working
 with lots of go routines that may need to print the same info over and over again.

However, being able to easily print to the console any message and the file/line number it came from with virtually 
 no cost is a big win when it comes to my style of debugging.

I believe that the tradeoff in overhead is negligible from the information and speedup gained from `Debug()` and `Print()`.

## Quick Debugging with Here() and DebugHere()

The package provides two convenient debugging functions:

### `Here()`
`Here()` prints "Here" to the console without line information.

You can customize the message:
```go
asynclog.SetHere("checkpoint") // Must be called before Start()
asynclog.Here() // Output: checkpoint
```

### `DebugHere()`

Same as `Here()` but includes file and line information.

```go
package main
// we are in cmd/main.go
import (
    "github.com/ninesl/asynclog-go"
)

func doSomeWork() {
    asynclog.Start()
    defer asynclog.Stop()
    
    asynclog.DebugHere() // Output: main.go:11 Here
    foo()
    asynclog.DebugHere() // Output: main.go:13 Here
    
    if err := bar(); err != nil {
        asynclog.Debug(err.Error()) // Output: main.go:16 "error message"
    }
}
```

Like `Print()` and `Debug()`, both functions are thread-safe and can be used in concurrent code with minimal overhead.

