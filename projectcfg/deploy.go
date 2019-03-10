package projectcfg

type Deploy struct {
	Provider    string        `json:"provider"`
	TargetName  string        `json:"target_name"`
	TargetId    string        `json:"target_id,omitempty"`
	WebRoot     string        `json:"web_root,omitempty"`
	Credentials Credentials   `json:"credentials,omitempty"`
	Vendors  VendorBag     `json:"vendors,omitempty"`
	Files       FileMap       `json:"files,omitempty"`
	Hosts       DeployHostMap `json:"hosts"`
}
