package modeler

const NoFilterPath FilterPath = "/"

var NoFilter Filter = Filter{
	Name: "No Filter",
	Path: NoFilterPath,
	Filter: func(ii Item) Item {
		return ii
	},
}
