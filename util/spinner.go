// Interface to spinner package.
// Created to dissociate from a specific spinner package.
package util

import (
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"os"
)

type Spinner struct {
	Style    spin.Name
	Text     string
	Instance *wow.Wow
	ExitOK   string
	ExitNOK  string
	Enable   bool
}
type SpinnerArgs Spinner

// //////////////////////////////////////////////////////////////////////////////
// Low-level related
func NewSpinner(args ...SpinnerArgs) *Spinner {
	var _args SpinnerArgs
	if len(args) > 0 {
		_args = args[0]
	}

	if _args.Style == 0 {
		_args.Style = spin.Clock
	}

	if _args.ExitOK == "" {
		_args.ExitOK = "OK"
	}

	if _args.ExitOK == "" {
		_args.ExitOK = "FAILED"
	}

	_args.Instance = wow.New(os.Stdout, spin.Get(_args.Style), _args.Text)

	spinner := &Spinner{}
	*spinner = Spinner(_args)

	return spinner
}

func (me *Spinner) Start() {

	if me.Instance.IsTerminal == true {
		me.Instance.Text(me.Text).Spinner(spin.Get(me.Style))
		me.Instance.Start()
	}

	return
}

func (me *Spinner) Stop(state bool) {

	if me.Instance.IsTerminal == true {
		if state == true {
			me.Instance.PersistWith(spin.Spinner{Frames: []string{"👍"}}, " "+me.ExitOK)
		} else {
			me.Instance.PersistWith(spin.Spinner{Frames: []string{"👎"}}, " "+me.ExitNOK)
		}

		me.Instance.Stop()
	}

	return
}

func (me *Spinner) Update(displayString string) {

	me.Text = displayString

	return
}

/*

	if _args.Enable == false {
		return nil
	}

*/
