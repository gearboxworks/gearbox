package projectcfg

type Dev struct {
	Stack        Stack          `json:"stack"`
	Environments EnvironmentMap `json:"environments,omitempty"`
}
