package tests

import (
	"gearbox"
	"gearbox/test"
	"testing"
)

func TestIdentity(t *testing.T) {
	test.StructMethodsTest(&IdTest{T: t})
}

func (me *IdTest) GetData() test.Table {
	return test.Table{
		test.NewFixture(&test.Fixture{
			Label: "No Org",
			In:    "php:7",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "php:7"},
				getRaw:     test.Args{Fail: false, Want: "php:7"},
				getGroup:   test.Args{Fail: false, Want: ""},
				getType:    test.Args{Fail: false, Want: ""},
				getName:    test.Args{Fail: false, Want: "php"},
				getVersion: test.Args{Fail: false, Want: "7"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Major",
			In:    "gearbox/php:7",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "gearbox/php:7"},
				getRaw:     test.Args{Fail: false, Want: "gearbox/php:7"},
				getGroup:   test.Args{Fail: false, Want: "gearbox"},
				getType:    test.Args{Fail: false, Want: ""},
				getName:    test.Args{Fail: false, Want: "php"},
				getVersion: test.Args{Fail: false, Want: "7"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Org/Type/Program/Version",
			In:    "wordpress/plugins/akismet:4.1.1",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "wordpress/plugins/akismet:4.1.1"},
				getRaw:     test.Args{Fail: false, Want: "wordpress/plugins/akismet:4.1.1"},
				getGroup:   test.Args{Fail: false, Want: "wordpress"},
				getType:    test.Args{Fail: false, Want: "plugins"},
				getName:    test.Args{Fail: false, Want: "akismet"},
				getVersion: test.Args{Fail: false, Want: "4.1.1"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full Identity",
			In:    "gearbox/php:7.2.1",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "gearbox/php:7.2.1"},
				getRaw:     test.Args{Fail: false, Want: "gearbox/php:7.2.1"},
				getGroup:   test.Args{Fail: false, Want: "gearbox"},
				getType:    test.Args{Fail: false, Want: ""},
				getName:    test.Args{Fail: false, Want: "php"},
				getVersion: test.Args{Fail: false, Want: "7.2.1"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full Identity with Revision",
			In:    "gearbox/php:7.2.1~r3",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "gearbox/php:7.2.1~r3"},
				getRaw:     test.Args{Fail: false, Want: "gearbox/php:7.2.1~r3"},
				getGroup:   test.Args{Fail: false, Want: "gearbox"},
				getType:    test.Args{Fail: false, Want: ""},
				getName:    test.Args{Fail: false, Want: "php"},
				getVersion: test.Args{Fail: false, Want: "7.2.1~r3"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Major/Minor",
			In:    "gearbox/php:7.2",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "gearbox/php:7.2"},
				getRaw:     test.Args{Fail: false, Want: "gearbox/php:7.2"},
				getGroup:   test.Args{Fail: false, Want: "gearbox"},
				getType:    test.Args{Fail: false, Want: ""},
				getName:    test.Args{Fail: false, Want: "php"},
				getVersion: test.Args{Fail: false, Want: "7.2"},
			},
		}),
	}
}

const (
	getId      = "GetId"
	getRaw     = "GetRaw"
	getGroup   = "GetGroup"
	getType    = "GetType"
	getName    = "GetName"
	getVersion = "GetVersion"
)

var _ test.StructMethodTester = (*IdTest)(nil)

type IdTest struct {
	T *testing.T
}

func (me *IdTest) GetT() *testing.T {
	return me.T
}

func (me *IdTest) MakeNewObject(f *test.Fixture) (obj interface{}, err error) {
	id := gearbox.NewIdentity()
	err = id.Parse(f.In)
	return id, err
}

func (me *IdTest) GetOutput(f *test.Fixture) (got string) {
	id := f.Obj.(*gearbox.Identity)
	switch f.Name {
	case getId:
		got = id.String()

	case getGroup:
		got = id.GetGroup()

	case getType:
		got = id.GetType()

	case getName:
		got = id.GetName()

	case getVersion:
		got = id.GetVersion().String()

	case getRaw:
		got = id.GetRaw()

	}
	return got
}
