package test

import (
	"gearbox/test"
	"gearbox/version"
	"github.com/gearboxworks/go-status"
	"testing"
)

func TestVersion(t *testing.T) {
	test.StructMethodsTest(&VerTest{T: t})
}

func (me *VerTest) GetData() test.Table {
	return test.Table{
		test.NewFixture(&test.Fixture{
			Label: "Full version plus Pre-release",
			In:    "1.27.3-alpha.1.a",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "1.27.3-alpha.1.a"},
				getRawDv:      test.Args{Fail: false, Want: "1.27.3-alpha.1.a"},
				getMajorMinor: test.Args{Fail: false, Want: "1.27"},
				getMajor:      test.Args{Fail: false, Want: "1"},
				getMinor:      test.Args{Fail: false, Want: "27"},
				getPatch:      test.Args{Fail: false, Want: "3"},
				getPrerelease: test.Args{Fail: false, Want: "alpha.1.a"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full version plus invalid Pre-release",
			In:    "1.27.3-beta!",
			Out: test.Out{
				parseDv: test.Args{Fail: true, Want: "pre-release in '1.27.3-beta!' is invalid semver"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full version plus invalid Pre-release",
			In:    "1.27.3+exp&sha.5114f85",
			Out: test.Out{
				parseDv: test.Args{Fail: true, Want: "build metadata in '1.27.3+exp&sha.5114f85' is not valid semver"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full version with Revision",
			In:    "1.27.3~r4",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "1.27.3~r4"},
				getRawDv:      test.Args{Fail: false, Want: "1.27.3~r4"},
				getMajorMinor: test.Args{Fail: false, Want: "1.27"},
				getMajor:      test.Args{Fail: false, Want: "1"},
				getMinor:      test.Args{Fail: false, Want: "27"},
				getPatch:      test.Args{Fail: false, Want: "3"},
				getRevision:   test.Args{Fail: false, Want: "4"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full version plus Build Metadata",
			In:    "1.27.3+exp.sha.5114f85",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "1.27.3+exp.sha.5114f85"},
				getRawDv:      test.Args{Fail: false, Want: "1.27.3+exp.sha.5114f85"},
				getMajorMinor: test.Args{Fail: false, Want: "1.27"},
				getMajor:      test.Args{Fail: false, Want: "1"},
				getMinor:      test.Args{Fail: false, Want: "27"},
				getPatch:      test.Args{Fail: false, Want: "3"},
				getMetadata:   test.Args{Fail: false, Want: "exp.sha.5114f85"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full version plus Pre-release and Build Metadata",
			In:    "1.27.3-beta+exp.sha.5114f85",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "1.27.3-beta+exp.sha.5114f85"},
				getRawDv:      test.Args{Fail: false, Want: "1.27.3-beta+exp.sha.5114f85"},
				getMajorMinor: test.Args{Fail: false, Want: "1.27"},
				getMajor:      test.Args{Fail: false, Want: "1"},
				getMinor:      test.Args{Fail: false, Want: "27"},
				getPatch:      test.Args{Fail: false, Want: "3"},
				getPrerelease: test.Args{Fail: false, Want: "beta"},
				getMetadata:   test.Args{Fail: false, Want: "exp.sha.5114f85"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Full version",
			In:    "1.27.3",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "1.27.3"},
				getRawDv:      test.Args{Fail: false, Want: "1.27.3"},
				getMajorMinor: test.Args{Fail: false, Want: "1.27"},
				getMajor:      test.Args{Fail: false, Want: "1"},
				getMinor:      test.Args{Fail: false, Want: "27"},
				getPatch:      test.Args{Fail: false, Want: "3"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Major/Minor version",
			In:    "1.11",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "1.11"},
				getRawDv:      test.Args{Fail: false, Want: "1.11"},
				getMajorMinor: test.Args{Fail: false, Want: "1.11"},
				getMajor:      test.Args{Fail: false, Want: "1"},
				getMinor:      test.Args{Fail: false, Want: "11"},
				getPatch:      test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Major version",
			In:    "12",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "12"},
				getRawDv:      test.Args{Fail: false, Want: "12"},
				getMajorMinor: test.Args{Fail: false, Want: "12"},
				getMajor:      test.Args{Fail: false, Want: "12"},
				getMinor:      test.Args{Fail: false, Want: ""},
				getPatch:      test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Major/Minor version",
			In:    "1.11",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "1.11"},
				getRawDv:      test.Args{Fail: false, Want: "1.11"},
				getMajorMinor: test.Args{Fail: false, Want: "1.11"},
				getMajor:      test.Args{Fail: false, Want: "1"},
				getMinor:      test.Args{Fail: false, Want: "11"},
				getPatch:      test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			Label: "Major version",
			In:    "12",
			Out: test.Out{
				parseDv:       test.Args{Fail: false, Want: "12"},
				getRawDv:      test.Args{Fail: false, Want: "12"},
				getMajorMinor: test.Args{Fail: false, Want: "12"},
				getMajor:      test.Args{Fail: false, Want: "12"},
				getMinor:      test.Args{Fail: false, Want: ""},
				getPatch:      test.Args{Fail: false, Want: ""},
			},
		}),
		test.NewFixture(&test.Fixture{
			In: "a.b.c",
			Out: test.Out{
				parseDv:       test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				getRawDv:      test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				getMajorMinor: test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				getMajor:      test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				getMinor:      test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				getPatch:      test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				getPrerelease: test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				getMetadata:   test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				//				parseDv: test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Fail: true,
			In:   "1.b.c",
			Out: test.Out{
				parseDv:       test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				getRawDv:      test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				getMajorMinor: test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				getMajor:      test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				getMinor:      test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				getPatch:      test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				getPrerelease: test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				getMetadata:   test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
				//				parseDv: test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
				//				parseDv: test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Fail: true,
			In:   "1.1.c",
			Out: test.Out{
				parseDv:       test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				getRawDv:      test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				getMajorMinor: test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				getMajor:      test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				getMinor:      test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				getPatch:      test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				getPrerelease: test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				getMetadata:   test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
				//				parseDv: test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
			},
		}),
	}

}

const (
	parseDv       = "Parse"
	getRawDv      = "GetRaw"
	getMajorMinor = "GetMajorMinor"
	getMajor      = "GetMajor"
	getMinor      = "GetMinor"
	getPatch      = "GetPatch"
	getPrerelease = "GetPrerelease"
	getMetadata   = "GetMetadata"
	getRevision   = "GetRevision"
)

var _ test.StructMethodTester = (*VerTest)(nil)

type VerTest struct {
	T *testing.T
}

func (me *VerTest) GetT() *testing.T {
	return me.T
}

func (me *VerTest) MakeNewObject(f *test.Fixture) (obj interface{}, sts status.Status) {
	var ver *version.Version
	ver = version.NewVersion()
	sts = ver.ParseString(f.In)
	return ver, sts
}

func (me *VerTest) GetOutput(f *test.Fixture) (got string) {
	ver := f.Obj.(*version.Version)
	switch f.Name {
	case parseDv:
		got = string(ver.GetIdentifier())

	case getMajorMinor:
		got = string(ver.GetMajorMinor())

	case getMajor:
		got = string(ver.GetMajor())

	case getMinor:
		got = string(ver.GetMinor())

	case getPatch:
		got = string(ver.GetPatch())

	case getPrerelease:
		got = string(ver.GetPrerelease())

	case getMetadata:
		got = string(ver.GetMetadata())

	case getRevision:
		got = string(ver.GetRevision())

	case getRawDv:
		got = string(ver.GetRaw())

	}
	return got
}
