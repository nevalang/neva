package core

import (
	"fmt"
)

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

func checkAllPorts(got, want Interface) error {
	if err := checkPorts(
		PortsInterface(got.In),
		PortsInterface(want.In),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	if err := checkPorts(
		PortsInterface(got.Out),
		PortsInterface(want.Out),
	); err != nil {
		return fmt.Errorf("incompatible inPorts: %w", err)
	}

	return nil
}

func checkPorts(got, want PortsInterface) error {
	if len(got) < len(want) {
		return fmt.Errorf(
			"not enough ports: got %d, want %d",
			len(got),
			len(want),
		)
	}

	for name := range want {
		if want[name] != got[name] {
			return fmt.Errorf(
				"incompatible types on port '%s': got '%s', want '%s'",
				name,
				want[name],
				got[name],
			)
		}
	}

	return nil
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
