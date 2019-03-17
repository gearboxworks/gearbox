package integration

import (
	"testing"
)

func TestHostApiResponses(t *testing.T) {
	gb, status := GearboxConstructAndInitialize(t)
	if status.IsError() {
		t.Error(status.Message)
	}
	hapi := gb.GetHostApi()
	for m, eps := range hapi.Api.MethodMap {
		t.Run(string(m), func(t *testing.T) {
			for n, url := range eps {
				t.Run(string(n), func(t *testing.T) {
					println(url)
				})
			}
		})
	}
}
