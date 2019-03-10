package test

import (
	"testing"
)

type Fixture struct {
	Label string
	Name  string
	Fail  bool
	In    string
	Out   Out
	T     *testing.T
	Obj   interface{}
}

type Out map[string]Args

type Args struct {
	Fail bool
	Want string
}

func NewFixture(args *Fixture) *Fixture {
	f := &Fixture{}
	*f = Fixture(*args)
	return f
}

type Table []*Fixture

func NoFailWrongOutput(t *testing.T, label string, name string, input string, got string, want Args, err error) bool {
	if !want.Fail && err == nil {
		if got != want.Want {
			t.Errorf("nofail: %s.%s, want: %s, got: %s, input: %s", label, name, want.Want, got, input)
		}
		return true
	}
	return false
}

func FailWrongError(t *testing.T, label string, name string, input string, want Args, err error) bool {
	if want.Fail {
		if err != nil && err.Error() != want.Want {
			t.Errorf("fail: %s.%s, want: %s, got: %s, input: %s", label, name, want.Want, err.Error(), input)
		}
		return true
	}
	return false

}

type StructMethodTester interface {
	MakeNewObject(f *Fixture) (obj interface{}, err error)
	GetOutput(f *Fixture) string
	GetT() *testing.T
	GetData() Table
}

func StructMethodsTest(smt StructMethodTester) {
	var got string
	var err error
	for _, fixture := range smt.GetData() {
		fixture.T = smt.GetT()
		if fixture.Label == "" {
			fixture.Label = fixture.In
		}
		t := fixture.T
		for name, out := range fixture.Out {
			fixture.Name = name
			fixture.Obj, err = smt.MakeNewObject(fixture)

			//fmt.Printf("Testing '%s'.%s\n", fixture.Label, name)

			if fixture.Obj == nil {
				t.Skipf("no object created; test '%s' for '%s' skipped", name, fixture.Label)
				continue
			}
			if err == nil {
				got = smt.GetOutput(fixture)
			}

			if NoFailWrongOutput(t, fixture.Label, name, fixture.In, got, out, err) {
				continue
			}

			if FailWrongError(t, fixture.Label, name, fixture.In, out, err) {
				continue
			}

			if err != nil {
				t.Errorf("error for '%s' test '%s' w/input '%s': %s", fixture.Label, name, fixture.In, err.Error())
			}
		}
	}
}

//func HandleError(f *Fixture, err error) {
//	if err != nil {
//		f.T.Errorf("testing %s for %s failed: %s",
//			f.Name,
//			f.Label,
//			err.Error(),
//		)
//	}
//}
