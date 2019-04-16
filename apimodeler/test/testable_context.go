package test

import (
	"gearbox/apimodeler"
)

var NilTestableContext = (*TestableContext)(nil)
var _ apimodeler.Contexter = NilTestableContext

type TestableContext struct{}

func (me *TestableContext) Param(name string) (value string) {
	switch name {
	case "foo":
		value = "alpha"
	case "bar":
		value = "beta"
	}
	return value
}
