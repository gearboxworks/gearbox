package is

import "gearbox/status"

type (
	Status = status.Status
)

func Error(sts Status) bool {
	return sts != nil && status.IsError(sts)
}
func Success(sts Status) bool {
	return sts == nil || status.IsSuccess(sts)
}
