package gearbox

import "encoding/json"

type ApiBaseUrls struct {
	HostUrl string `json:"host_url"`
	VmUrl   string `json:"vm_url"`
}

func NewApiBaseUrls(hostUrl, vmUrl string) *ApiBaseUrls {
	return &ApiBaseUrls{
		HostUrl: hostUrl,
		VmUrl:   vmUrl,
	}
}
func (me *ApiBaseUrls) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}
