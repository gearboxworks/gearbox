package util

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
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

func UnmarshalJson(j []byte,obj FilepathHelpUrlGetter) (err error) {
	for range only.Once {
		err := json.Unmarshal(j, &obj)
		if err != nil {
			err = AddHelpToError(
				fmt.Errorf("unable to load config file '%s'", obj.GetFilepath()),
				fmt.Sprintf("ensure '%s' is in correct format per %s",
					obj.GetFilepath(),
					obj.GetHelpUrl(),
				),
			)
			break
		}
	}
	return err
}

