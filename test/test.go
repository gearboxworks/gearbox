package test

import (
	"gearbox/stat"
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

func NoFailWrongOutput(t *testing.T, label string, name string, input string, got string, want Args, status stat.Status) bool {
	if !want.Fail && !status.IsError() {
		if got != want.Want {
			t.Errorf("nofail: %s.%s, want: %s, got: %s, input: %s", label, name, want.Want, got, input)
		}
		return true
	}
	return false
}

func FailWrongError(t *testing.T, label string, name string, input string, want Args, status stat.Status) bool {
	if want.Fail {
		if status.IsError() && status.Message != want.Want {
			t.Errorf("fail: %s.%s, want: %s, got: %s, input: %s", label, name, want.Want, status.Message, input)
		}
		return true
	}
	return false

}

type StructMethodTester interface {
	MakeNewObject(f *Fixture) (obj interface{}, status stat.Status)
	GetOutput(f *Fixture) string
	GetT() *testing.T
	GetData() Table
}

func StructMethodsTest(smt StructMethodTester) {
	var got string
	var status stat.Status
	for _, fixture := range smt.GetData() {
		fixture.T = smt.GetT()
		if fixture.Label == "" {
			fixture.Label = fixture.In
		}
		t := fixture.T
		for name, out := range fixture.Out {
			fixture.Name = name
			fixture.Obj, status = smt.MakeNewObject(fixture)

			//fmt.Printf("Testing '%s'.%s\n", fixture.Label, name)

			if fixture.Obj == nil {
				t.Skipf("no object created; test '%s' for '%s' skipped", name, fixture.Label)
				continue
			}
			if status.IsError() {
				got = smt.GetOutput(fixture)
			}

			if NoFailWrongOutput(t, fixture.Label, name, fixture.In, got, out, status) {
				continue
			}

			if FailWrongError(t, fixture.Label, name, fixture.In, out, status) {
				continue
			}

			if status.IsError() {
				t.Errorf("error for '%s' test '%s' w/input '%s': %s",
					fixture.Label,
					name,
					fixture.In,
					status.Message,
				)
			}
		}
	}
}

//func HandleError(f *Fixture, err error) {
//	if err != nil {
//		f.T.Errorf("testing %s for %s failed: %s",
//			f.Type,
//			f.Label,
//			err.Error(),
//		)
//	}
//}
