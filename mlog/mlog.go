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
	"strings"
	"sync"
	"sync/atomic"
)

const (
	Lbase  uint64 = 0
	Ltime  uint64 = 1 << iota // log the date+time
	Ldebug                    // enable debug level log
	Lsort                     // sort keys in output
	Lstd   = Lbase | Ltime
)

var (
	bufPool     = newBufferPool()
	NEWLINE     = []byte("\n")
	SPACE       = []byte{' '}
	COLON       = []byte{':'}
	QUOTE       = []byte{'"'}
	EQUAL_QUOTE = []byte{'=', '"'}
	QUOTE_SPACE = []byte{'"', ' '}
)

type LogMap map[string]interface{}

func (lm *LogMap) Keys() []string {
	var keys []string
	for k := range *lm {
		keys = append(keys, k)
	}
	return keys
}

func (lm *LogMap) WriteTo(w io.Writer) (int64, error) {
	i := 0
	ilen := len(*lm)
	for k, v := range *lm {
		w.Write([]byte(k))
		w.Write(EQUAL_QUOTE)
		fmt.Fprint(w, v)
		w.Write(QUOTE)
		if i < ilen-1 {
			w.Write(SPACE)
		}
		i++
	}
	// int64 to be compat with io.WriterTo
	return int64(ilen), nil
}

func (lm *LogMap) SortedWriteTo(w io.Writer) (int64, error) {
	keys := lm.Keys()
	sort.Strings(keys)

	i := 0
	ilen := len(keys)
	for _, k := range keys {
		w.Write([]byte(k))
		w.Write(EQUAL_QUOTE)
		fmt.Fprint(w, (*lm)[k])
		w.Write(QUOTE)
		if i < ilen-1 {
			w.Write(SPACE)
		}
		i++
	}
	// int64 to be compat with WriterTo above
	return int64(ilen), nil
}

func (lm *LogMap) String() string {
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	lm.WriteTo(buf)
	return buf.String()
}

func (lm *LogMap) SortedString() string {
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	lm.SortedWriteTo(buf)
	return buf.String()
}

// A Logger represents a logging object, that embeds log.Logger, and
// provides support for a toggle-able debug flag.
type Logger struct {
	mu    sync.Mutex // ensures atomic writes are synchronized
	out   io.Writer
	flags uint64
}

func (l *Logger) Output(depth int, level string, message string, data ...interface{}) {
	// get this as soon as possible
	now := formattedDate.String()

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	flags := atomic.LoadUint64(&l.flags)
	if flags&Ltime != 0 {
		buf.Write([]byte(`time="`))
		buf.WriteString(now)
		buf.Write(QUOTE_SPACE)
	}

	buf.WriteString(`level="`)
	buf.WriteString(level)
	buf.Write(QUOTE)

	if flags&Ldebug != 0 {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			file = "???"
			line = 0
		}
		buf.WriteString(` caller="`)
		buf.WriteString(file)
		buf.Write(COLON)
		buf.WriteString(strconv.Itoa(line))
		buf.Write(QUOTE)
	}

	var mapv []*LogMap
	var fmtv []interface{}
	if len(data) > 0 {
		for _, v := range data {
			switch x := v.(type) {
			case *LogMap:
				if fmtv == nil {
					fmtv = make([]interface{}, 0)
				}
				mapv = append(mapv, x)
			case LogMap:
				if fmtv == nil {
					fmtv = make([]interface{}, 0)
				}
				mapv = append(mapv, &x)
			default:
				if mapv == nil {
					mapv = make([]*LogMap, 0)
				}
				fmtv = append(fmtv, v)
			}
		}
	}

	buf.WriteString(` msg="`)
	buf.WriteString(strings.TrimSpace(message))
	buf.Write(QUOTE)

	if len(mapv) > 0 {
		for _, e := range mapv {
			buf.Write(SPACE)
			if flags&Lsort != 0 {
				e.SortedWriteTo(buf)
			} else {
				e.WriteTo(buf)
			}
		}
	}

	lfmtv := len(fmtv)
	if lfmtv > 0 {
		for i, f := range fmtv {
			buf.WriteString(` extra`)
			buf.WriteString(strconv.Itoa(i + 1))
			buf.Write(EQUAL_QUOTE)
			fmt.Fprint(buf, f)
			buf.Write(QUOTE)
		}
	}

	buf.Write(NEWLINE)

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

// Debugf calls log.Print if debug is true.
// If debug is false, does nothing.
func (l *Logger) Debug(message string, v ...interface{}) {
	if l.HasDebug() {
		l.Output(2, "debug", message, v...)
	}
}

// Print calls log.Print
func (l *Logger) Info(message string, v ...interface{}) {
	l.Output(2, "info", message, v...)
}

// Fatalf calls log.Print then calls os.Exit(1)
func (l *Logger) Fatal(message string, v ...interface{}) {
	l.Output(2, "fatal", message, v...)
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

// Logs to the default Logger. See Logger.Debug
func Debug(message string, v ...interface{}) {
	if DefaultLogger.HasDebug() {
		DefaultLogger.Output(2, "[D]", message, v...)
	}
}

// Logs to the default Logger. See Logger.Print
func Info(message string, v ...interface{}) {
	DefaultLogger.Output(2, "[I]", message, v...)
}

// Logs to the default Logger. See Logger.Fatalf
func Fatalf(message string, v ...interface{}) {
	DefaultLogger.Output(2, "[F]", message, v...)
	os.Exit(1)
}
