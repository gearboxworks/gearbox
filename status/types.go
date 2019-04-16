package status

type jsonS struct {
	Message string      `json:"message"`
	Help    string      `json:"help,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Status interface {
	IsSuccess() bool
	IsError() bool
	Error() string
	GetHelp(HelpType) string
	HttpStatus() int
	Message() string
	Cause() error
	SetData(interface{})
	SetCause(error)
	SetSuccess(bool)
	SetMessage(string)
	SetHttpStatus(int)
	SetHelp(HelpType, string)
	SetOtherHelp(HelpTypeMap)
	GetData() interface{}
	GetString() (string, Status)
}

type Args struct {
	Success    bool
	Help       string
	ApiHelp    string
	CliHelp    string
	OtherHelp  HelpTypeMap
	Message    string
	HttpStatus int
	Cause      error
	Data       interface{}
}

type SuccessInspector interface {
	IsSuccess() bool
}
