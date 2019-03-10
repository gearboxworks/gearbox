package util

import (
	"gearbox/only"
	"strings"
)

func CharAt(s string, b byte) byte {
	sa:=[]rune(s)
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

