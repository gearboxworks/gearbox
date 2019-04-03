package util

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
)

type FilepathHelpUrlGetter interface {
	FilepathGetter
	HelpUrlGetter
}

type FilepathGetter interface {
	GetFilepath() types.AbsoluteFilepath
}

type HelpUrlGetter interface {
	GetHelpUrl() string
}

func UnmarshalJson(j []byte, obj FilepathHelpUrlGetter) (sts status.Status) {
	for range only.Once {
		err := json.Unmarshal(j, &obj)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("failed to unmarshal JSON for '%s'", obj.GetFilepath()),
				Help: fmt.Sprintf("ensure '%s' is in correct format per %s",
					obj.GetFilepath(),
					obj.GetHelpUrl(), // @TODO Improve the accuracy of this help once we have docs online
				),
			})
			break
		}
	}
	return sts
}
