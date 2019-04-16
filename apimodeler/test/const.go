package test

import "gearbox/apimodeler"

const (
	FrobinatorType    = "frobinator"
	UnicornType       = "unicorn"
	FrobinatorsFilter = "/frobinators"
	UnicornFilter     = "/unicorns"
)

const (
	testableModelBasepath = "/foo"
)

var testableItemData = map[apimodeler.ItemId]apimodeler.ItemType{
	"foo": FrobinatorType,
	"bar": FrobinatorType,
	"baz": UnicornType,
}

var testableModelIdParams = apimodeler.IdParams{"foo", "bar"}
