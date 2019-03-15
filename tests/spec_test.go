package tests

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
			In:    "@github.com/gearboxworks/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: true, Want: "invalid host '@github.com' in '@github.com/gearboxworks/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority and ServiceType w/invalid StackName",
			In:    "github.com/gearbox.works/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: true, Want: "invalid stack name 'gearbox.works' in spec 'github.com/gearbox.works/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName w/invalid ServiceType",
			In:    "gearboxworks/@dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: true, Want: "invalid role '@dbserver' in 'gearboxworks/@dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "invalid StackName w/ServiceType",
			In:    "gearbox#works/dbserver",
			Out: test.Out{
				getSpec: test.Args{Fail: true, Want: "invalid stack name 'gearbox#works' in spec 'gearbox#works/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType/Revision w/invalid Authority",
			In:    "github!/gearboxworks/dbserver:2",
			Out: test.Out{
				getSpec: test.Args{Fail: true, Want: "invalid stack name 'github!/gearboxworks' in spec 'github!/gearboxworks/dbserver:2'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType w/invalid Revision",
			In:    "gearboxworks/dbserver:1a",
			Out: test.Out{
				getSpec: test.Args{Fail: true, Want: "invalid version in 'gearboxworks/dbserver:1a'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType/Revision",
			In:    "gearboxworks/dbserver:2",
			Out: test.Out{
				getSpec:        test.Args{Fail: false, Want: "gearboxworks/dbserver:2"},
				getRawSpec:     test.Args{Fail: false, Want: "gearboxworks/dbserver:2"},
				getAuthority:   test.Args{Fail: false, Want: ""},
				getStackName:   test.Args{Fail: false, Want: "gearboxworks"},
				getServiceType: test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: false, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority/StackName/ServiceType/Revision",
			In:    "github.com/gearboxworks/dbserver:2",
			Out: test.Out{
				getSpec:        test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver:2"},
				getRawSpec:     test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver:2"},
				getAuthority:   test.Args{Fail: false, Want: "github.com"},
				getStackName:   test.Args{Fail: false, Want: "gearboxworks"},
				getServiceType: test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: false, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "StackName/ServiceType",
			In:    "gearboxworks/dbserver",
			Out: test.Out{
				getSpec:        test.Args{Fail: false, Want: "gearboxworks/dbserver"},
				getRawSpec:     test.Args{Fail: false, Want: "gearboxworks/dbserver"},
				getAuthority:   test.Args{Fail: false, Want: ""},
				getStackName:   test.Args{Fail: false, Want: "gearboxworks"},
				getServiceType: test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "ServiceType Only",
			In:    "dbserver",
			Out: test.Out{
				getSpec:        test.Args{Fail: false, Want: "dbserver"},
				getRawSpec:     test.Args{Fail: false, Want: "dbserver"},
				getAuthority:   test.Args{Fail: false, Want: ""},
				getStackName:   test.Args{Fail: false, Want: ""},
				getServiceType: test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Authority/StackName/ServiceType",
			In:    "github.com/gearboxworks/dbserver",
			Out: test.Out{
				getSpec:        test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver"},
				getRawSpec:     test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver"},
				getAuthority:   test.Args{Fail: false, Want: "github.com"},
				getStackName:   test.Args{Fail: false, Want: "gearboxworks"},
				getServiceType: test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:     test.Args{Fail: false, Want: ""},
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
