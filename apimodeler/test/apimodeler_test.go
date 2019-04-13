package test

import (
	"gearbox/apimodeler"
	"gearbox/only"
	"gearbox/status/is"
	"gearbox/types"
	"testing"
)

const IdParamsSlugifyWanted = ":foo/:bar/:baz"

func TestIdParams(t *testing.T) {
	idp := apimodeler.IdParams{"foo", "bar", "baz"}
	if idp.Slugify() != IdParamsSlugifyWanted {
		t.Errorf("idparams.Slugify(); got '%s', wanted: '%s'", idp.Slugify(), IdParamsSlugifyWanted)
	}
}

func TestModels(t *testing.T) {

	t.Run("GetBasepath()", func(t *testing.T) {
		ms := apimodeler.NewController(NewTestableController())
		if ms.GetBasepath() != testableModelBasepath {
			t.Errorf("List basepath is not '%s'", testableModelBasepath)
		}
		if ms.GetBasepath() != testableModelBasepath {
			t.Errorf("List basepath is not '%s'", testableModelBasepath)
		}
	})

	t.Run("GetIdFromUrl()", func(t *testing.T) {
		ms := apimodeler.NewController(NewTestableController())
		ctx := &TestableContext{}
		itemid, sts := apimodeler.GetIdFromUrl(ctx, ms)
		if is.Error(sts) {
			t.Errorf("unable to get item Id from context: %s", sts.Message())
			return
		}
		wanted := apimodeler.ItemId("alpha/beta")
		if itemid != wanted {
			t.Errorf("item ID; got '%s', wanted: '%s'", itemid, wanted)
		}
	})

	t.Run("GetIdTemplate()", func(t *testing.T) {
		ms := apimodeler.NewController(NewTestableController())
		template := apimodeler.GetIdTemplate(ms)
		wanted := types.UrlTemplate(":foo/:bar")
		if template != wanted {
			t.Errorf("template; got '%s', wanted: '%s'", template, wanted)
		}

	})

	t.Run("GetIdParams()", func(t *testing.T) {
		ms := apimodeler.NewController(NewTestableController())
		params := apimodeler.GetIdParams(ms)
		if len(params) != len(testableModelIdParams) {
			t.Errorf("len(GetIdParams()); got '%d', wanted: '%d'", len(params), len(testableModelIdParams))
			return
		}
		for i, p := range params {
			if p != testableModelIdParams[i] {
				t.Errorf("ID params[%d]; got '%s', wanted: '%s'", i, p, testableModelIdParams[i])
			}
		}
	})
	t.Run("GetResourceUrlTemplate()", func(t *testing.T) {
		ms := apimodeler.NewController(NewTestableController())
		template := ms.GetResourceUrlTemplate()
		wanted := types.UrlTemplate("/foo/:foo/:bar")
		if template != wanted {
			t.Errorf("template; got '%s', wanted: '%s'", template, wanted)
		}
	})
	t.Run("GetRouteNamePrefix()", func(t *testing.T) {
		ms := apimodeler.NewController(NewTestableController())
		prefix := apimodeler.GetRouteNamePrefix(ms)
		wanted := "foo"
		if prefix != wanted {
			t.Errorf("template; got '%s', wanted: '%s'", prefix, wanted)
		}
		//@TODO Test when Controller have children
	})

}

func TestItem(t *testing.T) {
	ti := &TestableItem{
		Id:   "foo",
		Type: "bar",
	}
	t.Run("GetId()", func(t *testing.T) {
		if ti.GetId() != "foo" {
			t.Errorf("item.GetId(); got '%s', wanted: 'foo'", ti.GetId())
		}
	})
	t.Run("GetType()", func(t *testing.T) {
		if ti.GetType() != "bar" {
			t.Errorf("item.GetType(); got '%s', wanted: 'bar'", ti.GetType())
		}
	})
	t.Run("GetItem()", func(t *testing.T) {
		item, sts := ti.GetItem()
		if is.Error(sts) {
			t.Errorf("item.GetItem(); got error: %s", sts.Message())
		}
		if item != ti {
			t.Errorf("item.GetItem(); got '%+v', wanted: '%+v'", item, ti)
		}
	})

}

