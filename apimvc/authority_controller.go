package apimvc

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
	"sort"
)

const AuthorityControllerName types.RouteName = "authorities"
const AuthoritiesBasepath types.Basepath = "/authorities"

var NilAuthorityController = (*AuthorityController)(nil)
var _ ListController = NilAuthorityController

type AuthorityController struct {
	Controller
	Gearbox gearbox.Gearboxer
}

func NewAuthorityController(gb gearbox.Gearboxer) *AuthorityController {
	return &AuthorityController{
		Gearbox: gb,
	}
}

func (me *AuthorityController) GetRelatedFields() RelatedFields {
	return RelatedFields{}
}

func (me *AuthorityController) CanAddItem(*Context) bool {
	return false
}

func (me *AuthorityController) GetName() types.RouteName {
	return AuthorityControllerName
}

func (me *AuthorityController) GetListLinkMap(*Context, ...FilterPath) (lm LinkMap, sts Status) {
	return LinkMap{
		//RelatedRelType: Link("http://example.org"),
	}, sts
}

func (me *AuthorityController) GetBasepath() types.Basepath {
	return AuthoritiesBasepath
}

func (me *AuthorityController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *AuthorityController) GetIdParams() IdParams {
	return IdParams{
		AuthorityIdParam,
		StacknameIdParam,
		RoleIdParam,
	}
}

func (me *AuthorityController) GetList(ctx *Context, filterPath ...FilterPath) (list List, sts Status) {
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

func (me *AuthorityController) FilterList(ctx *Context, filterPath FilterPath) (list List, sts Status) {
	return me.GetList(ctx, filterPath)
}

func (me *AuthorityController) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts Status) {
	for range only.Once {
		if len(filterPath) == 0 {
			filterPath = []FilterPath{NoFilterPath}
		}
		list, sts := me.GetList(ctx, filterPath[0])
		if is.Error(sts) {
			break
		}
		itemids = make(ItemIds, len(list))
		i := 0
		for _, item := range list {
			itemids[i] = ItemId(item.GetId())
			i++
		}
	}
	return itemids, sts
}

func (me *AuthorityController) GetItem(ctx *Context, authorityid ItemId) (list ItemModeler, sts Status) {
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

func (me *AuthorityController) GetItemDetails(ctx *Context, itemid ItemId) (ItemModeler, Status) {
	return me.GetItem(ctx, itemid)
}

func (me *AuthorityController) FilterItem(in ItemModeler, filterPath FilterPath) (out ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *AuthorityController) GetFilterMap() FilterMap {
	return FilterMap{}
}

func assertAuthority(item ItemModeler) (s *AuthorityModel, sts Status) {
	s, ok := item.(*AuthorityModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a AuthorityModel: %v", item),
		})
	}
	return s, sts
}
