package ensure

import (
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"reflect"
)

func NotNil(obj interface{}, args ...interface{}) (err error) {
	for range only.Once {
		if obj != nil {
			break
		}
		v := reflect.ValueOf(obj)
		if v.Kind() != reflect.Ptr {
			break
		}
		if !v.IsNil() {
			break
		}
		msg := fmt.Sprintf("value of type '%T' is nil", obj)
		switch len(args) {
		case 0:
			err = fmt.Errorf(msg)
		case 1:
			err = fmt.Errorf("%s; %s", msg, args[0])
		default:
			extra := ""
			s, ok := args[0].(string)
			if !ok {
				err = fmt.Errorf("%s: %v", msg, args[0])
				break
			}
			extra = fmt.Sprintf(s, args[1:]...)
			err = fmt.Errorf("%s: %s", msg, extra)
		}
	}
	return err

}
