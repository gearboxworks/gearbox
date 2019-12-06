package network

import (
	"gearbox/eventbroker/msgs"
	"net/url"
)

type ServiceConfig struct {
	EntityId     msgs.Address
	EntityName   msgs.Address
	EntityParent *msgs.Address
	UrlString    string   `json:"urlstring"` //
	Url          *url.URL `json:"url"`       //

	Name   Name   `json:"name"`   // == Service.Entry.Instance
	Type   Type   `json:"type"`   // == Service.Entry.Service
	Domain Domain `json:"domain"` // == Service.Entry.Domain
	Port   Port   `json:"port"`   // == Service.Entry.Port
	Text   Text   `json:"text"`   // == Service.Entry.Text
	TTL    uint32 `json:"ttl"`    // == Service.Entry.TTL
}
