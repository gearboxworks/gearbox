package api

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"gearbox/status"
)

type ResponseMeta struct {
	Version   string    `json:"version"`
	Service   string    `json:"service"`
	DocsUrl   string    `json:"docs_url"`
	RouteName RouteName `json:"route"`
}

type Response struct {
	Success    bool         `json:"success"`
	StatusCode int          `json:"status_code"`
	Meta       ResponseMeta `json:"meta,omitempty"`
	Links      Links        `json:"links,omitempty"`
	Data       interface{}  `json:"data,omitempty"`
}

func (me *Response) GetUrlPathTemplate(resourceType RouteName) (url UriTemplate, sts status.Status) {
	for range only.Once {
		var ok bool
		url, ok = me.Links[resourceType].(UriTemplate)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("no '%s' in resource links", resourceType),
				Help:    ContactSupportHelp(),
			})
		}
	}
	return url, sts
}

func (me *Response) Clone() *Response {
	r := Response{}
	for range only.Once {
		b, err := json.Marshal(me)
		if err != nil {
			break
		}
		_ = json.Unmarshal(b, &r)
	}
	return &r
}
