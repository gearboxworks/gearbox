package jsonapi

import (
	"encoding/json"
	"fmt"
	"gearbox/apiworks"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"github.com/clbanning/checkjson"
	"net/http"
)

var NilResourceObjects = (ResourceObjects)(nil)
var _ ResourceContainer = NilResourceObjects
var _ RelationshipsLinkMapGetter = NilResourceObjects

type ResourceObjects []*ResourceObject

func (me ResourceObjects) GetRelationshipsLinkMap() (lm apiworks.LinkMap, sts status.Status) {
	lm = make(apiworks.LinkMap, 0)
	for fn, _lm := range me {
		for rel, link := range _lm.LinkMap {
			if rel != apiworks.SelfRelType {
				continue
			}
			lm[apiworks.RelType(fn)] = link
		}
	}
	return lm, sts
}

func (ResourceObjects) ContainsResource() {}

func (me ResourceObjects) SetAttributes(attrs interface{}) (sts status.Status) {
	panic("Not yet implemented")
	return nil
}

func (me ResourceObjects) AppendResourceObject(ro *ResourceObject) (ResourceObjects, status.Status) {
	return append(me, ro), nil
}

func (me ResourceObjects) SetIds(ids ResourceIds) (sts status.Status) {
	for i, ro := range me {
		sts = ro.SetId(apiworks.ItemId(ids[i]))
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func (me ResourceObjects) SetTypes(types ResourceTypes) (sts status.Status) {
	for i, ro := range me {
		sts = ro.SetType(apiworks.ItemType(types[i]))
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func (me *ResourceObject) Unmarshal(data []byte) (sts status.Status) {
	var err error
	var msg string
	//var mks []string
	for range only.Once {
		err = json.Unmarshal(data, &me)
		if err != nil {
			err = checkjson.ResolveJSONError(data, err)
			if err != nil {
				msg = fmt.Sprintf("unable to unmarshal JSON to type '%T'", me)
			}
			break
		}
		//mks, err = checkjson.MissingJSONKeys(data,me)
		//if err != nil {
		//	msg = err.Error()
		//	break
		//}
		//var obj map[string]json.RawMessage
		//err = json.Unmarshal(data, &obj)
		//if err != nil {
		//	err = checkjson.ResolveJSONError(data, err)
		//	if err != nil {
		//		msg = err.Error()
		//		break
		//	}
		//	msg = "unable to unmarshal JSON to map of json.RawMessage keyed by string"
		//	break
		//}
		//attrs,ok := obj["attributes"]
		//if !ok {
		//	break
		//}
		//data, err = json.Marshal(attrs)
		//if err != nil {
		//	err = checkjson.ResolveJSONError(data, err)
		//	if err != nil {
		//		msg = err.Error()
		//		break
		//	}
		//	msg = "unable to marshal properties of 'attributes' property"
		//	break
		//}
		//mks2, err := checkjson.MissingJSONKeys(attrs, me.AttributeMap)
		//if err != nil {
		//	msg = err.Error()
		//	break
		//}
		//if len(mks2) == 0 {
		//	break
		//}
		//mks = append(mks, mks2...)
	}

	for range only.Once {
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				HttpStatus: http.StatusUnprocessableEntity,
				Message:    msg,
			})
			break
		}
		//if len(mks) > 0 {
		//	sts = status.Fail(&status.Args{
		//		HttpStatus: http.StatusBadRequest,
		//		Message: fmt.Sprintf("missing keys in JSON: '%s'",strings.Join(mks,", ")),
		//	})
		//	break
		//}
	}
	return sts
}
