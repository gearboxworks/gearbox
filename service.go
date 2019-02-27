package gearbox

type Service struct {
	Label string `json:"label"`
}

func (me *Service) String() string {
	return me.Label
}
