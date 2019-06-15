package only

//
// This used for 1-time for loops to provide break-sequences for error handling, e.g.
//
//   func DoSeveralThings() (whatever, status.Status) {
//		var result whatever
//		for range only.Once {
//			sts := DoSomething
//			if is.Error(sts) {
//				break
//			}
//			result = DoSomethingElse()
//		}
//      return result,sts
//   }
//
var Once string = "1"
var Twice = []string{"", ""}
