package ab

type FilterName string
type FilterPath string
type FilterMap map[FilterPath]Filter
type Filters []Filter
type Critera interface{}
type Filter struct {
	Name   FilterName
	Path   FilterPath
	Filter func(Item) Item
}

const NoFilterPath FilterPath = "/"

var NoFilter Filter = Filter{
	Name: "No Filter",
	Path: NoFilterPath,
	Filter: func(ii Item) Item {
		return ii
	},
}
