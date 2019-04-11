package apimodels

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
	"sort"
)

const AuthoritiesName types.RouteName = "authorities"
const AuthoritysBasepath types.Basepath = "/authorities"

var NilAuthorityModel = (*AuthorityModel)(nil)
var _ apimodeler.ApiModeler = NilAuthorityModel

type AuthorityModel struct {
	apimodeler.BaseModel
	Gearbox gearbox.Gearboxer
}

func NewAuthorityModel(gb gearbox.Gearboxer) *AuthorityModel {
	return &AuthorityModel{
		Gearbox: gb,
	}
}

func (me *AuthorityModel) CanAddItem(*apimodeler.Context) bool {
	return false
}

func (me *AuthorityModel) GetName() types.RouteName {
	return AuthoritiesName
}

func (me *AuthorityModel) GetListLinkMap(*apimodeler.Context, ...apimodeler.FilterPath) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("http://example.org"),
	}, sts
}

func (me *AuthorityModel) GetBasepath() types.Basepath {
	return AuthoritysBasepath
}

func (me *AuthorityModel) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *AuthorityModel) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{
		AuthorityIdParam,
		StacknameIdParam,
		RoleIdParam,
	}
}

func (me *AuthorityModel) GetList(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	for range only.Once {
		gbgas, sts := me.Gearbox.GetGears().GetAuthorities()
		if is.Error(sts) {
			break
		}
		for _, gbga := range gbgas {
			ns, sts := NewFromGearsAuthority(ctx, gbga)
			if is.Error(sts) {
				break
			}
			list = append(list, ns)
		}
		sort.Slice(list, func(i, j int) bool {
			return list[i].GetId() < list[j].GetId()
		})
	}
	return list, sts
}

func (me *AuthorityModel) FilterList(ctx *apimodeler.Context, filterPath apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	return me.GetList(ctx, filterPath)
}

func (me *AuthorityModel) GetListIds(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (itemids apimodeler.ItemIds, sts status.Status) {
	for range only.Once {
		if len(filterPath) == 0 {
			filterPath = []apimodeler.FilterPath{apimodeler.NoFilterPath}
		}
		list, sts := me.GetList(ctx, filterPath[0])
		if is.Error(sts) {
			break
		}
		itemids = make(apimodeler.ItemIds, len(list))
		i := 0
		for _, item := range list {
			itemids[i] = apimodeler.ItemId(item.GetId())
			i++
		}
	}
	return itemids, sts
}

func (me *AuthorityModel) GetItem(ctx *apimodeler.Context, authorityid apimodeler.ItemId) (list apimodeler.ApiItemer, sts status.Status) {
	var ns *Authority
	for range only.Once {
		gbgs, sts := me.Gearbox.GetGears().FindAuthority(types.AuthorityDomain(authorityid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("Authority '%s' not found", authorityid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewFromGearsAuthority(ctx, gbgs)
		if is.Error(sts) {
			break
		}
		sts = status.Success("Authority '%s' found", authorityid)
	}
	return ns, sts
}

func (me *AuthorityModel) GetItemDetails(ctx *apimodeler.Context, itemid apimodeler.ItemId) (apimodeler.ApiItemer, status.Status) {
	return me.GetItem(ctx, itemid)
}

func (me *AuthorityModel) FilterItem(in apimodeler.ApiItemer, filterPath apimodeler.FilterPath) (out apimodeler.ApiItemer, sts status.Status) {
	out = in
	return out, sts
}

func (me *AuthorityModel) GetFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{}
}

func assertAuthority(item apimodeler.ApiItemer) (s *Authority, sts status.Status) {
	s, ok := item.(*Authority)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a Authority: %v", item),
		})
	}
	return s, sts
}