func TestModeler(t *testing.T) {

	t.Run("GetBasepath()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		if m.GetBasepath() != testableModelBasepath {
			t.Errorf("connections.GetBasepath(); got '%s', wanted: '%s'", m.GetBasepath(), testableModelBasepath)
		}
	})

	t.Run("GetIdParams()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		for range only.Once {
			idp := m.GetIdParams()

			if len(idp) != 2 {
				t.Errorf("len(GetIdParams()); got '%d', wanted: '%d'", len(idp), 2)
				break
			}
			if idp[0] != testableModelIdParams[0] {
				t.Errorf("GetIdParams()[0]; got '%s', wanted: '%s'", idp[0], testableModelIdParams[0])
			}

			if idp[1] != testableModelIdParams[1] {
				t.Errorf("GetIdParams()[1]; got '%s', wanted: '%s'", idp[1], testableModelIdParams[1])
			}
		}
	})

	t.Run("GetList()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		coll, sts := m.GetList()
		if is.Error(sts) {
			t.Errorf(sts.Message())
		}
		if len(testableItemData) != len(coll) {
			t.Errorf("expected same length as testable item data. Got ID: %d, wanted: %d",
				len(coll),
				len(testableItemData),
			)
			return
		}
		for i, ti := range coll {
			typ, ok := testableItemData[ti.GetId()]
			if !ok {
				t.Errorf("ID '%s' not found", ti.GetId())
				continue
			}
			if typ != coll[i].GetType() {
				t.Errorf("got Type: %s, wanted: %s", coll[i].GetType(), typ)
				continue
			}
			if typ != ti.GetType() {
				t.Errorf("got Type: %s, wanted: %s", ti.GetType(), typ)
			}
		}
	})

	t.Run("AddItem()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		ti := &TestableItem{
			Id:   "42",
			Type: "hitchhiker",
		}
		sts := m.AddItem(ti)
		if is.Error(sts) {
			t.Errorf("unable to add testable item ID %s: %s", ti.Id, sts.Message())
			return
		}
		coll, sts := m.GetList()
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		if len(coll) == 0 {
			t.Error("returned empty List")
			return
		}
		if coll[len(coll)-1] == nil {
			t.Error("returned List with nil items")
			return
		}
		ti2, ok := coll[len(coll)-1].(*TestableItem)
		if !ok {
			t.Error("elements of returned List not *TestableItem")
			return
		}
		if ti2.Id != ti.Id {
			t.Errorf("got ID: %s, wanted: %s", ti2.Id, ti.Id)
		}
		if ti2.Type != ti.Type {
			t.Errorf("got Type: %s, wanted: %s", ti2.Type, ti.Type)
			return
		}
	})

	t.Run("DeleteItem()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		coll, sts := m.GetList()
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		wantLen := len(coll) - 1
		item := coll[wantLen]
		sts = m.DeleteItem(item.GetId())
		if is.Error(sts) {
			t.Errorf("item '%s' not found in List", item.GetId())
			return
		}
		coll2, _ := m.GetList()
		if wantLen != len(coll2) {
			t.Errorf("got len: %d, wanted: %d", len(coll2), wantLen)
		}
		for _, ti := range coll2 {
			if ti.GetId() != item.GetId() {
				continue
			}
			t.Errorf("item '%s' not deleted", item.GetId())
		}
	})

	t.Run("UpdateItem()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		coll, sts := m.GetList()
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		coll = append(coll[:0:0], coll...) // Clone

		item := coll[len(coll)-1]
		var newtype string
		if item.GetType() == FrobinatorType {
			newtype = UnicornType
		} else {
			newtype = FrobinatorType
		}
		newitem := &TestableItem{
			Id:   item.GetId(),
			Type: apimodeler.ItemType(newtype),
		}
		sts = m.UpdateItem(newitem)
		if is.Error(sts) {
			t.Errorf("item '%s' not found", item.GetId())
			return
		}
		coll2, _ := m.GetList()
		if len(coll) != len(coll2) {
			t.Errorf("got len: %d, wanted: %d", len(coll2), len(coll))
			return
		}
		item2 := getItem(coll2, item.GetId())
		if item == item2 {
			t.Errorf("item not updated")
			return
		}
		if item2 == nil {
			t.Errorf("item '%s' not found after update", item.GetId())
			return
		}
		if item2.GetType() != apimodeler.ItemType(newtype) {
			t.Errorf("item '%s' not type '%s' after update", item.GetId(), newtype)
			return
		}
	})

	t.Run("GetFilterMap()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		fm := m.GetFilterMap()
		if len(fm) != 2 {
			t.Errorf("filter map len; got: %d, wanted: %d", len(fm), 2)
		}
		if f, ok := fm[FrobinatorsFilter]; !ok {
			t.Errorf("no '%s' filter found", FrobinatorsFilter)
		} else {
			if f.Path != FrobinatorsFilter {
				t.Errorf("path for filter '%s' not equal to '%s'", FrobinatorsFilter, FrobinatorsFilter)
			}
			if f.ItemFilter != nil {
				t.Errorf("item filter for '%s' not nil", FrobinatorsFilter)
			}
			if f.ListFilter == nil {
				t.Errorf("List filter for '%s' is nil", FrobinatorsFilter)
			}
		}
		if f, ok := fm[UnicornFilter]; !ok {
			t.Errorf("no '%s' filter found", UnicornFilter)
		} else {
			if f.Path != UnicornFilter {
				t.Errorf("path for filter '%s' not equal to '%s'", UnicornFilter, UnicornFilter)
			}
			if f.ItemFilter == nil {
				t.Errorf("item filter for '%s' is nil", UnicornFilter)
			}
			if f.ListFilter != nil {
				t.Errorf("List filter for '%s' not nil", UnicornFilter)
			}
		}

	})

	t.Run("GetListIds()", func(t *testing.T) {
		for range only.Once {
			m := apimodeler.NewController(NewTestableController())
			cids, sts := m.GetListIds()
			if is.Error(sts) {
				t.Error(sts.Message())
				break
			}
			if len(cids) != len(testableItemData) {
				t.Errorf("got len: %d, wanted: %d", len(cids), len(testableItemData))
				break
			}
			for _, cid := range cids {
				if _, ok := testableItemData[cid]; !ok {
					t.Errorf("item id '%s' not found", cid)
				}
			}
		}
	})

	t.Run("GetItem()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		coll, sts := m.GetList()
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		if len(coll) == 0 {
			t.Error("List is empty")
			return
		}
		ti := coll[0]
		ti2, sts := m.GetItem(ti.GetId())
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		if ti2.GetId() != ti.GetId() {
			t.Errorf("got ID: %s, wanted: %s", ti2.GetId(), ti.GetId())
		}
		if ti2.GetType() != ti.GetType() {
			t.Errorf("got Type: %s, wanted: %s", ti2.GetType(), ti.GetType())
			return
		}
	})

	t.Run("FilterList()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		fc, sts := m.FilterList(FrobinatorsFilter)
		if is.Error(sts) {
			t.Errorf("unable to filter List on '%s': %s", FrobinatorsFilter, sts.Message())
		}
		wantLen := countValues(FrobinatorType)
		if len(fc) != wantLen {
			t.Errorf("got len: %d, wanted: %d", len(fc), wantLen)
			return
		}
		for _, c := range fc {
			if c.GetType() == FrobinatorType {
				continue
			}
			t.Errorf("List contains type '%s'", c.GetType())
			return
		}
	})

	t.Run("FilterItem()", func(t *testing.T) {
		m := apimodeler.NewController(NewTestableController())
		coll, sts := m.GetList()
		if is.Error(sts) {
			t.Error("unable to get List")
			return
		}
		wantLen := countValues(UnicornType)
		filtered := make(apimodeler.List, 0)
		for i, item := range coll {
			fi, sts := m.FilterItem(item, UnicornFilter)
			if is.Error(sts) {
				t.Errorf("unable to filter item: %s", sts.Message())
				return
			}
			if fi == nil {
				if coll[i].GetType() != UnicornFilter {
					continue
				}
				t.Errorf("item '%s' of type '%s' filtered out", coll[i].GetId(), UnicornFilter)
			}
			if fi != item {
				t.Errorf("filtered item not the same;ID before: '%s' , after: '%s'", item.GetId(), fi.GetId())
			}
			if fi.GetType() != UnicornType {
				t.Errorf("filtered item has '%s' type", fi.GetType())
			}
			filtered = append(filtered, fi)
		}
		if len(filtered) != wantLen {
			t.Errorf("got len: %d, wanted: %d", len(filtered), wantLen)
			return
		}
	})

}

func getItem(coll apimodeler.List, itemid apimodeler.ItemId) (item apimodeler.Itemer) {
	for _, i := range coll {
		if i.GetId() != itemid {
			continue
		}
		item = i
		break
	}
	return item
}

func countIds(coll apimodeler.List, id apimodeler.ItemId) int {
	cnt := 0
	for _, c := range coll {
		if c.GetId() != id {
			continue
		}
		cnt++
	}
	return cnt
}

func countTypes(coll apimodeler.List, typ apimodeler.ItemType) int {
	cnt := 0
	for _, c := range coll {
		if c.GetType() != typ {
			continue
		}
		cnt++
	}
	return cnt
}

func countKeys(key apimodeler.ItemId) int {
	cnt := 0
	for k := range testableItemData {
		if k != key {
			continue
		}
		cnt++
	}
	return cnt
}

func countValues(value apimodeler.ItemType) int {
	cnt := 0
	for _, v := range testableItemData {
		if v != value {
			continue
		}
		cnt++
	}
	return cnt
}
