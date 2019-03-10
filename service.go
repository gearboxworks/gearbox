package gearbox

type ServiceName string

type ServiceMap map[ServiceName]*Service

type Services []*Service

type Service struct {
	Name        ServiceName `json:"name"`
	Label       string      `json:"label"`
	ShortLabel  string      `json:"short_label"`
	Examples    []string    `json:"examples"`
	StackName   StackName   `json:"stack"`
	ServiceType string      `json:"service_type"`
	Optional    bool        `json:"optional"`
}
type ServiceArgs Service

func (me *Service) String() string {
	return me.Label
}
