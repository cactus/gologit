package mlog

import (
	"bytes"
	"fmt"
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
		`level="info" msg="test one" x="1"`,
		"test one",
		[]interface{}{
			&ExtraMap{"x": 1},
		},
	},
}

func testPrint(t *testing.T, format string, arguments []interface{}, pattern string) {
	buf := new(bytes.Buffer)
	logger := New(buf, false, false)

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
	fmt.Println(pattern, line)
}

func TestAll(t *testing.T) {
	for _, testcase := range tests {
		testPrint(t, testcase.format, testcase.arguments, testcase.pattern)
		testPrint(t, testcase.format, testcase.arguments, testcase.pattern)
	}
}
