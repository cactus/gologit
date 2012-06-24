package gologit

import (
	"log"
	"os"
)

type DebugLogger struct {
	*log.Logger
	debug bool
}

func New(debug bool) *DebugLogger {
	return &DebugLogger{log.New(os.Stderr, "", log.LstdFlags), debug}
}

func (l *DebugLogger) Toggle() {
	if l.debug == false {
		l.debug = true
	} else {
		l.debug = false
	}
}

func (l *DebugLogger) Set(v bool) {
	l.debug = v
}

func (l *DebugLogger) Debugf(format string, v ...interface{}) {
	if l.debug == true {
		l.Printf(format, v...)
	}
}

func (l *DebugLogger) Debug(v ...interface{}) {
	if l.debug == true {
		l.Print(v...)
	}
}

func (l *DebugLogger) Debugln(v ...interface{}) {
	if l.debug == true {
		l.Println(v...)
	}
}
