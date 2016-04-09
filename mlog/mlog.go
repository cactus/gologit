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
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	bufPool = newBufferPool()
)

const (
	Lbase  uint64 = 0
	Ltime  uint64 = 1 << iota // log the date+time
	Ldebug                    // enable debug level log
	Lsort                     // sort keys in output
	Lstd   = Lbase | Ltime
)

type ExtraMap map[string]interface{}

// A Logger represents a logging object, that embeds log.Logger, and
// provides support for a toggle-able debug flag.
type Logger struct {
	mu    sync.Mutex // ensures atomic writes are synchronized
	out   io.Writer
	flags uint64
}

func (l *Logger) Output(depth int, level string, format string, v ...interface{}) {
	// get this as soon as possible
	now := formattedDate.String()

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	flags := atomic.LoadUint64(&l.flags)
	if flags&Ltime != 0 {
		buf.Write([]byte(`time="`))
		buf.WriteString(now)
		buf.WriteString(`" `)
	}

	buf.WriteString(`level="`)
	buf.WriteString(level)
	buf.WriteByte('"')

	if flags&Ldebug != 0 {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			file = "???"
			line = 0
		}
		buf.WriteString(` caller="`)
		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
		buf.WriteByte('"')
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

	buf.WriteString(` msg="`)
	if format != "" {
		fmt.Fprintf(buf, format, fmtv...)
	}
	buf.WriteByte('"')

	if len(mapv) > 0 {
		for _, e := range mapv {
			if flags&Lsort != 0 {
				var keys []string
				for k := range *e {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					buf.WriteByte(' ')
					buf.WriteString(k)
					buf.WriteString(`="`)
					fmt.Fprint(buf, (*e)[k])
					buf.WriteByte('"')
				}
			} else {
				for k, v := range *e {
					buf.WriteByte(' ')
					buf.WriteString(k)
					buf.WriteString(`="`)
					fmt.Fprint(buf, v)
					buf.WriteByte('"')
				}
			}
		}
	}

	buf.WriteByte('\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	buf.WriteTo(l.out)
}

func (l *Logger) SetFlags(flags uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	atomic.StoreUint64(&l.flags, flags)
}

func (l *Logger) HasDebug() bool {
	flags := atomic.LoadUint64(&l.flags)
	return flags&Ldebug != 0
}

// Debugf calls log.Printf if debug is true.
// If debug is false, does nothing.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.HasDebug() {
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
func New(out io.Writer, flags uint64) *Logger {
	return &Logger{
		out:   out,
		flags: flags,
	}
}

// default Logger
var DefaultLogger = New(os.Stderr, Lstd)

func SetFlags(flags uint64) {
	DefaultLogger.SetFlags(flags)
}

// Logs to the default Logger. See Logger.Debugf
func Debugf(format string, v ...interface{}) {
	if DefaultLogger.HasDebug() {
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
