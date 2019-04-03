package unit

import (
	"gearbox/gearid"
	"gearbox/status"
	"gearbox/test"
	"testing"
)

func TestIdentity(t *testing.T) {
	test.StructMethodsTest(&IdTest{T: t})
}

func (me *IdTest) GetData() test.Table {
	return test.Table{
		test.NewFixture(&test.Fixture{
			Label: "No OrgName",
			In:    "php:7",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "php:7"},
				getRaw:     test.Args{Fail: false, Want: "php:7"},
				getOrgName: test.Args{Fail: false, Want: ""},
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
				getOrgName: test.Args{Fail: false, Want: "gearbox"},
				getType:    test.Args{Fail: false, Want: ""},
				getName:    test.Args{Fail: false, Want: "php"},
				getVersion: test.Args{Fail: false, Want: "7"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "OrgName/Type/Program/Version",
			In:    "wordpress/plugins/akismet:4.1.1",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "wordpress/plugins/akismet:4.1.1"},
				getRaw:     test.Args{Fail: false, Want: "wordpress/plugins/akismet:4.1.1"},
				getOrgName: test.Args{Fail: false, Want: "wordpress"},
				getType:    test.Args{Fail: false, Want: "plugins"},
				getName:    test.Args{Fail: false, Want: "akismet"},
				getVersion: test.Args{Fail: false, Want: "4.1.1"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full GearId",
			In:    "gearbox/php:7.2.1",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "gearbox/php:7.2.1"},
				getRaw:     test.Args{Fail: false, Want: "gearbox/php:7.2.1"},
				getOrgName: test.Args{Fail: false, Want: "gearbox"},
				getType:    test.Args{Fail: false, Want: ""},
				getName:    test.Args{Fail: false, Want: "php"},
				getVersion: test.Args{Fail: false, Want: "7.2.1"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full GearId with Revision",
			In:    "gearbox/php:7.2.1~r3",
			Out: test.Out{
				getId:      test.Args{Fail: false, Want: "gearbox/php:7.2.1~r3"},
				getRaw:     test.Args{Fail: false, Want: "gearbox/php:7.2.1~r3"},
				getOrgName: test.Args{Fail: false, Want: "gearbox"},
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
				getOrgName: test.Args{Fail: false, Want: "gearbox"},
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
	getOrgName = "GetOrgName"
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

func (me *IdTest) MakeNewObject(f *test.Fixture) (obj interface{}, sts status.Status) {
	gid := gearid.NewGearId()
	sts = gid.ParseString(f.In)
	return gid, sts
}

func (me *IdTest) GetOutput(f *test.Fixture) (got string) {
	gid := f.Obj.(*gearid.GearId)
	switch f.Name {
	case getId:
		got = gid.String()

	case getOrgName:
		got = string(gid.GetOrgName())

	case getType:
		got = string(gid.GetType())

	case getName:
		got = string(gid.GetName())

	case getVersion:
		got = gid.GetVersion().String()

	case getRaw:
		got = string(gid.GetRaw())

	}
	return got
}
