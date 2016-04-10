package mlog

import (
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkSLoggingBase(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingBaseSortedKeys(b *testing.B) {
	logger := New(ioutil.Discard, Lsort)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingDatetime(b *testing.B) {
	logger := New(ioutil.Discard, Ldatetime)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingShortfile(b *testing.B) {
	logger := New(ioutil.Discard, Lshortfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingLongfile(b *testing.B) {
	logger := New(ioutil.Discard, Llongfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingDebugWithEnabled(b *testing.B) {
	logger := New(ioutil.Discard, Ldebug)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingDebugWithDisabled(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingLikeStdlib(b *testing.B) {
	logger := New(ioutil.Discard, Ldatetime)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", &LogMap{"x": 42})
	}
}

func BenchmarkSStdlibLog(b *testing.B) {
	logger := log.New(ioutil.Discard, "info: ", log.LstdFlags)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print(`msg="%s" %s="%d"`, "this is a test: %s", "x", 42)
	}
}

func BenchmarkSStdlibLogShortfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "info: ", log.LstdFlags|log.Lshortfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print(`msg="%s" %s="%d"`, "this is a test: %s", "x", 42)
	}
}

func BenchmarkSStdlibLogLongfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "info: ", log.LstdFlags|log.Llongfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print(`msg="%s" %s="%d"`, "this is a test: %s", "x", 42)
	}
}
