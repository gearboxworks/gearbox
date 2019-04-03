package main

import (
	"fmt"
)

func main() {
	fmt.Printf("%t\n", IsError(nil))
}

type Status struct {
	Cause error
}

func (s Status) Error() string {
	return s.Cause.Error()
}
func (s *Status) IsError() bool {
	return s.Cause != nil
}
func IsError(err error) bool {
	status, ok := err.(Status)
	if !ok {
		return err != nil
	}
	return sts.IsError(sts)
}
