package apimvc

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

var NilAuthorityController = (*AuthorityController)(nil)
var _ apimodeler.ApiController = NilAuthorityController

type AuthorityController struct {
	apimodeler.Controller
	Gearbox gearbox.Gearboxer
}

func NewAuthorityController(gb gearbox.Gearboxer) *AuthorityController {
	return &AuthorityController{
		Gearbox: gb,
	}
}

func (me *AuthorityController) CanAddItem(*apimodeler.Context) bool {
	return false
}

func (me *AuthorityController) GetName() types.RouteName {
	return AuthoritiesName
}

func (me *AuthorityController) GetListLinkMap(*apimodeler.Context, ...apimodeler.FilterPath) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("http://example.org"),
	}, sts
}

func (me *AuthorityController) GetBasepath() types.Basepath {
	return AuthoritysBasepath
}

func (me *AuthorityController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *AuthorityController) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{
		AuthorityIdParam,
		StacknameIdParam,
		RoleIdParam,
	}
}

func (me *AuthorityController) GetList(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
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

func (me *AuthorityController) FilterList(ctx *apimodeler.Context, filterPath apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	return me.GetList(ctx, filterPath)
}

func (me *AuthorityController) GetListIds(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (itemids apimodeler.ItemIds, sts status.Status) {
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

func (me *AuthorityController) GetItem(ctx *apimodeler.Context, authorityid apimodeler.ItemId) (list apimodeler.Itemer, sts status.Status) {
	var ns *AuthorityModel
	for range only.Once {
		gbgs, sts := me.Gearbox.GetGears().FindAuthority(types.AuthorityDomain(authorityid))
		if is.Error(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    fmt.Sprintf("AuthorityModel '%s' not found", authorityid),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		ns, sts = NewFromGearsAuthority(ctx, gbgs)
		if is.Error(sts) {
			break
		}
		sts = status.Success("AuthorityModel '%s' found", authorityid)
	}
	return ns, sts
}

func (me *AuthorityController) GetItemDetails(ctx *apimodeler.Context, itemid apimodeler.ItemId) (apimodeler.Itemer, status.Status) {
	return me.GetItem(ctx, itemid)
}

func (me *AuthorityController) FilterItem(in apimodeler.Itemer, filterPath apimodeler.FilterPath) (out apimodeler.Itemer, sts status.Status) {
	out = in
	return out, sts
}

func (me *AuthorityController) GetFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{}
}

func assertAuthority(item apimodeler.Itemer) (s *AuthorityModel, sts status.Status) {
	s, ok := item.(*AuthorityModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a AuthorityModel: %v", item),
		})
	}
	return s, sts
}
