package apimodeler

const NoFilterPath FilterPath = "/"

var NoFilter Filter = Filter{
	Label: "No Filter",
	Path:  NoFilterPath,
	ItemFilter: func(ii Itemer) Itemer {
		return ii
	},
}
