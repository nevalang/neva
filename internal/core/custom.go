package core

type CustomModule struct {
	deps    Deps
	in      InportsInterface
	out     OutportsInterface
	workers Workers
	net     Net
}

func (cm CustomModule) Interface() Interface {
	return Interface{
		In:  cm.in,
		Out: cm.out,
	}
}

type Workers map[string]string

type Net []Subscription

type Subscription struct {
	Sender    PortPoint
	Recievers []PortPoint
}

type PortPoint struct {
	Node string
	Port string
}

type Connection struct {
	Sender    chan Msg
	Receivers []chan Msg
}

func NewCustomModule(deps Deps, in InportsInterface, out OutportsInterface, workers Workers, net Net) Module {
	return CustomModule{
		deps:    deps,
		in:      in,
		out:     out,
		workers: workers,
		net:     net,
	}
}
