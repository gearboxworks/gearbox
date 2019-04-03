package main

import (
	"fmt"
	"github.com/gedex/inflector"
)

func main() {

	fmt.Printf("%s\n", inflector.Singularize("/projects"))
	fmt.Printf("%s\n", inflector.Pluralize("/project"))
}
