package mlog

import (
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
)

type tester struct {
	pattern   string
	message   string
	arguments []interface{}
}

var tests = []tester{
	{
		`level="info" msg="test one %d" x="y" extra1="1"`,
		"test one %d",
		[]interface{}{
			1,
			&LogMap{"x": "y"},
		},
	},
	{
		`level="info" msg="test one %d" x="y" extra1="1"`,
		"test one %d",
		[]interface{}{
			&LogMap{"x": "y"},
			1,
		},
	},
	{
		`level="info" msg="test one" x="y" y="z" t="u" u="v"`,
		"test one",
		[]interface{}{
			LogMap{"x": "y", "y": "z"},
			&LogMap{"t": "u", "u": "v"},
		},
	},
	{
		`level="info" msg="test one" x="y" y="z" t="u" u="v"`,
		"test one",
		[]interface{}{
			LogMap{"y": "z", "x": "y"},
			&LogMap{"u": "v", "t": "u"},
		},
	},
	{
		`level="info" msg="test one" x="1" y="2" z="3"`,
		"test one",
		[]interface{}{
			&LogMap{
				"x": 1,
				"y": 2,
				"z": 3,
			},
		},
	},
	{
		`level="info" msg="test: %s %d" extra1="test" extra2="1"`,
		"test: %s %d",
		[]interface{}{
			"test", 1,
		},
	},
}

func testInfo(t *testing.T, message string, arguments []interface{}, pattern string) {
	buf := new(bytes.Buffer)
	logger := New(buf, Lbase|Lsort)

	logger.Info(message, arguments...)
	line := buf.String()
	line = line[0 : len(line)-1]
	pattern = "^" + pattern + "$"
	matched, err := regexp.MatchString(pattern, line)
	if err != nil {
		t.Fatal("pattern did not compile:", err)
	}
	if !matched {
		t.Errorf("log output should match\n%12s %q\n%12s %q",
			"expected:", pattern[1:len(pattern)-1],
			"actual:", line)
	}
}

func TestAll(t *testing.T) {
	for _, testcase := range tests {
		testInfo(t, testcase.message, testcase.arguments, testcase.pattern)
		testInfo(t, testcase.message, testcase.arguments, testcase.pattern)
	}
}

func BenchmarkSLoggingBase(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", "test", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingTime(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", "test", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingSortedKeys(b *testing.B) {
	logger := New(ioutil.Discard, Lsort)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("this is a test: %s", "test", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingDebugEnabled(b *testing.B) {
	logger := New(ioutil.Discard, Ldebug)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("this is a test: %s", "test", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingDebugDisabled(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("this is a test: %s", "test", &LogMap{"x": 42})
	}
}

func BenchmarkSLoggingLikeStdlib(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(`this is a test: %s x=%d`, "test", 42)
	}
}

func BenchmarkSStdlibLog(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print(`this is a test: %s x=%d`, "test", 42)
	}
}

func BenchmarkSStdlibLogShortfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags|log.Lshortfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print(`this is a test: %s x=%d`, "test", 42)
	}
}

func BenchmarkPLoggingBase(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("this is a test: %s", "test", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingTime(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("this is a test: %s", "test", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingSortedKeys(b *testing.B) {
	logger := New(ioutil.Discard, Lsort)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("this is a test: %s", "test", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingDebugEnabled(b *testing.B) {
	logger := New(ioutil.Discard, Ldebug)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("this is a test: %s", "test", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingDebugDisabled(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("this is a test: %s", "test", &LogMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingLikeStdlib(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(`this is a test: %s x=%d`, "test", 42)
		}
	})
}

func BenchmarkPStdlibLog(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print(`this is a test: %s x=%d`, "test", 42)
		}
	})
}

func BenchmarkPStdlibLogShortfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags|log.Lshortfile)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print(`this is a test: %s x=%d`, "test", 42)
		}
	})
}
