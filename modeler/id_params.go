package modeler

import (
	"gearbox/types"
	"strings"
)

type IdParams []IdParam
type IdParam string

func (me IdParams) Slugify() types.UrlTemplate {
	ss := make([]string, len(me))
	for i, idp := range me {
		ss[i] = ":" + string(idp)
	}
	return types.UrlTemplate(strings.Join(ss, "/"))
}
