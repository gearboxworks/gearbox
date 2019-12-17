package test

import "gearbox/apiworks"

const (
	FrobinatorType    = "frobinator"
	UnicornType       = "unicorn"
	FrobinatorsFilter = "/frobinators"
	UnicornFilter     = "/unicorns"
)

const (
	testableModelBasepath = "/foo"
)

var testableItemData = map[apiworks.ItemId]apiworks.ItemType{
	"foo": FrobinatorType,
	"bar": FrobinatorType,
	"baz": UnicornType,
}

var testableModelIdParams = apiworks.IdParams{"foo", "bar"}
