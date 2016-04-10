package mlog

import (
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkPLoggingBase(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("this is a test: %s", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingBaseSortedKeys(b *testing.B) {
	logger := New(ioutil.Discard, Lsort)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("this is a test: %s", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingDatetime(b *testing.B) {
	logger := New(ioutil.Discard, Ldatetime)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("this is a test: %s", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingDebugWithEnabled(b *testing.B) {
	logger := New(ioutil.Discard, Ldebug)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("this is a test: %s", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingDebugWithDisabled(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("this is a test: %s", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingLikeStdlib(b *testing.B) {
	logger := New(ioutil.Discard, Ldatetime)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("this is a test: %s", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPStdlibLog(b *testing.B) {
	logger := log.New(ioutil.Discard, "info: ", log.LstdFlags)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print(`msg="%s" %s="%d"`, "this is a test: %s", "x", 42)
		}
	})
}

func BenchmarkPStdlibLogShortfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "info: ", log.LstdFlags|log.Lshortfile)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print(`msg="%s" %s="%d"`, "this is a test: %s", "x", 42)
		}
	})
}

func BenchmarkPStdlibLogLongfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "info: ", log.LstdFlags|log.Llongfile)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print(`msg="%s" %s="%d"`, "this is a test: %s", "x", 42)
		}
	})
}
