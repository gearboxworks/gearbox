package integration

import (
	"gearbox"
	"gearbox/api"
	"gearbox/only"
	"gearbox/test"
	"testing"
)

func TestHostApiStackDetails(t *testing.T) {
	var url api.UriTemplate
	for range only.Once {
		gb, status := GearboxConstructAndInitialize(t)
		if status.IsError() {
			t.Error(status.Message)
		}
		hapi := gb.GetHostApi()

		uritemplate := hapi.Api.Defaults.Links[gearbox.StackDetailsResource]
		valuesFunc := gb.GetHostApi().Api.ValuesFuncMap[gearbox.StackDetailsResource]
		sns, status := valuesFunc()
		if status.IsError() {
			break
		}
		for _, sn := range sns {
			url = api.ExpandUriTemplate(uritemplate, api.UriTemplateVars{
				gearbox.StackNameResourceVar: sn,
			})
			hc := &test.HttpClient{}
			status = hc.GET(string(url))
			if status.IsError() {
				t.Error(status.Message)
				break
			}
			var body []byte
			body, status = hc.GetBody()
			if status.IsError() {
				t.Error(status.Message)
				break
			}
			println(string(body))
		}
	}
	return

}

//func TestHostApiResponses(t *testing.T) {
//	gb, status := GearboxConstructAndInitialize(t)
//	if status.IsError() {
//		t.Error(status.Message)
//	}
//	hapi := gb.GetHostApi()
//	for m, mm := range hapi.Api.MethodMap {
//		if m != http.MethodGet {
//			continue
//		}
//		t.Run(string(m), func(t *testing.T) {
//			for n, _ := range mm {
//				t.Run(string(n), func(t *testing.T) {
//					for range only.Once {
//						url, status := hapi.Api.GetUrl(n, api.UriTemplateVars{})
//						if status.IsError() {
//							t.Error(status.Message)
//							break
//						}
//						hc := &test.HttpClient{}
//						status = hc.GET(string(url))
//						if status.IsError() {
//							t.Error(status.Message)
//							break
//						}
//						var body []byte
//						body,status = hc.GetBody()
//						if status.IsError() {
//							t.Error(status.Message)
//							break
//						}
//						println(string(body))
//					}
//				})
//			}
//		})
//	}
//}
