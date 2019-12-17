package gearbox

import (
	"encoding/json"
	"gearbox/types"
)

type ApiBaseurls struct {
	HostUrl types.UrlTemplate `json:"host_url"`
	VmUrl   types.UrlTemplate `json:"vm_url"`
}

func NewApiBaseurls(hostUrl, vmUrl types.UrlTemplate) *ApiBaseurls {
	return &ApiBaseurls{
		HostUrl: hostUrl,
		VmUrl:   vmUrl,
	}
}
func (me *ApiBaseurls) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}
