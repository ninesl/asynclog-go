# go-debug-logger

A lightweight, concurrent logging package for Go applications with support for debug line information.

## Features
- Non-blocking message queue
- Concurrent message processing
- Debug mode with file/line information

## Planned Features
- Custom output streams
- Other configurability (number of workers, etc.)

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
    
    // Debug logging includes file and line information
    gologger.Debug("Processing request")
}
```

## Concurrent Usage
The logger is designed for concurrent environments such as:
- web scraping
- data parsing
- http requests
```go
func worker(id int) {
    for i := 0; i < 1000; i++ {
        gologger.Print("Processing item " + strconv.Itoa(i))
        fooBar()
        if err != nil {
            gologger.Debug("")
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