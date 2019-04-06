package ja

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"reflect"
)

const SchemaVersion = "1.0"

type ResponseType string

const (
	DatasetResponse    ResponseType = "dataset"
	CollectionResponse ResponseType = "collection"
)

type RootDocument struct {
	ResponseType ResponseType      `json:"-"`
	JsonApi      *JsonApi          `json:"jsonapi,omitempty"`
	MetaMap      MetaMap           `json:"meta,omitempty"`
	LinkMap      LinkMap           `json:"links,omitempty"`
	Data         ResourceContainer `json:"data,omitempty"`
	Included     ResourceObjects   `json:"included,omitempty"`
	Errors       Errors            `json:"errors,omitempty"`
}

type RootDocArgs RootDocument

func NewRootDocument(responseType ResponseType, args ...*RootDocArgs) *RootDocument {
	var _args *RootDocArgs
	if len(args) == 0 {
		_args = &RootDocArgs{}
	} else {
		_args = args[0]
	}
	rd := RootDocument{}
	rd = RootDocument(*_args)
	rd.ResponseType = responseType
	if rd.Data == nil {
		switch responseType {
		case DatasetResponse:
			rd.Data = &ResourceObject{}
		case CollectionResponse:
			rd.Data = make(ResourceObjects, 0)
		default:
			panic(fmt.Sprintf("invalid response type '%s'", responseType))
		}
	}
	if rd.JsonApi == nil {
		rd.JsonApi = &JsonApi{}
	}
	rd.JsonApi.Version = SchemaVersion

	if rd.MetaMap == nil {
		rd.MetaMap = make(MetaMap, 0)
	}
	rd.MetaMap[MetaDcFormat] = ContentType
	rd.MetaMap[MetaDcType] = responseType

	if rd.Included == nil {
		rd.Included = make(ResourceObjects, 0)
	}

	if rd.LinkMap == nil {
		rd.LinkMap = make(LinkMap, 0)
	}

	if rd.Errors != nil {
		rd.Data = nil
	}
	return &rd
}

func (me *RootDocument) AddResourceObject(ro *ResourceObject) (sts status.Status) {
	for range only.Once {
		appender, ok := me.Data.(ResourceObjectAppender)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("cannot add resource objects to a '%s' response type", me.ResponseType),
			})
			break
		}
		me.Data, sts = appender.AppendResourceObject(ro)
	}
	return sts
}

//func (me *RootDocument) SetItemData(data interface{}) (sts status.Status) {
//	return sts
//}

//func (me *RootDocument) SetId(id ResourceId) (sts status.Status) {
//	setter, ok := me.Data.(ResourceIdSetter)
//	if ok {
//		sts = setter.SetId(id)
//	}
//	return sts
//}
//
//func (me *RootDocument) SetIds(ids ResourceIds) (sts status.Status) {
//	setter, ok := me.Data.(ResourceIdsSetter)
//	if ok {
//		sts = setter.SetIds(ids)
//	}
//	return sts
//}
//
//func (me *RootDocument) SetTypes(types ResourceTypes) (sts status.Status) {
//	setter, ok := me.Data.(ResourceTypesSetter)
//	if ok {
//		sts = setter.SetTypes(types)
//	}
//	return sts
//}
//
//func (me *RootDocument) SetType(_typ ResourceType) (sts status.Status) {
//	setter, ok := me.Data.(ResourceTypeSetter)
//	if ok {
//		sts = setter.SetType(_typ)
//	}
//	return sts
//}
func (me *RootDocument) SetMeta(meta MetaMap) (sts status.Status) {
	me.MetaMap = meta
	return nil
}
func (me *RootDocument) SetIncluded(included ResourceObjects) (sts status.Status) {
	me.Included = included
	return nil
}
func (me *RootDocument) SetLinks(linkmap LinkMap) (sts status.Status) {
	me.LinkMap = linkmap
	return nil
}
func (me *RootDocument) SetError(sts status.Status) status.Status {
	me.Errors = Errors{sts}
	me.Data = nil
	return nil
}
func (me *RootDocument) SetAttributes(attrs interface{}) (sts status.Status) {
	setter, ok := me.Data.(AttributesSetter)
	if ok {
		switch reflect.TypeOf(attrs).Kind() {
		case reflect.Slice, reflect.Array:
			me.Data = make(ResourceObjects, 0)

		case reflect.Struct:
			me.Data = &ResourceObject{}
		}
		sts = setter.SetAttributes(attrs)
	}
	return sts
}