package eventbroker

import (
	"fmt"
	"gearbox/eventbroker/states"
)


func defaultCallback(state states.Status) error {

	var err error

	fmt.Printf("CB state: %s\n", state.ShortString())

	return err
}

