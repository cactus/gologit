// Package gogitlog implements a very simple wrapper around the
// Go "log" package, providing support for a toggle-able debug flag
// and a couple of functions that log or not based on that flag.
package gologit

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var Logger = New(false)

// A DebugLogger represents a logging object, that embeds log.Logger, and
// provides support for a toggle-able debug flag.
type DebugLogger struct {
	*log.Logger
	debug bool
	mx sync.Mutex
}

// New creates a new DebugLogger.
// The debug argument specifies whether debug should be set or not.
func New(debug bool) *DebugLogger {
	return &DebugLogger{log.New(os.Stderr, "", log.LstdFlags), debug, sync.Mutex{}}
}

// Toggles the debug state.
// If debug is true, sets it to false.
// If debug is false, sets it to true.
func (l *DebugLogger) Toggle() {
	l.mx.Lock()
	defer l.mx.Unlock()
	if l.debug == false {
		l.debug = true
	} else {
		l.debug = false
	}
}

func (l *DebugLogger) ToggleOnSignal(sig os.Signal) {
	debugSig := make(chan os.Signal, 1)
	// spawn goroutine to handle signal/toggle of debug logging
	go func() {
		for {
			<-debugSig
			l.Toggle()
			if l.State() {
				l.Printf("Debug logging enabled")
			} else {
				l.Printf("Debug logging disabled")
			}
		}
	}()
	// notify send to debug sign channel on signusr1
	signal.Notify(debugSig, sig)
}

func (l *DebugLogger) State() bool {
	return l.debug
}

// Set the debug state directly.
func (l *DebugLogger) Set(debug bool) {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.debug = debug
}

// Debugf calls log.Printf if debug is true.
// If debug is false, does nothing.
func (l *DebugLogger) Debugf(format string, v ...interface{}) {
	if l.debug == true {
		l.Printf(format, v...)
	}
}

// Debug calls log.Print if debug is true.
// If debug is false, does nothing.
func (l *DebugLogger) Debug(v ...interface{}) {
	if l.debug == true {
		l.Print(v...)
	}
}

// Debugln calls log.Println if debug is true.
// If debug is false, does nothing.
func (l *DebugLogger) Debugln(v ...interface{}) {
	if l.debug == true {
		l.Println(v...)
	}
}
