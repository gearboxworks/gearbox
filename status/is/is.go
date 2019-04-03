package is

import "gearbox/status"

func Error(sts status.Status) bool {
	return sts != nil && status.IsError(sts)
}
func Success(sts status.Status) bool {
	return sts == nil || status.IsSuccess(sts)
}
