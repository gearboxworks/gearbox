package tests

import (
	"gearbox"
	"gearbox/test"
	"testing"
)

func TestDottedVersion(t *testing.T) {
	test.StructMethodsTest(&DvTest{T: t})
}

func (me *DvTest) GetData() test.Table {
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
				parseDv: test.Args{Fail: true, Want: "non-integer major version in 'a.b.c'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Fail: true,
			In:   "1.b.c",
			Out: test.Out{
				parseDv: test.Args{Fail: true, Want: "non-integer minor version in '1.b.c'"},
			},
		}),
		test.NewFixture(&test.Fixture{
			Fail: true,
			In:   "1.1.c",
			Out: test.Out{
				parseDv: test.Args{Fail: true, Want: "non-integer patch version in '1.1.c'"},
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

var _ test.StructMethodTester = (*DvTest)(nil)

type DvTest struct {
	T *testing.T
}

func (me *DvTest) GetT() *testing.T {
	return me.T
}

func (me *DvTest) MakeNewObject(f *test.Fixture) (obj interface{}, err error) {
	var dv *gearbox.DottedVersion
	dv = gearbox.NewDottedVersion()
	err = dv.Parse(f.In)
	return dv, err
}

func (me *DvTest) GetOutput(f *test.Fixture) (got string) {
	dv := f.Obj.(*gearbox.DottedVersion)
	switch f.Name {
	case parseDv:
		got = dv.GetVersion()

	case getMajorMinor:
		got = dv.GetMajorMinor()

	case getMajor:
		got = dv.GetMajor()

	case getMinor:
		got = dv.GetMinor()

	case getPatch:
		got = dv.GetPatch()

	case getPrerelease:
		got = dv.GetPrerelease()

	case getMetadata:
		got = dv.GetMetadata()

	case getRevision:
		got = dv.GetRevision()

	case getRawDv:
		got = dv.GetRaw()

	}
	return got
}
