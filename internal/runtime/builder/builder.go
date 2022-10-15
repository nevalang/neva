package builder

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
	"golang.org/x/sync/errgroup"
)

type Builder struct{}

func (b Builder) Build(prog src.Program) (runtime.Executable, error) {
	var (
		g           errgroup.Group
		ports       = b.buildPorts(prog.Ports)
		connections []runtime.Connection
		effects     runtime.Effects
	)

	g.Go(func() error {
		var err error
		connections, err = b.buildConnections(ports, prog.Connections)
		if err != nil {
			return fmt.Errorf("build connections: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		effects, err = b.buildEffects(ports, prog.Fx)
		if err != nil {
			return fmt.Errorf("build effects: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return runtime.Executable{}, fmt.Errorf("wait: %w", err)
	}

	return runtime.Executable{
		Start: prog.Start,
		Ports: ports,
		Net:   connections,
		Fx:    effects,
	}, nil
}

func (b Builder) buildPorts(in src.PortSet) runtime.Ports {
	out := make(
		runtime.Ports,
		len(in),
	)
	for addr, buf := range in {
		out[addr] = make(chan core.Msg, buf)
	}
	return out
}

func (b Builder) buildConnections(ports runtime.Ports, srcConns []src.Connection) ([]runtime.Connection, error) {
	cc := make([]runtime.Connection, 0, len(srcConns))

	for _, srcConn := range srcConns {
		c, err := b.buildConnection(ports, srcConn)
		if err != nil {
			return nil, fmt.Errorf("build connection: err %w, conn %v", err, srcConn)
		}

		cc = append(cc, c)
	}

	return cc, nil
}

func (b Builder) buildConnection(ports runtime.Ports, srcConn src.Connection) (runtime.Connection, error) {
	senderPort, ok := ports[srcConn.SenderSide]
	if !ok {
		return runtime.Connection{}, fmt.Errorf("%w: %v", core.ErrPortNotFound, srcConn.SenderSide)
	}

	rr := make([]chan core.Msg, 0, len(srcConn.ReceiverSides))
	for _, srcReceiverPoint := range srcConn.ReceiverSides {
		receiverPort, ok := ports[srcReceiverPoint.PortAddr]
		if !ok {
			return runtime.Connection{}, fmt.Errorf("%w: %v", core.ErrPortNotFound, srcConn.SenderSide)
		}

		rr = append(rr, receiverPort)
	}

	return runtime.Connection{
		Src:       srcConn,
		Sender:    senderPort,
		Receivers: rr,
	}, nil
}

func (b Builder) buildEffects(ports runtime.Ports, effects src.Fx) (runtime.Effects, error) {
	c, err := b.buildConstEffects(ports, effects.Const)
	if err != nil {
		return runtime.Effects{}, fmt.Errorf("build const effects: %w", err)
	}

	o, err := b.buildOperatorEffects(ports, effects.Func)
	if err != nil {
		return runtime.Effects{}, fmt.Errorf("build operator effects: %w", err)
	}

	t, err := b.buildTriggerEffects(ports, effects.Trigger)
	if err != nil {
		return runtime.Effects{}, fmt.Errorf("build operator effects: %w", err)
	}

	return runtime.Effects{
		Const:   c,
		Func:    o,
		Trigger: t,
	}, nil
}

var ErrPortNotFound = errors.New("port not found")

func (b Builder) buildConstEffects(
	ports runtime.Ports,
	in map[src.PortAddr]src.Msg,
) ([]runtime.ConstFx, error) {
	result := make([]runtime.ConstFx, 0, len(in))

	for addr, msg := range in {
		port, ok := ports[addr]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
		}

		msg, err := b.buildCoreMsg(msg)
		if err != nil {
			return nil, fmt.Errorf("build core msg: %w", err)
		}

		result = append(result, runtime.ConstFx{
			OutPort: port,
			Msg:     msg,
		})
	}

	return result, nil
}

func (b Builder) buildOperatorEffects(
	ports runtime.Ports,
	ops []src.FuncFx,
) ([]runtime.FuncFx, error) {
	result := make([]runtime.FuncFx, 0, len(ops))

	for _, srcOpEffect := range ops {
		io := core.IO{
			In:  make(core.Ports, len(srcOpEffect.Ports.In)),
			Out: make(core.Ports, len(srcOpEffect.Ports.Out)),
		}

		for _, addr := range srcOpEffect.Ports.In {
			port, ok := ports[addr]
			if !ok {
				return nil, fmt.Errorf("%w: %v", core.ErrPortNotFound, addr)
			}
			relativeAddr := core.PortAddr{
				Port: addr.Port,
				Idx:  addr.Idx,
			}
			io.In[relativeAddr] = port
		}

		for _, addr := range srcOpEffect.Ports.Out {
			port, ok := ports[addr]
			if !ok {
				return nil, fmt.Errorf("%w: %v", core.ErrPortNotFound, addr)
			}
			relativeAddr := core.PortAddr{
				Port: addr.Port,
				Idx:  addr.Idx,
			}
			io.Out[relativeAddr] = port
		}

		result = append(result, runtime.FuncFx{
			Ref: srcOpEffect.Ref,
			IO:  io,
		})
	}

	return result, nil
}

var ErrUnknownMsgType = errors.New("unknown message type")

func (b Builder) buildCoreMsg(in src.Msg) (core.Msg, error) {
	var out core.Msg

	switch in.Type {
	case src.IntMsg:
		out = core.NewIntMsg(in.Int)
	case src.BoolMsg:
		out = core.NewBoolMsg(in.Bool)
	case src.StrMsg:
		out = core.NewStrMsg(in.Str)
	case src.DictMsg:
		structMsg := make(map[string]core.Msg, len(in.Dict))
		for field, value := range in.Dict {
			v, err := b.buildCoreMsg(value)
			if err != nil {
				return nil, fmt.Errorf("core msg: %w", err)
			}
			structMsg[field] = v
		}
		out = core.NewDictMsg(structMsg)
	default:
		return nil, fmt.Errorf("%w: %v", ErrUnknownMsgType, in.Type)
	}

	return out, nil
}

func (b Builder) buildTriggerEffects(
	ports runtime.Ports,
	in []src.TriggerFx,
) ([]runtime.TriggerFx, error) {
	result := make([]runtime.TriggerFx, 0, len(in))

	for _, effect := range in {
		inPort, ok := ports[effect.In]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, effect.In)
		}

		outPort, ok := ports[effect.Out]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrPortNotFound, effect.In)
		}

		msg, err := b.buildCoreMsg(effect.Msg)
		if err != nil {
			return nil, fmt.Errorf("build core msg: %w", err)
		}

		result = append(result, runtime.TriggerFx{
			InPort:  inPort,
			OutPort: outPort,
			Msg:     msg,
		})
	}

	return result, nil
}
