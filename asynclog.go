// asynclog provides a simple logging mechanism with support for
// concurrent message processing, custom output destinations, and debugging
// information in the format of "pkg/filename.go:line"
//
// The package allows setting a buffer size for the message channel and
// redirecting log output to a specified io.Writer. It also supports starting
// the logger with a custom number of worker goroutines for message consumption.
//
// Usage:
//
//	// Configuration (optional):
//	   asynclog.SetBuffer(b int)     // 100 by default
//	   asynclog.SetWorkers(w int)    // 15 by default
//	   asynclog.SetOutput(io.Writer) // os.Stdout by default
//
//	// Start the logger:
//	   asynclog.Start()
//
//	// Send messages to your logger:
//	   asynclog.Print(msg string)
//	   asynclog.Debug(msg string) // includes file and line number
//
//	// Stop the logger to ensure all messages are consumed before the program exits:
//	   asynclog.Stop() // defer after Start()
//
// Example:
//
//	// Setting up the logger
//
//		asynclog.SetOutput(w io.Writer)
//		asynclog.SetBuffer(b int)
//		asynclog.SetWorkers(w int)
//		func SetWorkers(w int)
//		asynclog.Start()
//		defer asynclog.Stop()
//
//	//some work while calling these thread safe logging functions:
//		asynclog.Print(s string)
//		asynclog.Debug(s string)
//
// The logger is safe for concurrent use. In fact, that is why you would want to use this package.
//
// Configuring the logger can help you avoid I/O overhead.
//
// There is no significant speed up in benchmarks when compared to regular fmt.Printf() calls
// EXCEPT in cases where the logger is used in large numbers of concurrent goroutines.
package asynclog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	buffer     = 100
	messages   chan string
	workers              = 15
	isStarted            = false
	output     io.Writer = os.Stdout // Change type to io.Writer
	debugCache sync.Map
)

// DebugInfo represents debugging information that includes the file name, line number, and a string message.
// This struct is used to store and convey detailed debugging information within the logging system.
type DebugInfo struct {
	pc   uintptr
	file string
	line int
	str  string
}

func (info *DebugInfo) String() string {
	if info.str == "" {
		info.str = fmt.Sprintf("%s:%d", info.file, info.line)
	}
	return info.str
}

// Sets the buffer limit to the messages channel. Default is 100.
//
// Must be called before
//
//	Start()
//
// If the logger is already started, this function does nothing.
func SetBuffer(b int) {
	if isStarted {
		return
	}
	buffer = b
}

// Takes an io.Writer to redirect logs to a file or other destination. Default is
//
//	os.Stdout //the console
//
// Has to be called before
//
//	Start()
//
// If the logger is already started, this function does nothing.
func SetOutput(w io.Writer) {
	if isStarted {
		return
	}
	output = w
}

// Sets the number of worker goroutines for message consumption.
//
// Default is 15.
//
// Has to be called before
//
//	Start()
//
// If the logger is already started, this function does nothing.
func SetWorkers(w int) {
	if isStarted {
		return
	}
	workers = w
}

// Returns the file and line number of the caller.
//
// Uses the debugCache to avoid recomputing the same info.
func debugInfo() *DebugInfo {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return nil
	}

	if cached, ok := debugCache.Load(pc); ok {
		return cached.(*DebugInfo)
	}

	// Cache miss - compute and store
	_, file = filepath.Split(file)
	info := &DebugInfo{
		pc:   pc,
		file: file,
		line: line,
	}
	debugCache.Store(pc, info)
	return info
}

// Start initializes the logger by setting up the message channel, debug cache, and worker goroutines for concurrent message processing.
//
// If the logger is already started, it returns immediately. This function must be called before sending any messages to the logger.
//
// Example:
//
//	asynclog.SetOutput(w io.Writer)
//	asynclog.SetBuffer(b int)
//	func SetWorkers(w int)
//	asynclog.Start()
//	defer asynclog.Stop()
//
//	//some work while calling these thread safe logging functions:
//	asynclog.Print(s string)
//	asynclog.Debug(s string)
func Start() {
	if isStarted {
		return
	}
	messages = make(chan string, buffer)
	debugCache = sync.Map{}
	isStarted = true
	for i := 0; i < workers; i++ {
		go consumeMessages()
	}
}

// Closes the messages channel for a graceful shutdown
//
// Does nothing if Start() was not called
func Stop() {
	if !isStarted {
		return
	}
	isStarted = false
	close(messages)
}

// Convert any type to string efficiently
func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case error:
		return val.Error()
	default:
		return fmt.Sprint(val)
	}
}

// Print sends a string to the messages channel if the logger is started.
func Print(msg string) {
	if !isStarted {
		return
	}
	messages <- msg
}

// PrintArgs takes and sends a string to the messages channel if the logger is started.
//
// It takes a variable number of string arguments, concatenates them, and sends the result.
//
// If the logger is not started, the function returns immediately without sending any message.
//
// Parameters:
//
//	args ...any //A variadic parameter of strings to be concatenated and sent to the messages channel.
//
// This function is not faster than Print(). It is an ease of use function for concatenating multiple strings to the logger.
func PrintArgs(args ...any) {
	if !isStarted {
		return
	}

	if len(args) == 1 {
		messages <- toString(args[0])
		return
	}

	sargs := make([]string, len(args))
	for i, arg := range args {
		sargs[i] = toString(arg)
	}

	totalLen := 0
	for _, s := range sargs {
		totalLen += len(s)
	}

	var sb strings.Builder
	sb.Grow(totalLen)
	for _, arg := range args {
		sb.WriteString(toString(arg))
	}

	messages <- sb.String()
}

// Sends a string to the logger prepended with the file and line number of the caller.
//
// If the logger is not started, the message is ignored.
//
// FIXME: add % increase from the benchmarks
// Tip: fmt.Sprintf() _% slower than a basic string concatenation.
//
// Thread safe!
func Debug(msg string) {
	if !isStarted {
		return
	}
	info := debugInfo()

	if info != nil {
		msg = info.String() + " " + msg
	} else {
		msg = "ISSUE DETERMINING RUNTIME CALLER: " + msg
	}
	messages <- msg
}

var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// TODO: improvements
func consumeMessages() {
	const (
		batchSize     = 256       // Larger batches for better throughput
		bufferSize    = 1024 * 64 // 64KB buffer
		flushInterval = 500 * time.Millisecond
	)

	buf := make([]byte, 0, bufferSize)
	w := bufio.NewWriterSize(output, bufferSize)
	defer w.Flush()

	timer := time.NewTimer(flushInterval)
	defer timer.Stop()

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				if len(buf) > 0 {
					w.Write(buf)
					w.Flush()
				}
				return
			}

			buf = append(buf, msg...)
			buf = append(buf, '\n')

			if len(buf) >= batchSize {
				w.Write(buf)
				w.Flush()
				buf = buf[:0]
				timer.Reset(flushInterval)
			}

		case <-timer.C:
			if len(buf) > 0 {
				w.Write(buf)
				w.Flush()
				buf = buf[:0]
			}
			timer.Reset(flushInterval)
		}
	}
}

// SetHere sets the string message to be used by the Here() function.
//
// If the logger is already started, this function does nothing.
func SetHere(msg string) {
	if isStarted {
		return
	}
	here = msg
}

var (
	here = "Here"
)

// Here() sends the default "Here" message to the messages channel if the logger is started.
func Here() {
	if !isStarted {
		return
	}
	messages <- here
}

// DebugHere() is a convenience function that calls Debug() with whatever is set to SetHere() default "Here".
func DebugHere() {
	if !isStarted {
		return
	}
	Debug(here)
}
