package gearbox

import (
	"encoding/json"
	"gearbox/types"
)

type ApiBaseUrls struct {
	HostUrl types.UrlTemplate `json:"host_url"`
	VmUrl   types.UrlTemplate `json:"vm_url"`
}

func NewApiBaseUrls(hostUrl, vmUrl types.UrlTemplate) *ApiBaseUrls {
	return &ApiBaseUrls{
		HostUrl: hostUrl,
		VmUrl:   vmUrl,
	}
}
func (me *ApiBaseUrls) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}
