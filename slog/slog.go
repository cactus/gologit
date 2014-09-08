// Package gogitlog implements a very simple wrapper around the
// Go "log" package, providing support for a toggle-able debug flag
// and a couple of functions that log or not based on that flag.
package slog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var timeFormat = "2006-01-02T15:04:05.000000"
var logger = New(INFO)

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
	*log.Logger
	severity severity
	mx       sync.Mutex
}

func (l *LeveledLogger) header(s severity, t *time.Time) *bytes.Buffer {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "%s %-5.5s ", t.Format(timeFormat), severityName[s])
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

func New(level severity) *LeveledLogger {
	return &LeveledLogger{
		log.New(os.Stderr, "", log.LstdFlags),
		level,
		sync.Mutex{},
	}
}

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

func GetLevel() severity {
	return logger.severity
}

func SetLevel(s severity) {
	logger.severity = s
}

func IsDebug() bool {
	if logger.severity == DEBUG {
		return true
	}
	return false
}

func Debugf(format string, v ...interface{}) {
	logger.logf(DEBUG, format, v...)
}

func Debugln(v ...interface{}) {
	logger.logln(DEBUG, v...)
}

func Infof(format string, v ...interface{}) {
	logger.logf(INFO, format, v...)
}

func Infoln(v ...interface{}) {
	logger.logln(INFO, v...)
}

func Warningf(format string, v ...interface{}) {
	logger.logf(WARNING, format, v...)
}

func Warningln(v ...interface{}) {
	logger.logln(WARNING, v...)
}

func Errorf(format string, v ...interface{}) {
	logger.logf(ERROR, format, v...)
}

func Errorln(v ...interface{}) {
	logger.logln(ERROR, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.logf(FATAL, format, v...)
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	logger.logln(FATAL, v...)
	os.Exit(1)
}

// Logs to the default Logger. See Logger.Panicf
func Panicf(format string, v ...interface{}) {
	logger.logf(FATAL, format, v...)
	panic(fmt.Sprintf(format, v...))
}

// Logs to the default Logger. See Logger.Panicln
func Panicln(v ...interface{}) {
	logger.logln(FATAL, v...)
	panic(fmt.Sprintln(v...))
}
