package util

import (
	"github.com/gearboxworks/go-status/only"
	"strings"
)

func CharAt(s string, b byte) byte {
	sa := []rune(s)
	return byte(sa[b])
}

func After(str string, sub string) (s string) {
	for range only.Once {
		pos := strings.LastIndex(str, sub)
		if pos == -1 {
			break
		}
		pos += len(sub)
		if pos >= len(str) {
			break
		}
		s = str[pos:]
	}
	return s
}

func Before(str string, sub string) string {
	pos := strings.Index(str, sub)
	if pos == -1 {
		return ""
	}
	return str[0:pos]
}

type UniqueStrings map[string]bool

func NewUniqueStrings(size int) UniqueStrings {
	return make(UniqueStrings, size)
}

func (me UniqueStrings) ToSlice() []string {
	keys := make([]string, len(me))
	i := 0
	for k := range me {
		keys[i] = k
		i++
	}
	return keys
}

func Dashify(s string) string {
	r := strings.NewReplacer(" ", "-", "_", "-")
	return r.Replace(s)
}
