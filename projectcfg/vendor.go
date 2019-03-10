package projectcfg

type VendorName string

type VendorBag map[VendorName]interface{}

type VendorMap map[VendorName]Vendor

type Vendors []Vendor

type Vendor interface {
	GetName() string
	GetInstance() interface{}
}

//var _ json.Unmarshaler = (*VendorBag)(nil)
//func (me *VendorBag) UnmarshalJSON(data []byte) error {
//	fwb := make(map[string]interface{},0)
//	err := json.Unmarshal(data,&fwb)
//	if err != nil {
//		panic(err)
//	}
//	return err
//}
