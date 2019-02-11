package gearbox

type WordPress struct {
	RootPath    string `json:"root_path"`
	CorePath    string `json:"core_path"`
	ContentPath string `json:"content_path"`
	VendorPath  string `json:"vendor_path"`
}
