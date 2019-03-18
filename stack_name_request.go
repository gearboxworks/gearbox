package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"github.com/labstack/echo"
	"net/http"
)

type StackNameRequest struct {
	StackName StackName `json:"stack"`
}

func getStackName(rc *api.RequestContext) (sn StackName, status stat.Status) {
	for range only.Once {
		if rc.Context.Request().Method == echo.GET {
			sn = StackName(rc.Param(StackNameResourceVar))
			a := StackName(rc.Param(AuthorityResourceVar))
			if a != "" {
				sn = StackName(fmt.Sprintf("%s/%s", a, sn))
			}
			break
		}
		snr := StackNameRequest{}
		status := rc.UnmarshalFromRequest(&snr)
		if status.IsError() {
			status.PriorStatus = status.String()
			status.Message = fmt.Sprintf("invalid request format for '%s' resource", rc.ResourceName)
			status.HttpStatus = http.StatusBadRequest
			status.ApiHelp = api.GetApiHelp("rc.ResourceName", "correct request format")
			break
		}
		sn = snr.StackName
	}
	if sn == "" {
		status = stat.NewStatus(&stat.Args{
			Message:    "stack name is empty",
			Help:       api.GetApiHelp(rc.ResourceName),
			HttpStatus: http.StatusNotFound,
			Error:      stat.IsStatusError,
		})
	}
	return sn, status
}
