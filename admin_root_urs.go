package gearbox

import "encoding/json"

type ApiRootUrls struct {
	HostUrl string `json:"host_url"`
	VmUrl   string `json:"vm_url"`
}

func NewApiRootUrls(hostUrl, vmUrl string) *ApiRootUrls {
	return &ApiRootUrls{
		HostUrl: hostUrl,
		VmUrl:   vmUrl,
	}
}
func (me *ApiRootUrls) Bytes() []byte {
	b, _ := json.Marshal(me)
	return b
}
