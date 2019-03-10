package util

import "fmt"

type OxfordCommaArgs struct {
	SingleQuote bool
	Conjunction string
}

func OxfordComma(ss []string, args ...*OxfordCommaArgs) string {
	var ocs string
	var _args *OxfordCommaArgs
	if len(args)==0 {
		_args = &OxfordCommaArgs{}
	} else {
		_args = args[0]
	}
	penultimate := len(ss)-2
	switch penultimate {
	case -2:
		break
	case -1:
		ocs = ss[0]
	case 0:
		ocs = fmt.Sprintf("%s %s %s", ss[0], _args.Conjunction, ss[1])
	default:
		for i,s := range ss {
			if _args.SingleQuote && i == penultimate {
				ocs = fmt.Sprintf("%s, %s %s", ocs, _args.Conjunction, s)
			} else {
				ocs = fmt.Sprintf("%s, %s", ocs, s)
			}
		}
	}
	return ocs
}

