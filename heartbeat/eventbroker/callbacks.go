package eventbroker

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/states"
)


func defaultCallback(state states.Status) error {

	var err error

	fmt.Printf("CB state: %s\n", state.ShortString())

	return err
}

