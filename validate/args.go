package validate

import "gearbox/config"

type Args struct {
	MustBeEmpty    bool
	MustNotBeEmpty bool
	MustNotExist   bool
	MustExist      bool
	MustSucceed    bool
	MustFail       bool
	MustExpire     bool
	MustNotExpire  bool
	ApiHelpUrl     string
	Config         config.Configer
}
