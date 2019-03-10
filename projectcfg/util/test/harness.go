package test

//import "testing"
//
//const ErrExpected = "~"
//
//type InOut interface {
//	Input() string
//	Output() string
//}
//
//type Harness struct {
//	T *testing.T
//	InOut
//	Item interface{}
//}
//
//func NewHarness(t *testing.T, td InOut, i interface{}) *Harness {
//	return &Harness{
//		T:     t,
//		InOut: td,
//		Item:  i,
//	}
//}
//func (th *Harness) Run(fn func()) {
//	th.T.Run(th.Input(), func(t *testing.T) {
//		fn()
//	})
//}
