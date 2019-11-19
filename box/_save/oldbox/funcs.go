package box

import (
	"strings"
	"unicode"
)

type KeyValue struct {
	Key	string
	Value string
}
type KeyValues []KeyValue

type KeyValueMap map[string]string
type KeyValuesMap map[string]KeyValueMap


func lineToKey(s string, splitOn rune) (ld KeyValue, ok bool) {

	ok = false
	foundKey := false
	stripSpace := false
	f := func(c rune) bool {
		switch {
		case c == splitOn:
			if foundKey {
				return false
			} else {
				foundKey = true
				stripSpace = true
				return true
			}

		case unicode.IsSpace(c) && stripSpace:
			return true

		default:
			stripSpace = false
			return false
			//return unicode.IsSpace(c)
		}
	}

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(s, f)
	if len(items) == 1 {
		ld.Key = items[0]
		ld.Value = ""
		ok = true
	} else if len(items) == 2 {
		ld.Key = items[0]
		ld.Value = items[1]
		ok = true
	}
	ld.Key = strings.TrimSuffix(strings.TrimPrefix(ld.Key, `"`), `"`)
	ld.Key = strings.TrimSuffix(strings.TrimPrefix(ld.Key, `{`), `}`)
	ld.Value = strings.TrimSuffix(strings.TrimPrefix(ld.Value, `"`), `"`)
	ld.Value = strings.TrimSuffix(strings.TrimPrefix(ld.Value, `{`), `}`)

	return
}


func decodeResponse(s string, splitOn rune) (dr KeyValues, ok bool) {

	ok = false

	lines := strings.Split(s, "\n")
	for _, l := range lines {
		kv, lineOk := lineToKey(l, splitOn)
		if lineOk == false {
			continue

		} else {
			// Return true if we have at least one key/value pair found.
			ok = true
			dr = append(dr, kv)
		}
		// fmt.Printf("items[%d]: '%s' = '%s'\n", i, foo.Key, foo.Value)
	}

	return
}

