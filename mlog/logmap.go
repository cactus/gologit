package mlog

import (
	"bytes"
	"fmt"
	"io"
	"sort"
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
	var buf bytes.Buffer
	lm.WriteTo(&buf)
	return buf.String()
}

func (lm *LogMap) SortedString() string {
	var buf bytes.Buffer
	lm.SortedWriteTo(&buf)
	return buf.String()
}
