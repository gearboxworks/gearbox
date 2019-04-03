package gearbox

import (
	"encoding/json"
	"gearbox/api"
)

type ApiBaseUrls struct {
	HostUrl api.UriTemplate `json:"host_url"`
	VmUrl   api.UriTemplate `json:"vm_url"`
}

func NewApiBaseUrls(hostUrl, vmUrl api.UriTemplate) *ApiBaseUrls {
	return &ApiBaseUrls{
		HostUrl: hostUrl,
		VmUrl:   vmUrl,
	}
}
func (me *ApiBaseUrls) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}
