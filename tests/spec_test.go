package tests

import (
	"gearbox"
	"gearbox/test"
	"testing"
)

func TestSpec(t *testing.T) {
	test.StructMethodsTest(&SpecTest{T: t})
}

func (me *SpecTest) GetData() test.Table {
	return test.Table{
		test.NewFixture(&test.Fixture{
			Label: "Host and Role w/invalid Namespace",
			In:    "@github.com/gearboxworks/dbserver",
			Out: test.Out{
				getSpec:      test.Args{Fail: true, Want: "invalid host '@github.com' in '@github.com/gearboxworks/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Host and Role w/invalid Namespace",
			In:    "github.com/gearbox.works/dbserver",
			Out: test.Out{
				getSpec:      test.Args{Fail: true, Want: "invalid namespace 'gearbox.works' in spec 'github.com/gearbox.works/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Namespace w/invalid Role",
			In:    "gearboxworks/@dbserver",
			Out: test.Out{
				getSpec:      test.Args{Fail: true, Want: "invalid role '@dbserver' in 'gearboxworks/@dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "invalid Namespace w/Role",
			In:    "gearbox#works/dbserver",
			Out: test.Out{
				getSpec:      test.Args{Fail: true, Want: "invalid namespace 'gearbox#works' in spec 'gearbox#works/dbserver'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Namespace/Role/Version w/invalid Host",
			In:    "github!/gearboxworks/dbserver:2",
			Out: test.Out{
				getSpec:      test.Args{Fail: true, Want: "invalid namespace 'github!/gearboxworks' in spec 'github!/gearboxworks/dbserver:2'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Namespace/Role w/invalid Version",
			In:    "gearboxworks/dbserver:1a",
			Out: test.Out{
				getSpec:      test.Args{Fail: true, Want: "invalid version in 'gearboxworks/dbserver:1a'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Namespace/Role/Version",
			In:    "gearboxworks/dbserver:2",
			Out: test.Out{
				getSpec:      test.Args{Fail: false, Want: "gearboxworks/dbserver:2"},
				getRawSpec:   test.Args{Fail: false, Want: "gearboxworks/dbserver:2"},
				getHost:      test.Args{Fail: false, Want: ""},
				getNamespace: test.Args{Fail: false, Want: "gearboxworks"},
				getRole:      test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:   test.Args{Fail: false, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Host/Namespace/Role/Version",
			In:    "github.com/gearboxworks/dbserver:2",
			Out: test.Out{
				getSpec:      test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver:2"},
				getRawSpec:   test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver:2"},
				getHost:      test.Args{Fail: false, Want: "github.com"},
				getNamespace: test.Args{Fail: false, Want: "gearboxworks"},
				getRole:      test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:   test.Args{Fail: false, Want: "2"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Namespace/Role",
			In:    "gearboxworks/dbserver",
			Out: test.Out{
				getSpec:      test.Args{Fail: false, Want: "gearboxworks/dbserver"},
				getRawSpec:   test.Args{Fail: false, Want: "gearboxworks/dbserver"},
				getHost:      test.Args{Fail: false, Want: ""},
				getNamespace: test.Args{Fail: false, Want: "gearboxworks"},
				getRole:      test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:   test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Role Only",
			In:    "dbserver",
			Out: test.Out{
				getSpec:      test.Args{Fail: false, Want: "dbserver"},
				getRawSpec:   test.Args{Fail: false, Want: "dbserver"},
				getHost:      test.Args{Fail: false, Want: ""},
				getNamespace: test.Args{Fail: false, Want: ""},
				getRole:      test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:   test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Host/Namespace/Role",
			In:    "github.com/gearboxworks/dbserver",
			Out: test.Out{
				getSpec:      test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver"},
				getRawSpec:   test.Args{Fail: false, Want: "github.com/gearboxworks/dbserver"},
				getHost:      test.Args{Fail: false, Want: "github.com"},
				getNamespace: test.Args{Fail: false, Want: "gearboxworks"},
				getRole:      test.Args{Fail: false, Want: "dbserver"},
				getSpecVer:   test.Args{Fail: false, Want: ""},
			},
		}),
	}
}

const (
	getSpec      = "GetSpec"
	getRawSpec   = "GetRawSpec"
	getHost      = "GetHost"
	getNamespace = "GetNamespace"
	getRole      = "GetRole"
	getSpecVer   = "GetSpecVer"
)

var _ test.StructMethodTester = (*SpecTest)(nil)

type SpecTest struct {
	T *testing.T
}

func (me *SpecTest) GetT() *testing.T {
	return me.T
}

func (me *SpecTest) MakeNewObject(f *test.Fixture) (obj interface{}, err error) {
	spec := gearbox.NewSpec()
	err = spec.Parse(f.In)
	return spec, err
}

func (me *SpecTest) GetOutput(f *test.Fixture) (got string) {
	spec := f.Obj.(*gearbox.Spec)
	switch f.Name {
	case getSpec:
		got = spec.String()

	case getHost:
		got = spec.GetHost()

	case getNamespace:
		got = spec.GetNamespace()

	case getRole:
		got = spec.GetRole()

	case getSpecVer:
		got = spec.GetVersion()

	case getRawSpec:
		got = spec.GetRaw()

	}
	return got
}
