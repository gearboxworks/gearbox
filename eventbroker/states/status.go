package states

import (
	"gearbox/eventbroker/msgs"
	"sync"
	"time"
)

type Status struct {
	EntityId   *msgs.Address
	EntityName *msgs.Address
	ParentId   *msgs.Address
	Current    State
	Want       State
	Last       State
	LastWhen   time.Time
	Attempts   int
	Error      error
	Action     Action

	mutex *sync.RWMutex // Mutex control for map.
}
