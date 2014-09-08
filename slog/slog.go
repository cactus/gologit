// Package gogitlog implements a very simple wrapper around the
// Go "log" package, providing support for a toggle-able debug flag
// and a couple of functions that log or not based on that flag.
package slog

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

var timeFormat = "2006-01-02T15:04:05.000000"
var logger = New(INFO, "")

type severity int32 // sync/atomic int32

const (
	DEBUG severity = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var severityName = []string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

type LeveledLogger struct {
	prefix   string
	severity severity
	mx       sync.Mutex
}

func (l *LeveledLogger) header(s severity, t *time.Time) *bytes.Buffer {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%s %-5.5s %s", t.Format(timeFormat), severityName[s], l.prefix)
	return b
}

func (l *LeveledLogger) logln(s severity, v ...interface{}) {
	if l.severity >= s {
		t := time.Now()
		buf := l.header(s, &t)
		fmt.Fprintln(buf, v...)
		buf.WriteTo(os.Stderr)
	}
}

func (l *LeveledLogger) logf(s severity, format string, v ...interface{}) {
	if l.severity >= s {
		t := time.Now()
		buf := l.header(s, &t)
		fmt.Fprintf(buf, format, v...)
		if buf.Bytes()[buf.Len()-1] != '\n' {
			buf.WriteByte('\n')
		}
		buf.WriteTo(os.Stderr)
	}
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

func (l *LeveledLogger) Warningf(format string, v ...interface{}) {
	l.logf(WARNING, format, v...)
}

func (l *LeveledLogger) Warningln(v ...interface{}) {
	l.logln(WARNING, v...)
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

func New(level severity, prefix string) *LeveledLogger {
	return &LeveledLogger{
		prefix,
		level,
		sync.Mutex{},
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
	return logger.GetLevel()
}

func SetLevel(s severity) {
	logger.SetLevel(s)
}

func IsDebug() bool {
	return logger.IsDebug()
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Debugln(v ...interface{}) {
	logger.Debugln(v...)
}

func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Infoln(v ...interface{}) {
	logger.Infoln(v...)
}

func Warningf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}

func Warningln(v ...interface{}) {
	logger.Warningln(v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	logger.Errorln(v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	logger.Panicln(v...)
}
