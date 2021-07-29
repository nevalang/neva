package core

import "errors"

type customModule struct {
	deps    Deps
	in      InportsInterface
	out     OutportsInterface
	workers Workers
	net     Net
}

func (cm customModule) Interface() Interface {
	return Interface{
		In:  cm.in,
		Out: cm.out,
	}
}

func (mod customModule) Validate() error {
	// mod.deps.validate()
	// mod.in.validate()
	// mod.out.validate()
	// mod.workers.validate()
	// mod.net.validate()
	return nil
}

type Deps map[string]Interface

func (d Deps) compat(name string, io Interface) error {
	for port, t := range io.In {
		if err := d[name].In[port].Compare(t); err != nil {
			return err
		}
	}
	for port, t := range io.Out {
		if err := d[name].Out[port].Compare(t); err != nil {
			return err
		}
	}
	return nil
}

type Workers map[string]string

func (w Workers) Interface(name string, deps Deps) (Interface, error) {
	i, ok := deps[name]
	if !ok {
		return Interface{}, errors.New("..")
	}
	return i, nil
}

type Net []Subscription

type Subscription struct {
	Sender    PortPoint
	Recievers []PortPoint
}

// PortPoint represents NormPortPoint and ArrPortPoint.
type PortPoint interface{}

type NormPortPoint struct {
	Node string
	Port string
}

type ArrPortPoint struct {
	Node  string
	Port  string
	Index uint8
}

func NewCustomModule(
	deps Deps,
	in InportsInterface,
	out OutportsInterface,
	workers Workers,
	net Net,
) (Module, error) {
	m := customModule{
		deps:    deps,
		in:      in,
		out:     out,
		workers: workers,
		net:     net,
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	return m, nil
}
