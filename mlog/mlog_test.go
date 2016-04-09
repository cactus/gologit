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
	format    string
	arguments []interface{}
}

var tests = []tester{
	{
		`level="info" msg="test one 1" x="y"`,
		"test one %d",
		[]interface{}{
			1,
			&ExtraMap{"x": "y"},
		},
	},
	{
		`level="info" msg="test one 1" x="y"`,
		"test one %d",
		[]interface{}{
			&ExtraMap{"x": "y"},
			1,
		},
	},
	{
		`level="info" msg="test one %!s\(int=1\)" x="y"`,
		"test one %s",
		[]interface{}{
			1,
			&ExtraMap{"x": "y"},
		},
	},
	{
		`level="info" msg="test one" x="y" t="u"`,
		"test one",
		[]interface{}{
			ExtraMap{"x": "y"},
			&ExtraMap{"t": "u"},
		},
	},
	{
		`level="info" msg="test one" x="1" y="2" z="3"`,
		"test one",
		[]interface{}{
			&ExtraMap{
				"x": 1,
				"y": 2,
				"z": 3,
			},
		},
	},
	{
		`level="info" msg="test: test 1"`,
		"test: %s %d",
		[]interface{}{
			"test", 1,
		},
	},
}

func testPrint(t *testing.T, format string, arguments []interface{}, pattern string) {
	buf := new(bytes.Buffer)
	logger := New(buf, Lbase|Lsort)

	logger.Printf(format, arguments...)
	line := buf.String()
	line = line[0 : len(line)-1]
	pattern = "^" + pattern + "$"
	matched, err := regexp.MatchString(pattern, line)
	if err != nil {
		t.Fatal("pattern did not compile:", err)
	}
	if !matched {
		t.Errorf("log output should match %q is %q", pattern, line)
	}
}

func TestAll(t *testing.T) {
	for _, testcase := range tests {
		testPrint(t, testcase.format, testcase.arguments, testcase.pattern)
		testPrint(t, testcase.format, testcase.arguments, testcase.pattern)
	}
}

func BenchmarkSLoggingBase(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf("this is a test: %s", "test", &ExtraMap{"x": 42})
	}
}

func BenchmarkSLoggingTime(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf("this is a test: %s", "test", &ExtraMap{"x": 42})
	}
}

func BenchmarkSLoggingSortedKeys(b *testing.B) {
	logger := New(ioutil.Discard, Lsort)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf("this is a test: %s", "test", &ExtraMap{"x": 42})
	}
}

func BenchmarkSLoggingDebugEnabled(b *testing.B) {
	logger := New(ioutil.Discard, Ldebug)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debugf("this is a test: %s", "test", &ExtraMap{"x": 42})
	}
}

func BenchmarkSLoggingDebugDisabled(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debugf("this is a test: %s", "test", &ExtraMap{"x": 42})
	}
}

func BenchmarkSLoggingLikeStdlib(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf(`this is a test: %s x=%d`, "test", 42)
	}
}

func BenchmarkSStdlibLog(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf(`this is a test: %s x=%d`, "test", 42)
	}
}

func BenchmarkSStdlibLogShortfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags|log.Lshortfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Printf(`this is a test: %s x=%d`, "test", 42)
	}
}

func BenchmarkPLoggingBase(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("this is a test: %s", "test", &ExtraMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingTime(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("this is a test: %s", "test", &ExtraMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingSortedKeys(b *testing.B) {
	logger := New(ioutil.Discard, Lsort)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("this is a test: %s", "test", &ExtraMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingDebugEnabled(b *testing.B) {
	logger := New(ioutil.Discard, Ldebug)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debugf("this is a test: %s", "test", &ExtraMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingDebugDisabled(b *testing.B) {
	logger := New(ioutil.Discard, Lbase)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debugf("this is a test: %s", "test", &ExtraMap{"x": 42})
		}
	})
}

func BenchmarkPLoggingLikeStdlib(b *testing.B) {
	logger := New(ioutil.Discard, Ltime)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf(`this is a test: %s x=%d`, "test", 42)
		}
	})
}

func BenchmarkPStdlibLog(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf(`this is a test: %s x=%d`, "test", 42)
		}
	})
}

func BenchmarkPStdlibLogShortfile(b *testing.B) {
	logger := log.New(ioutil.Discard, "debug: ", log.LstdFlags|log.Lshortfile)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf(`this is a test: %s x=%d`, "test", 42)
		}
	})
}
