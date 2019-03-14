package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/only"
	"gearbox/util"
)

type ResourceType string

type ResponseMeta struct {
	Version  string       `json:"version"`
	Service  string       `json:"service"`
	DocsUrl  string       `json:"docs_url"`
	Resource ResourceName `json:"resource"`
}

type Response struct {
	Success    bool         `json:"success"`
	StatusCode int          `json:"status_code"`
	Meta       ResponseMeta `json:"meta"`
	Links      Links        `json:"links"`
	Data       interface{}  `json:"data,omitempty"`
}

func (me *Response) GetApiSelfLink(resourceType ResourceName) (url string, err error) {
	for range only.Once {
		var ok bool
		url, ok = me.Links[resourceType]
		if !ok {
			err = util.AddHelpToError(
				errors.New(fmt.Sprintf("no '%s' in resource links", resourceType)),
				ContactSupportHelp(),
			)
		}
	}
	return url, err
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