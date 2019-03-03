package gearbox

type validateArgs struct {
	MustNotBeEmpty bool
	MustNotExist   bool
	MustExist      bool
	ApiHelpUrl     string
	Gearbox        *Gearbox
}
