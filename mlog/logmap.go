package mlog

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"unsafe"
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

func stringtoslicebytetmp(s *string) []byte {
	// Return a slice referring to the actual string bytes.
	// This is only for use by internal compiler optimizations
	// that know that the slice won't be mutated.
	// The only such case today is:
	// for i, c := range []byte(str)

	sh := (*reflect.SliceHeader)(unsafe.Pointer(s))
	sh.Len = len(*s)
	sh.Cap = sh.Len
	return *(*[]byte)(unsafe.Pointer(sh))
}

func (lm *LogMap) SortedWriteTo(w io.Writer) (int64, error) {
	keys := lm.Keys()
	sort.Strings(keys)

	i := 0
	ilen := len(keys)
	for _, k := range keys {
		//w.Write([]byte(k))
		// this is a bit grotesque, but it avoids
		// an allocation. Since write will not mutate
		// the string, this *should* be safe.
		w.Write(stringtoslicebytetmp(&k))
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
