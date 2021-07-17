package runtime

// CustomModule composes other modules.
type CustomModule struct {
	in      InportsInterface
	out     OutportsInterface
	workers map[string]Module
}

func (cm CustomModule) Interface() (InportsInterface, OutportsInterface) {
	return cm.in, cm.out
}

func (m CustomModule) Run() {
	m.Worker().Run()
}

func (m CustomModule) Worker() Worker {
	// cr := []ChanRelation{}
	return Worker{}
}

type Env map[string]Module

// ChanRelation represents one-to-many relation between sender and receiver channels.
type ChanRelation struct {
	Sender    <-chan Msg
	Receivers []chan Msg
}

func NewCustomModule(
	in InportsInterface,
	out OutportsInterface,
) CustomModule {
	return CustomModule{
		in:  in,
		out: out,
	}
}
