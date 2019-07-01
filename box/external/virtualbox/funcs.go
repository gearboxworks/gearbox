package virtualbox

import (
	"bytes"
	"fmt"
	"gearbox/ensure"
	"gearbox/eventbroker/eblog"
	"github.com/gearboxworks/go-status/only"
	"strings"
	"unicode"
)

func decodeResponse(s *bytes.Buffer, splitOn rune) (dr KeyValues, ok bool) {

	ok = false

	lines := strings.Split(s.String(), "\n")
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

func findFirstNic(vm VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		_, err = RunVm(vm, "list", "bridgedifs", "-s")
		if err != nil {
			break
		}

		var nic KeyValueMap
		dr, ok := decodeResponse(logger.GetStdout(), ':')
		if ok == true {
			var nics KeyValuesMap
			nics, ok = dr.decodeBridgeIfs()
			if ok == false {
				err = fmt.Errorf("no NICs found for VM '%s'", vm.GetName())
				break
			}

			for _, nic = range nics {
				if nic["FirstNic"] == "Yes" {
					break
				}
			}
		}

		if nic == nil {
			err = fmt.Errorf("no NICs found for VM '%s'", vm.GetName())
			break
		}

		logger.Debug("using NIC '%s' for VM '%s'", nic, vm.GetName())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}
