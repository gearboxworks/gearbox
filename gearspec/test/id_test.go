package test

import (
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/status"
	"gearbox/test"
	"testing"
)

func TestId(t *testing.T) {
	test.StructMethodsTest(&SpecTest{T: t})
}

func (me *SpecTest) GetData() test.Table {
	return test.Table{
		test.NewFixture(&test.Fixture{
			Label: "Authority and Role w/invalid Stackname",
			In:    "@gearbox.works/wordpress/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid authority '@gearbox.works' in '@gearbox.works/wordpress/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority and Role w/invalid Stackname",
			In:    "gearbox.works/word.press/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid stack name 'word.press' in stack ID 'gearbox.works/word.press/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Stackname w/invalid Role",
			In:    "wordpress/@dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid role '@dbserver' in 'wordpress/@dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "invalid Stackname w/Role",
			In:    "gearbox.works/word#press/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid stack name 'word#press' in stack ID 'gearbox.works/word#press/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Stackname/Role/Revision w/invalid Authority",
			In:    "gearbox.works!/wordpress/dbserver:2",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid authority 'gearbox.works!' in 'gearbox.works!/wordpress/dbserver:2'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Stackname/Role w/invalid Revision",
			In:    "wordpress/dbserver:1a",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid version in 'wordpress/dbserver:1a'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Stackname/Role/Revision",
			In:    "wordpress/dbserver:2",
			Out: test.Out{
				getSpec:         test.Args{Fail: F, Want: global.DefaultAuthority + "/wordpress/dbserver:2"},
				getRawSpec:      test.Args{Fail: F, Want: "wordpress/dbserver:2"},
				getAuthority:    test.Args{Fail: F, Want: global.DefaultAuthority},
				getStackName:    test.Args{Fail: F, Want: "wordpress"},
				getServiceType:  test.Args{Fail: F, Want: "dbserver"},
				getSpecRevision: test.Args{Fail: F, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority/Stackname/Role/Revision",
			In:    global.DefaultAuthority + "/wordpress/dbserver:2",
			Out: test.Out{
				getSpec:         test.Args{Fail: F, Want: global.DefaultAuthority + "/wordpress/dbserver:2"},
				getRawSpec:      test.Args{Fail: F, Want: global.DefaultAuthority + "/wordpress/dbserver:2"},
				getAuthority:    test.Args{Fail: F, Want: global.DefaultAuthority},
				getStackName:    test.Args{Fail: F, Want: "wordpress"},
				getServiceType:  test.Args{Fail: F, Want: "dbserver"},
				getSpecRevision: test.Args{Fail: F, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Stackname/Role",
			In:    "wordpress/dbserver",
			Out: test.Out{
				getSpec:         test.Args{Fail: F, Want: global.DefaultAuthority + "/wordpress/dbserver"},
				getRawSpec:      test.Args{Fail: F, Want: "wordpress/dbserver"},
				getAuthority:    test.Args{Fail: F, Want: global.DefaultAuthority},
				getStackName:    test.Args{Fail: F, Want: "wordpress"},
				getServiceType:  test.Args{Fail: F, Want: "dbserver"},
				getSpecRevision: test.Args{Fail: F, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Role Only",
			In:    "dbserver",
			Out: test.Out{
				getSpec:         test.Args{Fail: T, Want: "invalid gearspec ID 'dbserver'"},
				getRawSpec:      test.Args{Fail: T, Want: "invalid gearspec ID 'dbserver'"},
				getAuthority:    test.Args{Fail: T, Want: "invalid gearspec ID 'dbserver'"},
				getStackName:    test.Args{Fail: T, Want: "invalid gearspec ID 'dbserver'"},
				getServiceType:  test.Args{Fail: T, Want: "invalid gearspec ID 'dbserver'"},
				getSpecRevision: test.Args{Fail: T, Want: "invalid gearspec ID 'dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority/Stackname/Role",
			In:    "gearbox.works/wordpress/dbserver",
			Out: test.Out{
				getSpec:         test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver"},
				getRawSpec:      test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver"},
				getAuthority:    test.Args{Fail: F, Want: "gearbox.works"},
				getStackName:    test.Args{Fail: F, Want: "wordpress"},
				getServiceType:  test.Args{Fail: F, Want: "dbserver"},
				getSpecRevision: test.Args{Fail: F, Want: ""},
			},
		}),
	}
}

const (
	getSpec         = "GetSpec"
	getRawSpec      = "GetRawSpec"
	getAuthority    = "GetAuthority"
	getStackName    = "GetStackname"
	getServiceType  = "GetType"
	getSpecRevision = "GetSpecVer"
)

var _ test.StructMethodTester = (*SpecTest)(nil)

type SpecTest struct {
	T *testing.T
}

func (me *SpecTest) GetT() *testing.T {
	return me.T
}

func (me *SpecTest) MakeNewObject(f *test.Fixture) (obj interface{}, sts status.Status) {
	gsi := gearspec.NewGearspec()
	sts = gsi.ParseString(f.In)
	return gsi, sts
}

func (me *SpecTest) GetOutput(f *test.Fixture) (got string) {
	gsi := f.Obj.(*gearspec.Gearspec)
	switch f.Name {
	case getSpec:
		got = gsi.String()

	case getAuthority:
		got = string(gsi.GetAuthority())

	case getStackName:
		got = string(gsi.GetStackname())

	case getServiceType:
		got = string(gsi.GetRole())

	case getSpecRevision:
		got = string(gsi.GetRevision())

	case getRawSpec:
		got = string(gsi.GetRaw())

	}
	return got
}
