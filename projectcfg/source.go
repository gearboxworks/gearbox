package projectcfg

type Source struct {
	WebRoot      string       `json:"web_root,omitempty"`
	Vendors   VendorBag    `json:"vendors,omitempty"`
	Repositories Repositories `json:"repositories,omitempty"`
	Branches     BranchMap    `json:"branches,omitempty"`
}
