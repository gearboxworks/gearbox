package api

type Status struct {
	Success    bool
	StatusCode int
	Help       string
	Error      error
}
type StatusResponse struct {
	Help  string `json:"help"`
	Error string `json:"error"`
}

func (me *Status) ToResponse() *StatusResponse {
	var msg string
	if me.Error != nil {
		msg = me.Error.Error()
	}
	return &StatusResponse{
		Help:  me.Help,
		Error: msg,
	}
}
