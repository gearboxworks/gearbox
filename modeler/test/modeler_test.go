package test

import (
	"gearbox/modeler"
	"gearbox/only"
	"gearbox/status/is"
	"testing"
)

const IdParamsSlugifyWanted = ":foo/:bar/:baz"

func TestIdParams(t *testing.T) {
	idp := modeler.IdParams{"foo", "bar", "baz"}
	if idp.Slugify() != IdParamsSlugifyWanted {
		t.Errorf("idparams.Slugify(); got '%s', wanted: '%s'", idp.Slugify(), IdParamsSlugifyWanted)
	}
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

func TestModels(t *testing.T) {

	t.Run("GetBasepath()", func(t *testing.T) {
		m := modeler.NewModels(NewTestableModel())
		if m.GetBasepath() != TestableModelBasepath {
			t.Errorf("connections.GetBasepath(); got '%s', wanted: '%s'", m.GetBasepath(), TestableModelBasepath)
		}
	})

	t.Run("GetIdParams()", func(t *testing.T) {
		m := modeler.NewModels(NewTestableModel())
		for range only.Once {
			idp := m.GetIdParams()

			if len(idp) != 2 {
				t.Errorf("len(GetIdParams()); got '%d', wanted: '%d'", len(idp), 2)
				break
			}
			if idp[0] != TestableModelIdParams[0] {
				t.Errorf("GetIdParams()[0]; got '%s', wanted: '%s'", idp[0], TestableModelIdParams[0])
			}

			if idp[1] != TestableModelIdParams[1] {
				t.Errorf("GetIdParams()[0]; got '%s', wanted: '%s'", idp[1], TestableModelIdParams[1])
			}
		}
	})

	t.Run("GetCollection()", func(t *testing.T) {
		m := modeler.NewModels(NewTestableModel())
		coll, sts := m.Self.GetCollection(modeler.NoFilterPath)
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
		m := modeler.NewModels(NewTestableModel())
		ti := &TestableItem{
			Id:   "42",
			Type: "hitchhiker",
		}
		sts := m.Self.AddItem(ti)
		if is.Error(sts) {
			t.Errorf("unable to add testable item ID %s: %s", ti.Id, sts.Message())
			return
		}
		coll, sts := m.Self.GetCollection(modeler.NoFilterPath)
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		if len(coll) == 0 {
			t.Error("returned empty collection")
			return
		}
		if coll[len(coll)-1] == nil {
			t.Error("returned collection with nil items")
			return
		}
		ti2, ok := coll[len(coll)-1].(*TestableItem)
		if !ok {
			t.Error("elements of returned collection not *TestableItem")
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
		m := modeler.NewModels(NewTestableModel())
		coll, sts := m.Self.GetCollection(modeler.NoFilterPath)
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		wantLen := len(coll) - 1
		item := coll[wantLen]
		sts = m.Self.DeleteItem(item.GetId())
		if is.Error(sts) {
			t.Errorf("item '%s' not found in collection", item.GetId())
			return
		}
		coll2, _ := m.Self.GetCollection(modeler.NoFilterPath)
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

	t.Run("FilterItem()", func(t *testing.T) {
		//m := modeler.NewModels(NewTestableModel())

	})

	t.Run("GetCollectionFilterMap()", func(t *testing.T) {
		//m := modeler.NewModels(NewTestableModel())

	})

	t.Run("GetCollectionIds()", func(t *testing.T) {
		//m := modeler.NewModels(NewTestableModel())

	})

	t.Run("GetItem()", func(t *testing.T) {
		m := modeler.NewModels(NewTestableModel())
		coll, sts := m.Self.GetCollection(modeler.NoFilterPath)
		if is.Error(sts) {
			t.Error(sts.Message())
			return
		}
		if len(coll) == 0 {
			t.Error("collection is empty")
			return
		}
		ti := coll[0]
		ti2, sts := m.Self.GetItem(ti.GetId())
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

	t.Run("UpdateItem()", func(t *testing.T) {
		//m := modeler.NewModels(NewTestableModel())

	})

}

func TestCollection(t *testing.T) {

	t.Run("Collection", func(t *testing.T) {
	})

}
