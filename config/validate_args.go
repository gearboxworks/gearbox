package config

type ValidateArgs struct {
	MustNotBeEmpty bool
	MustNotExist   bool
	MustExist      bool
	ApiHelpUrl     string
	Config         Configer
}
