package gearbox



type DefaultGroupGetter interface {
	GetDefaultGroup() string
}

var _ DefaultGroupGetter = (*ServiceId)(nil)
func (me *ServiceId) GetDefaultGroup() string {
	return "gearbox"
}


type ServiceId struct {
	*Identity
}

//func NewServiceId() (sid *ServiceId) {
//	return &ServiceId{
//		Identity: NewIdentity(),
//	}
//
//}
//
//func (me *ServiceId) SetDefaultGroup(defaultGroup string) {
//	me.Identity.defaultGroup = defaultGroup
//}
