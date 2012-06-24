// Package gogitlog implements a very simple wrapper around the
// Go "log" package, providing support for a toggle-able debug flag
// and a couple of functions that log or not based on that flag.
package gologit

import (
	"log"
	"os"
)

// A DebugLogger represents a logging object, that embeds log.Logger, and
// provides support for a toggle-able debug flag.
type DebugLogger struct {
	*log.Logger
	debug bool
}

// New creates a new DebugLogger.
// The debug argument specifies whether debug should be set or not.
func New(debug bool) *DebugLogger {
	return &DebugLogger{log.New(os.Stderr, "", log.LstdFlags), debug}
}

// Toggles the debug state.
// If debug is true, sets it to false.
// If debug is false, sets it to true.
func (l *DebugLogger) Toggle() {
	if l.debug == false {
		l.debug = true
	} else {
		l.debug = false
	}
}

// Set the debug state directly.
func (l *DebugLogger) Set(debug bool) {
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
