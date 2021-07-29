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

type PortPoint struct {
	Node string
	Port string
}

func NewCustomModule(
	deps Deps,
	in InportsInterface,
	out OutportsInterface,
	workers Workers,
	net Net,
) Module {
	return customModule{
		deps:    deps,
		in:      in,
		out:     out,
		workers: workers,
		net:     net,
	}
}
