// Copyright (c) 2012-2016 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package gologit implements a very simple wrapper around the
// Go "log" package, providing support for a toggle-able debug flag
// and a couple of functions that log or not based on that flag.
package mlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

var bufPool = newBufferPool()

type ExtraMap map[string]interface{}

// A Logger represents a logging object, that embeds log.Logger, and
// provides support for a toggle-able debug flag.
type Logger struct {
	mu          sync.Mutex // ensures atomic writes to out
	out         io.Writer
	enableDebug bool
	showTime    bool
}

func (l *Logger) Output(depth int, level string, format string, v ...interface{}) {
	// get this as soon as possible
	now := formattedDate.String()

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	if l.showTime {
		fmt.Fprintf(buf, `time="%s" `, now)
	}

	if level != "" {
		fmt.Fprintf(buf, `level="%s" `, level)
	}

	if l.enableDebug {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			file = "???"
			line = 0
		}
		fmt.Fprintf(buf, `caller="%s:%s" `, file, line)
	}

	var mapv []*ExtraMap
	var fmtv []interface{}
	if len(v) > 0 {
		mapv = make([]*ExtraMap, 0)
		fmtv = make([]interface{}, 0)
		for _, x := range v {
			if y, ok := x.(*ExtraMap); ok {
				mapv = append(mapv, y)
			} else if y, ok := x.(ExtraMap); ok {
				mapv = append(mapv, &y)
			} else {
				fmtv = append(fmtv, x)
			}
		}
	}

	if format == "" {
		fmt.Fprint(buf, `msg="" `)
	} else {
		fmt.Fprint(buf, `msg="`)
		fmt.Fprintf(buf, format, fmtv...)
		buf.Write([]byte(`" `))
	}

	if len(mapv) > 0 {
		for _, e := range mapv {
			for k, v := range *e {
				fmt.Fprintf(buf, `%s="`, k)
				fmt.Fprint(buf, v)
				buf.Write([]byte(`" `))
			}
		}
	}

	// trim off trailing space
	buf.Truncate(buf.Len() - 1)
	buf.WriteByte('\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	buf.WriteTo(l.out)
}

func (l *Logger) HasTimestamp() bool {
	return l.showTime
}

func (l *Logger) HasDebugging() bool {
	return l.enableDebug
}

// Debugf calls log.Printf if debug is true.
// If debug is false, does nothing.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.enableDebug == true {
		l.Output(2, "debug", format, v...)
	}
}

// Printf calls log.Printf
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, "info", format, v...)
}

// Fatalf calls log.Printf then calls os.Exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, "fatal", format, v...)
	os.Exit(1)
}

// New creates a new Logger.
// The debug argument specifies whether debug should be set or not.
func New(out io.Writer, enableDebug, showTime bool) *Logger {
	flags := log.LstdFlags
	if enableDebug == true {
		flags = flags | log.Lshortfile
	}
	return &Logger{out: out, enableDebug: enableDebug, showTime: showTime}
}

// default Logger
var DefaultLogger = New(os.Stderr, false, false)

// returns whether debugging is enabled for the defualt Logger.
func HasDebugging() bool {
	return DefaultLogger.HasDebugging()
}

// returns whether timestamp output is enabled by the defualt Logger.
func HasTimestamp() bool {
	return DefaultLogger.HasTimestamp()
}

// Logs to the default Logger. See Logger.Debugf
func Debugf(format string, v ...interface{}) {
	if DefaultLogger.HasDebugging() {
		DefaultLogger.Output(2, "[D]", format, v...)
	}
}

// Logs to the default Logger. See Logger.Printf
func Printf(format string, v ...interface{}) {
	DefaultLogger.Output(2, "[I]", format, v...)
}

// Logs to the default Logger. See Logger.Fatalf
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Output(2, "[F]", format, v...)
	os.Exit(1)
}
