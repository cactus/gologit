// Copyright (c) 2012-2016 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package slog implements a very simple levelled logger
package slog

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

var timeFormat = "2006-01-02T15:04:05.000000"
var Logger = New(INFO, timeFormat, "", nil)
var bufpool = &sync.Pool{New: func() interface{} {
	return bytes.NewBuffer(make([]byte, 0, 128))
}}

func getBuffer() *bytes.Buffer {
	buf := bufpool.Get().(*bytes.Buffer)
	return buf
}

func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufpool.Put(buf)
}

type severity int32 // sync/atomic int32

const (
	DEBUG severity = iota
	INFO
	ERROR
	FATAL
)

var severityName = []string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

type LeveledLogger struct {
	timeformat string
	prefix     string
	severity   severity
	mx         sync.Mutex
	output     *os.File
}

func (l *LeveledLogger) header(s severity, t *time.Time) *bytes.Buffer {
	b := getBuffer()
	if l.timeformat != "" {
		fmt.Fprintf(b, "%s ", t.Format(l.timeformat))
	}
	fmt.Fprintf(b, "%-5.5s %s", severityName[s], l.prefix)
	return b
}

func (l *LeveledLogger) logln(s severity, v ...interface{}) {
	if s >= l.severity {
		t := time.Now()
		buf := l.header(s, &t)
		defer putBuffer(buf)
		fmt.Fprintln(buf, v...)
		buf.WriteTo(l.output)
	}
}

func (l *LeveledLogger) logf(s severity, format string, v ...interface{}) {
	if s >= l.severity {
		t := time.Now()
		buf := l.header(s, &t)
		defer putBuffer(buf)
		fmt.Fprintf(buf, format, v...)
		if buf.Bytes()[buf.Len()-1] != '\n' {
			buf.WriteByte('\n')
		}
		buf.WriteTo(l.output)
	}
}

func (l *LeveledLogger) Write(p []byte) (n int, err error) {
	t := time.Now()
	buf := l.header(INFO, &t)
	defer putBuffer(buf)
	buf.Write(p)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	written, err := buf.WriteTo(l.output)
	if err != nil {
		return int(written), err
	}
	return int(written), nil
}

func (l *LeveledLogger) SetOutput(output *os.File) {
	l.output = output
}

func (l *LeveledLogger) GetLevel() severity {
	return l.severity
}

func (l *LeveledLogger) SetLevel(s severity) {
	l.severity = s
}

func (l *LeveledLogger) IsDebug() bool {
	if l.severity == DEBUG {
		return true
	}
	return false
}

func (l *LeveledLogger) Debugf(format string, v ...interface{}) {
	l.logf(DEBUG, format, v...)
}

func (l *LeveledLogger) Debugln(v ...interface{}) {
	l.logln(DEBUG, v...)
}

func (l *LeveledLogger) Infof(format string, v ...interface{}) {
	l.logf(INFO, format, v...)
}

func (l *LeveledLogger) Infoln(v ...interface{}) {
	l.logln(INFO, v...)
}

func (l *LeveledLogger) Errorf(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
}

func (l *LeveledLogger) Errorln(v ...interface{}) {
	l.logln(ERROR, v...)
}

func (l *LeveledLogger) Fatalf(format string, v ...interface{}) {
	l.logf(FATAL, format, v...)
	os.Exit(1)
}

func (l *LeveledLogger) Fatalln(v ...interface{}) {
	l.logln(FATAL, v...)
	os.Exit(1)
}

func (l *LeveledLogger) Panicf(format string, v ...interface{}) {
	l.logf(FATAL, format, v...)
	panic(fmt.Sprintf(format, v...))
}

func (l *LeveledLogger) Panicln(v ...interface{}) {
	l.logln(FATAL, v...)
	panic(fmt.Sprintln(v...))
}

func New(level severity, timeformat string, prefix string, output *os.File) *LeveledLogger {
	if output == nil {
		output = os.Stderr
	}
	return &LeveledLogger{
		timeformat,
		prefix,
		level,
		sync.Mutex{},
		output,
	}
}

/*
// isatty returns true if f is a TTY, false otherwise.
func isatty(f *os.File) bool {
	switch runtime.GOOS {
	case "darwin":
	case "linux":
	default:
		return false
	}
	var t [2]byte
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		f.Fd(), syscall.TIOCGPGRP,
		uintptr(unsafe.Pointer(&t)))
	return errno == 0
}
*/

func GetLevel() severity {
	return Logger.GetLevel()
}

func SetLevel(s severity) {
	Logger.SetLevel(s)
}

func IsDebug() bool {
	return Logger.IsDebug()
}

func Debugf(format string, v ...interface{}) {
	Logger.Debugf(format, v...)
}

func Debugln(v ...interface{}) {
	Logger.Debugln(v...)
}

func Infof(format string, v ...interface{}) {
	Logger.Infof(format, v...)
}

func Infoln(v ...interface{}) {
	Logger.Infoln(v...)
}

func Errorf(format string, v ...interface{}) {
	Logger.Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	Logger.Errorln(v...)
}

func Fatalf(format string, v ...interface{}) {
	Logger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	Logger.Fatalln(v...)
}

func Panicf(format string, v ...interface{}) {
	Logger.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	Logger.Panicln(v...)
}
