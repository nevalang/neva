package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime/src"
	"golang.org/x/sync/errgroup"
)

type Runtime struct {
	decoder   Decoder
	connector Connector
	effector  Effector
}

var (
	ErrDecoder           = errors.New("decoder")
	ErrStartPortNotFound = errors.New("start port not found")
	ErrConnector         = errors.New("connector")
	ErrEffector          = errors.New("effector")
	ErrUnknownMsgType    = errors.New("unknown message type")
	ErrStartPortBlocked  = errors.New("start port blocked")
)

func (r Runtime) Run(ctx context.Context, bb []byte) error {
	prog, err := r.decoder.Decode(bb)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDecoder, err)
	}

	_, ok := prog.Ports[prog.StartPort]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPort)
	}

	ports := r.buildPorts(prog.Ports)

	conns, err := r.buildConnections(ports, prog.Connections)
	if err != nil {
		return fmt.Errorf("build connections: %w", err)
	}

	effects, err := r.buildEffects(ports, prog.Effects)
	if err != nil {
		return fmt.Errorf("build effects: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := r.connector.Connect(ctx, conns); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.effector.MakeEffects(ctx, effects); err != nil {
			return fmt.Errorf("%w: %v", ErrEffector, err)
		}
		return nil
	})

	return g.Wait()
}

func (r Runtime) buildPorts(in map[src.AbsolutePortAddr]uint8) map[src.AbsolutePortAddr]chan core.Msg {
	out := make(
		map[src.AbsolutePortAddr]chan core.Msg,
		len(in),
	)
	for addr, buf := range in {
		out[addr] = make(chan core.Msg, buf)
	}
	return out
}

func (r Runtime) buildConnections(ports map[src.AbsolutePortAddr]chan core.Msg, in []src.Connection) ([]Connection, error) {
	cc := make([]Connection, 0, len(in))
	for _, srcConn := range in {
		senderPort, ok := ports[srcConn.SenderPortAddr]
		if !ok {
			return nil, fmt.Errorf("%w: %v", core.ErrPortNotFound, srcConn.SenderPortAddr)
		}

		sender := Sender{
			addr: srcConn.SenderPortAddr,
			port: senderPort,
		}

		rr := make([]Receiver, 0, len(srcConn.ReceiversConnectionPoints))
		for _, srcReceiverPoint := range srcConn.ReceiversConnectionPoints {
			receiverPort, ok := ports[srcReceiverPoint.PortAddr]
			if !ok {
				return nil, fmt.Errorf("%w: %v", core.ErrPortNotFound, srcConn.SenderPortAddr)
			}

			rr = append(rr, Receiver{
				point: srcReceiverPoint,
				port:  receiverPort,
			})
		}

		cc = append(cc, Connection{
			sender:    sender,
			receivers: rr,
		})
	}

	return cc, nil
}

func (r Runtime) buildEffects(ports map[src.AbsolutePortAddr]chan core.Msg, in src.Effects) (Effects, error) {
	consts, err := r.buildConstEffects(ports, in.Constants)
	if err != nil {
		return Effects{}, fmt.Errorf("build const effects: %w", err)
	}

	ops, err := r.buildOperatorEffects(ports, in.Operators)
	if err != nil {
		return Effects{}, fmt.Errorf("build operator effects: %w", err)
	}

	return Effects{consts, ops}, nil
}

func (r Runtime) buildConstEffects(
	ports map[src.AbsolutePortAddr]chan core.Msg,
	in map[src.AbsolutePortAddr]src.Msg,
) ([]ConstEffect, error) {
	result := make([]ConstEffect, 0, len(in))

	for addr, msg := range in {
		port, ok := ports[addr]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrStartPortNotFound, addr)
		}

		msg, err := r.buildCoreMsg(msg)
		if err != nil {
			return nil, fmt.Errorf("build core msg: %w", err)
		}

		result = append(result, ConstEffect{
			Port: port,
			Msg:  msg,
		})
	}

	return result, nil
}

func (r Runtime) buildOperatorEffects(
	ports map[src.AbsolutePortAddr]chan core.Msg,
	in []src.OperatorEffect,
) ([]OperatorEffect, error) {
	result := make([]OperatorEffect, 0, len(in))

	for _, srcOpEffect := range in {
		io := core.IO{
			In:  make(core.Ports, len(srcOpEffect.PortAddrs.In)),
			Out: make(core.Ports, len(srcOpEffect.PortAddrs.Out)),
		}

		for _, addr := range srcOpEffect.PortAddrs.In {
			port, ok := ports[addr]
			if !ok {
				return nil, fmt.Errorf("%w: %v", core.ErrPortNotFound, addr)
			}
			relativeAddr := core.RelativePortAddr{
				Port: addr.Port,
				Idx:  addr.Idx,
			}
			io.In[relativeAddr] = port
		}

		for _, addr := range srcOpEffect.PortAddrs.Out {
			port, ok := ports[addr]
			if !ok {
				return nil, fmt.Errorf("%w: %v", core.ErrPortNotFound, addr)
			}
			relativeAddr := core.RelativePortAddr{
				Port: addr.Port,
				Idx:  addr.Idx,
			}
			io.Out[relativeAddr] = port
		}

		result = append(result, OperatorEffect{
			Ref: srcOpEffect.Ref,
			IO:  io,
		})
	}

	return result, nil
}

func (r Runtime) buildCoreMsg(in src.Msg) (core.Msg, error) {
	var out core.Msg

	switch in.Type {
	case src.IntMsg:
		out = core.NewIntMsg(in.Int)
	case src.BoolMsg:
		out = core.NewBoolMsg(in.Bool)
	case src.StrMsg:
		out = core.NewStrMsg(in.Str)
	case src.StructMsg:
		structMsg := make(map[string]core.Msg, len(in.Struct))
		for field, value := range in.Struct {
			v, err := r.buildCoreMsg(value)
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

func MustNew(
	decoder Decoder,
	connector Connector,
	effector Effector,
) Runtime {
	utils.NilPanic(decoder, effector, connector)

	return Runtime{
		decoder:   decoder,
		effector:  effector,
		connector: connector,
	}
}
