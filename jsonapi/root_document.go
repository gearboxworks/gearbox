package jsonapi

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"github.com/apcera/util/uuid"
	"reflect"
	"strings"
)

const SchemaVersion = "1.0"
const ResponseTypeKey = "response_type"

type ResponseType string

const (
	DatasetResponse    ResponseType = "dataset"
	CollectionResponse ResponseType = "collection"
)

var _ apimodeler.ItemModeler = (*IncludedItem)(nil)

type IncludedList []*IncludedItem
type IncludedItem ResourceObject

func (me *IncludedItem) GetAttributeMap() AttributeMap {
	return me.AttributeMap
}

func (me *IncludedItem) GetRelatedItems(ctx *apimodeler.Context) (list apimodeler.List, sts status.Status) {
	panic("implement me")
}

func (me *IncludedItem) GetId() apimodeler.ItemId {
	panic("implement me")
}
func (me *IncludedItem) SetId(apimodeler.ItemId) status.Status {
	panic("implement me")
}
func (me *IncludedItem) GetType() apimodeler.ItemType {
	panic("implement me")
}
func (me *IncludedItem) GetItem() (apimodeler.ItemModeler, status.Status) {
	panic("implement me")
}
func (me *IncludedItem) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	panic("implement me")
}

func (me IncludedList) AppendItem(item apimodeler.ItemModeler) (inc IncludedList, sts status.Status) {
	inc = me
	for range only.Once {
		ii, ok := item.(*IncludedItem)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("item '%s' does not implement ja.ResourceObject", item.GetId()),
			})
			break
		}
		inc = append(me, ii)
	}
	return inc, sts
}

var _ apimodeler.RootDocumentor = (*RootDocument)(nil)

type RootDocument struct {
	ResponseType types.ResponseType `json:"-"`
	JsonApi      *JsonApi           `json:"jsonapi,omitempty"`
	MetaMap      MetaMap            `json:"meta,omitempty"`
	LinkMap      apimodeler.LinkMap `json:"links,omitempty"`
	Data         ResourceContainer  `json:"data,omitempty"`
	Included     IncludedList       `json:"included,omitempty"`
	Errors       ErrorObjects       `json:"errors,omitempty"`
}

func (me *RootDocument) GetDataRelationshipsLinkMap() (rlm apimodeler.LinkMap, sts status.Status) {
	for range only.Once {
		getter, ok := me.Data.(RelationshipsLinkMapGetter)
		if !ok {
			rlm = apimodeler.LinkMap{}
		}
		rlm, sts = getter.GetRelationshipsLinkMap()
	}
	return rlm, sts
}

type RootDocArgs RootDocument

func NewRootDocument(ctx Contexter, responseType types.ResponseType, args ...*RootDocArgs) *RootDocument {
	ctx.Set(ResponseTypeKey, responseType)
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
		case global.ItemResponse:
			rd.Data = &ResourceObject{}
		case global.ListResponse:
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
	rd.MetaMap[apimodeler.MetaDcFormat] = ContentType
	rd.MetaMap[apimodeler.MetaDcType] = responseType

	if rd.Included == nil {
		rd.Included = make(IncludedList, 0)
	}

	if rd.LinkMap == nil {
		rd.LinkMap = make(apimodeler.LinkMap, 0)
	}

	if rd.Errors != nil {
		rd.Data = nil
	}
	return &rd
}

func (me *RootDocument) AddMeta(name apimodeler.Metaname, value apimodeler.MetaValue) {
	me.MetaMap[name] = value
}

func (me *RootDocument) AddLink(rel apimodeler.RelType, link apimodeler.LinkImplementor) {
	me.LinkMap[rel] = link
}

func (me *RootDocument) AddLinks(links apimodeler.LinkMap) {
	for rel, link := range links {
		me.AddLink(rel, link)
	}
}

func (me *RootDocument) SetRelated(ctx *apimodeler.Context, list apimodeler.List) (sts status.Status) {
	for range only.Once {
		inc := make(IncludedList, 0, len(list))
		for _, item := range list {
			inc, sts = inc.AppendItem(item)
			if is.Error(sts) {
				break
			}
		}
		me.Included = inc
	}
	return sts
}

func (me *RootDocument) GetResponseType() types.ResponseType {
	return me.ResponseType
}

func (me *RootDocument) GetRootDocument() interface{} {
	return me
}

func (me *RootDocument) AddResponseItem(item apimodeler.ItemModeler) (sts status.Status) {
	for range only.Once {
		appender, ok := me.Data.(ResourceObjectAppender)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("cannot add resource objects to a '%s' response type", me.ResponseType),
			})
			break
		}
		ro, ok := item.(*ResourceObject)
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

func (me *RootDocument) SetLinks(linkmap apimodeler.LinkMap) (sts status.Status) {
	me.LinkMap = linkmap
	return nil
}
func (me *RootDocument) SetErrors(err error) {
	sts, ok := err.(status.Status)
	if ok {
		err = sts.GetFullError()
	}
	if err == nil {
		err = fmt.Errorf("unexpected error")
	}
	if sts == nil {
		sts = status.OurBad("unexpected error")
	}
	msg := err.Error()
	me.Errors = ErrorObjects{NewErrorObject(&ErrorObjectArgs{
		ErrorId:    ErrorId(uuid.Generate().String()),
		ErrorCode:  ErrorCode(fmt.Sprintf("%s-error", me.ResponseType)),
		Title:      msg[:strings.IndexByte(msg+";", ';')],
		Detail:     err.Error(),
		HttpStatus: HttpStatus(sts.HttpStatus()),
		ErrorSource: NewErrorSource(&ErrorSourceArgs{
			JsonPointer:  Link("@TODO: json pointer goes here"),
			UrlParameter: UrlParameter("@TODO: url parameter goes here"),
		}),
		LinkMap: LinkMap{
			AboutRelType: Link("@TODO: about error link goes here"),
			TypeRelType:  Links{"@TODO: array error type links go here"},
		},
	})}
	me.Data = nil
	return
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
