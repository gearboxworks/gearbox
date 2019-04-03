package routes

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"github.com/labstack/echo"
	"net/http"
)

const AuthorityResourceVar api.ResourceVarName = "authority"
const StackNameResourceVar api.ResourceVarName = "stack"

type NamedStackIdRequest struct {
	NamedStackId types.StackId `json:"stack"`
}

func GetNamedStackIdFromRequest(rc *api.RequestContext) (nsid types.StackId, sts status.Status) {
	for range only.Once {
		if rc.Context.Request().Method == echo.GET {
			nsid = types.StackId(rc.Param(StackNameResourceVar))
			a := types.Stackname(rc.Param(AuthorityResourceVar))
			if a != "" {
				nsid = types.StackId(fmt.Sprintf("%s/%s", a, nsid))
			}
			break
		}
		snr := NamedStackIdRequest{}
		sts := rc.UnmarshalFromRequest(&snr)
		if status.IsError(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("invalid request format for '%s' routes", rc.RouteName),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    api.GetApiHelp("rc.RouteName", "correct request format"),
			})
			break
		}
		nsid = snr.NamedStackId
	}
	if nsid == "" {
		sts = status.Fail(&status.Args{
			Message:    "named stack ID is empty",
			Help:       api.GetApiHelp(rc.RouteName),
			HttpStatus: http.StatusNotFound,
		})
	}
	return nsid, sts
}
