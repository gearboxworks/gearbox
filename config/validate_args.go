package config

type ValidateArgs struct {
	MustBeEmpty    bool
	MustNotBeEmpty bool
	MustNotExist   bool
	MustExist      bool
	ApiHelpUrl     string
	Config         Configer
}
