package service

import (
	"encoding/json"
	"gearbox/only"
)

type ServicerProxy struct {
	Servicer
}

func ProxyServicer(servicer Servicer) *ServicerProxy {
	psp := &ServicerProxy{}
	psp.Servicer = servicer
	return psp
}

func (me *ServicerProxy) UnmarshalJSON(b []byte) (err error) {
	for range only.Once {
		data := make(map[string]json.RawMessage, 0)
		err = json.Unmarshal(b, &data)
		if err != nil {
			break
		}
	}
	return err
}
