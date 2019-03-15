package util

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/stat"
)

type FilepathHelpUrlGetter interface {
	FilepathGetter
	HelpUrlGetter
}

type FilepathGetter interface {
	GetFilepath() string
}

type HelpUrlGetter interface {
	GetHelpUrl() string
}

func UnmarshalJson(j []byte, obj FilepathHelpUrlGetter) (status stat.Status) {
	for range only.Once {
		err := json.Unmarshal(j, &obj)
		if err != nil {
			status = stat.NewFailedStatus(&stat.Args{
				Error:   err,
				Message: fmt.Sprintf("failed to unmarshal JSON for '%s'", obj.GetFilepath()),
				Help: fmt.Sprintf("ensure '%s' is in correct format per %s",
					obj.GetFilepath(),
					obj.GetHelpUrl(), // @TODO Improve the accuracy of this help once we have docs online
				),
			})
			break
		}
	}
	return status
}
