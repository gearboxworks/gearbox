package apimvc

import (
	"encoding/json"
	"fmt"
	"gearbox/jsonapi"
	"gearbox/types"
	"gearbox/util"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
	"github.com/mitchellh/go-homedir"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const DirectoryControllerName types.RouteName = "directories"
const DirectoriesBasepath types.Basepath = "/directories"

var NilDirectoryController = (*DirectoryController)(nil)
var _ ListController = NilDirectoryController

type DirectoryController struct {
	Controller
}

func NewDirectoryController() *DirectoryController {
	return &DirectoryController{}
}

func (me *DirectoryController) GetNilItem(ctx *Context) ItemModeler {
	return NilDirectoryModel
}

func (me *DirectoryController) GetName() types.RouteName {
	return DirectoryControllerName
}

func (me *DirectoryController) GetBasepath() types.Basepath {
	return DirectoriesBasepath
}

func (me *DirectoryController) GetListIds(ctx *Context, filterPath ...FilterPath) (itemids ItemIds, sts Status) {
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

func (me *DirectoryController) GetItem(ctx *Context, dir ItemId) (item ItemModeler, sts Status) {
	for range only.Once {
		var d types.Dir
		d, sts = UnescapeDirectory(types.Dir(dir))
		if is.Error(sts) {
			break
		}
		item, sts = FindDirectory(d)
	}
	return item, sts
}

func (me *DirectoryController) FilterItem(in ItemModeler, filterPath FilterPath) (out ItemModeler, sts Status) {
	out = in
	return out, sts
}

func (me *DirectoryController) GetFilterMap() FilterMap {
	return GetDirectoryFilterMap()
}

func (me *DirectoryController) AddItem(ctx *Context, item ItemModeler) (im ItemModeler, sts Status) {
	for range only.Once {
		var dm *DirectoryModel
		dm, sts = me.getDirectoryModelFromItem(item)
		if is.Error(sts) {
			break
		}
		sts = CanAddDirectory(dm.Directory)
		if is.Error(sts) {
			break
		}
		sts = AddDirectory(dm.Directory)
		if is.Error(sts) {
			break
		}
		im = dm
		sts = status.Success("directory '%s' added", dm.Directory).
			SetHttpStatus(http.StatusCreated)
	}
	return im, sts
}

func (me *DirectoryController) DeleteItem(ctx *Context, itemid ItemId) (sts Status) {
	return status.Fail().
		SetHttpStatus(http.StatusMethodNotAllowed).
		SetMessage("cannot delete directory '%s'", itemid) // @TODO capture the dir to display
}

func (me *DirectoryController) UpdateItem(ctx *Context, item ItemModeler) (im ItemModeler, sts Status) {
	return im, status.Fail().
		SetHttpStatus(http.StatusMethodNotAllowed).
		SetMessage("cannot delete directory '%s'", item.GetId())
}

func (me *DirectoryController) maybeGetDirectoryDirectory(dm *DirectoryModel, itemid ItemId) types.Dir {
	dir := types.Dir(dm.GetId())
	for range only.Once {
		if dir != "" {
			break
		}
		if itemid != "" {
			dir = types.Dir(itemid)
			break
		}
	}
	return dir
}

func (me *DirectoryController) setDirectoryModelId(dm *DirectoryModel, itemid ItemId) (sts Status) {
	for range only.Once {
		dir := me.maybeGetDirectoryDirectory(dm, itemid)
		if dir == "" {
			break
		}
		sts = dm.SetId(ItemId(dir))
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func (me *DirectoryController) getDirectoryModelFromItem(item ItemModeler) (dm *DirectoryModel, sts Status) {
	for range only.Once {
		var ro *jsonapi.ResourceObject
		ro, sts = jsonapi.AssertResourceObject(item)
		if is.Error(sts) {
			break
		}
		if ro.GetType() != DirectoryModelType {
			sts = status.Fail().
				SetHttpStatus(http.StatusBadRequest).
				SetMessage("invalid request type '%s'; should be '%s'",
					ro.GetType(),
					DirectoryModelType,
				)
			break
		}
		var b []byte
		b, sts = ro.MarshalAttributeMap()
		if is.Error(sts) {
			break
		}
		err := json.Unmarshal(b, &dm)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				HttpStatus: http.StatusBadRequest,
				Message: fmt.Sprintf("unable to marshal AttributeMap for Directory '%s'",
					item.GetId(),
				),
			})
			break
		}
		sts = me.setDirectoryModelId(dm, item.GetId())
		if is.Error(sts) {
			break
		}
	}
	return dm, sts
}

func AssertDirectoryModel(item ItemModeler) (dm *DirectoryModel, sts Status) {
	dm, ok := item.(*DirectoryModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a %T: %s",
				(*DirectoryModel)(nil),
				item.GetId(),
			),
		})
	}
	return dm, sts
}

func GetDirectoryFilterMap() FilterMap {
	return FilterMap{}
}

func FindDirectory(dir types.Dir) (item ItemModeler, sts Status) {
	d := &DirectoryModel{}
	for range only.Once {
		if !util.DirExists(dir) {
			sts = status.Fail().
				SetMessage("directory '%s' does not exist", dir).
				SetHttpStatus(http.StatusNotFound)
			break
		}
		d.Directory = dir
		sts = status.Success("directory '%s' exists", dir)
	}
	return d, sts
}

func CanAddDirectory(dir types.Dir) (sts Status) {
	for range only.Once {
		hd, err := homedir.Dir()
		if err != nil {
			sts = status.Wrap(err)
			break
		}
		if !strings.HasPrefix(string(dir), hd) {
			sts = status.Fail().
				SetMessage("cannot add directory '%s'", dir).
				SetDetail("directory '%s' is not within your home directory '%s'", dir, hd).
				SetHttpStatus(http.StatusBadRequest)
			break
		}
		if util.DirExists(dir) {
			sts = status.Fail().
				SetMessage("cannot add directory '%s'", dir).
				SetDetail("directory '%s' already exists", dir).
				SetHttpStatus(http.StatusConflict)

			break
		}
	}
	return sts
}

func AddDirectory(dir types.Dir) (sts Status) {
	for range only.Once {
		sts = CanAddDirectory(dir)
		if is.Error(sts) {
			break
		}
		err := os.Mkdir(string(dir), os.ModePerm)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to add directory '%s'", dir)
			break
		}
		sts = status.Success("directory '%s' added", dir)
	}
	return sts
}

func UnescapeDirectory(dir types.Dir) (d types.Dir, sts Status) {
	for range only.Once {
		cleandir, err := url.QueryUnescape(string(dir))
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to unencode directory '%d'", dir).
				SetHttpStatus(http.StatusBadRequest)
		}
		d = types.Dir(cleandir)
	}
	return d, sts
}
