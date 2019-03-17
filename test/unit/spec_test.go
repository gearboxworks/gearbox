package unit

import (
	"gearbox"
	"gearbox/stat"
	"gearbox/test"
	"testing"
)

func TestSpec(t *testing.T) {
	test.StructMethodsTest(&SpecTest{T: t})
}

func (me *SpecTest) GetData() test.Table {
	return test.Table{
		test.NewFixture(&test.Fixture{
			Label: "Authority and ServiceType w/invalid StackName",
			In:    "@gearbox.works/wordpress/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid authority '@gearbox.works' in '@gearbox.works/wordpress/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority and ServiceType w/invalid StackName",
			In:    "gearbox.works/word.press/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid stack name 'word.press' in spec 'gearbox.works/word.press/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName w/invalid ServiceType",
			In:    "wordpress/@dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid role '@dbserver' in 'wordpress/@dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "invalid StackName w/ServiceType",
			In:    "gearbox.works/word#press/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid stack name 'word#press' in spec 'gearbox.works/word#press/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType/Revision w/invalid Authority",
			In:    "gearbox.works!/wordpress/dbserver:2",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid authority 'gearbox.works!' in 'gearbox.works!/wordpress/dbserver:2'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType w/invalid Revision",
			In:    "wordpress/dbserver:1a",
			Out: test.Out{
				getSpec: test.Args{Fail: T, Want: "invalid version in 'wordpress/dbserver:1a'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType/Revision",
			In:    "wordpress/dbserver:2",
			Out: test.Out{
				getSpec:        test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver:2"},
				getRawSpec:     test.Args{Fail: F, Want: "wordpress/dbserver:2"},
				getAuthority:   test.Args{Fail: F, Want: "gearbox.works"},
				getStackName:   test.Args{Fail: F, Want: "wordpress"},
				getServiceType: test.Args{Fail: F, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: F, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority/StackName/ServiceType/Revision",
			In:    "gearbox.works/wordpress/dbserver:2",
			Out: test.Out{
				getSpec:        test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver:2"},
				getRawSpec:     test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver:2"},
				getAuthority:   test.Args{Fail: F, Want: "gearbox.works"},
				getStackName:   test.Args{Fail: F, Want: "wordpress"},
				getServiceType: test.Args{Fail: F, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: F, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType",
			In:    "wordpress/dbserver",
			Out: test.Out{
				getSpec:        test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver"},
				getRawSpec:     test.Args{Fail: F, Want: "wordpress/dbserver"},
				getAuthority:   test.Args{Fail: F, Want: "gearbox.works"},
				getStackName:   test.Args{Fail: F, Want: "wordpress"},
				getServiceType: test.Args{Fail: F, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: F, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "ServiceType Only",
			In:    "dbserver",
			Out: test.Out{
				getSpec:        test.Args{Fail: T, Want: "invalid spec 'dbserver'"},
				getRawSpec:     test.Args{Fail: T, Want: "invalid spec 'dbserver'"},
				getAuthority:   test.Args{Fail: T, Want: "invalid spec 'dbserver'"},
				getStackName:   test.Args{Fail: T, Want: "invalid spec 'dbserver'"},
				getServiceType: test.Args{Fail: T, Want: "invalid spec 'dbserver'"},
				getSpecVer:     test.Args{Fail: T, Want: "invalid spec 'dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority/StackName/ServiceType",
			In:    "gearbox.works/wordpress/dbserver",
			Out: test.Out{
				getSpec:        test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver"},
				getRawSpec:     test.Args{Fail: F, Want: "gearbox.works/wordpress/dbserver"},
				getAuthority:   test.Args{Fail: F, Want: "gearbox.works"},
				getStackName:   test.Args{Fail: F, Want: "wordpress"},
				getServiceType: test.Args{Fail: F, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: F, Want: ""},
			},
		}),
	}
}

const (
	getSpec        = "GetSpec"
	getRawSpec     = "GetRawSpec"
	getAuthority   = "GetAuthority"
	getStackName   = "GetStackName"
	getServiceType = "GetType"
	getSpecVer     = "GetSpecVer"
)

var _ test.StructMethodTester = (*SpecTest)(nil)

type SpecTest struct {
	T *testing.T
}

func (me *SpecTest) GetT() *testing.T {
	return me.T
}

func (me *SpecTest) MakeNewObject(f *test.Fixture) (obj interface{}, status stat.Status) {
	spec := gearbox.NewSpec()
	status = spec.Parse(f.In)
	return spec, status
}

func (me *SpecTest) GetOutput(f *test.Fixture) (got string) {
	spec := f.Obj.(*gearbox.Spec)
	switch f.Name {
	case getSpec:
		got = spec.String()

	case getAuthority:
		got = string(spec.GetAuthority())

	case getStackName:
		got = string(spec.GetStackName())

	case getServiceType:
		got = spec.GetType()

	case getSpecVer:
		got = spec.GetVersion()

	case getRawSpec:
		got = spec.GetRaw()

	}
	return got
}
