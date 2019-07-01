package ebutil

import "fmt"

func LogError(err error) {
	fmt.Printf(err.Error())
}
