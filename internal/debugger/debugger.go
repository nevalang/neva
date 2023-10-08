package debugger

import (
	"sync"

	"github.com/nevalang/neva/internal/runtime"
	"golang.org/x/exp/maps"
)

// Debugger implements runtime.Listener interface in a way that it can be used to debug programs.
// It allows to set breakpoints on program's events, observe program's state and time-travel.
// WARNING: every `Send` call locks debugger's mutex, so it dramatically slows down the runtime.
type Debugger struct {
	mu    sync.Mutex
	state State
}

// TODO figure out best shape
type State map[runtime.PortAddr]runtime.Event

func (d *Debugger) Send(event runtime.Event, msg runtime.Msg) runtime.Msg {
	// lock other `Send` calls from runtime so we don't miss anything in the background
	d.mu.Lock()
	defer d.mu.Unlock()

	// TODO: Update state
	// TODO: handle breakpoint

	return msg
}

// State returns copy of debuger's state.
func (d *Debugger) State() State {
	var state State
	maps.Copy(state, d.state)
	return state
}

func New(breakpoints map[runtime.PortAddr]runtime.Event) Debugger {
	return Debugger{
		state: map[runtime.PortAddr]runtime.Event{},
	}
}
