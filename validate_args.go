package gearbox

type ValidateArgs struct {
	MustNotBeEmpty bool
	MustNotExist   bool
	MustExist      bool
	ApiHelpUrl     string
	Gearbox        Gearbox
}
