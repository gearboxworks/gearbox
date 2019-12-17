package eventbroker

import (
	"fmt"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
)

type SimpleStatus map[msgs.Address]states.State

func (me SimpleStatus) String() string {
	var ret string
	for range only.Once {
		for k, v := range me {
			ret += fmt.Sprintf("%s %s\n", k, v)
		}
	}
	return ret
}
