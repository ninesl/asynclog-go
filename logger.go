package gologger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	buffer    = 100
	messages  chan string
	workers             = 15
	isStarted           = false
	output    io.Writer = os.Stdout // Change type to io.Writer
	errors              = make(chan error)
)

// Takes an io.Writer to redirect logs to a file or other destination.
// The default is os.Stdout (the console).
func SetOutput(w io.Writer) {
	output = w
}

// Must be called before sending any messages.
// It is recommended to defer Stop() after Start() in main()
// to ensure all messages are consumed before the program exits.
// If Start() is not called, messages are ignored.
func Start() {
	if isStarted {
		return
	}
	messages = make(chan string, buffer)
	isStarted = true
	for i := 0; i < workers; i++ {
		go consumeMessages()
	}
}

// Closes the messages channel for a graceful shutdown
// Does nothing if Start() was not called
func Stop() {
	if !isStarted {
		return
	}
	isStarted = false
	close(messages)
}

// Sends a plain string to the logger.
// If the logger is not started, the message is ignored.
func Print(msg string) {
	if !isStarted {
		return
	}
	messages <- msg
}

// Sends a string to the logger prepended with the file and line number.
// If the logger is not started, the message is ignored.
func Debug(msg string) {
	if !isStarted {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fParts := strings.Split(file, "/")
		fileStr := fParts[len(fParts)-2] + "/" + fParts[len(fParts)-1]
		msg = fmt.Sprintf("%s:%d %s", fileStr, line, msg)
	} else {
		msg = "ISSUE DETERMINING RUNTIME CALLER: " + msg
	}
	messages <- msg
}

func consumeMessages() {
	for msg := range messages {
		log.Println(msg)
	}
}
