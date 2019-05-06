package config

type ValidateArgs struct {
	MustBeEmpty     bool
	MustNotBeEmpty  bool
	MustNotExist    bool
	MustExist       bool
	MustBeOnDisk    bool
	MustNotBeOnDisk bool
	MustSucceed     func() Status
	MustNotEqual    interface{}
	MustBeIn        interface{}
	MustNotBeIn     interface{}
	IgnoreCurrent   bool
	ApiHelpUrl      string
	Config          Configer
}
