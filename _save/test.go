package _save

import (
	"fmt"
	"github.com/gedex/inflector"
)

func main_test() {

	fmt.Printf("%s\n", inflector.Singularize("/projects"))
	fmt.Printf("%s\n", inflector.Pluralize("/project"))
}
