package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

/* === PROGRAM === */

type (
	Program struct {
		StartPortAddr PortAddr
		Ports         Ports
		Net           []Connection
		Nodes         Routines
	}

	PortAddr struct {
		Path, Name string
		Idx        uint8
	}

	Ports map[PortAddr]chan any

	Connection struct {
		Sender    ConnectionSide
		Receivers []ConnectionSide
	}

	ConnectionSide struct {
		Port chan any
		Meta ConnectionSideMeta
	}

	ConnectionSideMeta struct {
		PortAddr  PortAddr
		Selectors []Selector
	}

	Selector struct {
		RecField string // "" means use ArrIdx
		ArrIdx   int
	}

	Routines struct {
		Giver     []GiverRoutine
		Component []ComponentRoutine
	}

	GiverRoutine struct {
		OutPort chan any
		Msg     any
	}

	ComponentRoutine struct {
		Ref ComponentRef
		IO  IO
	}

	IO struct { // to core?
		In, Out map[string]chan any
	}

	ComponentRef struct {
		Pkg, Name string
	}
)

/* === RUNTIME === */

type Runtime struct {
	connector     Connector
	routineRunner RoutineRunner
}

type (
	Connector interface {
		Connect(context.Context, []Connection) error
	}
	RoutineRunner interface {
		Run(context.Context, Routines) error
	}
)

var (
	ErrStartPortNotFound = errors.New("start port not found")
	ErrConnector         = errors.New("connector")
	ErrRoutineRunner     = errors.New("routine runner")
)

func (r Runtime) Run(ctx context.Context, prog Program) error {
	startPort, ok := prog.Ports[prog.StartPortAddr]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPortAddr)
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := r.connector.Connect(gctx, prog.Net); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})
	g.Go(func() error {
		if err := r.routineRunner.Run(gctx, prog.Nodes); err != nil {
			return fmt.Errorf("%w: %v", ErrRoutineRunner, err)
		}
		return nil
	})

	startPort <- nil
	return g.Wait()
}

/*  === ROUTINE-RUNNER === */

type Runner struct {
	constant GiverRunner
	operator ComponentRunner
}

type (
	GiverRunner interface {
		Run(context.Context, []GiverRoutine) error
	}
	ComponentRunner interface {
		Run(context.Context, []ComponentRoutine) error
	}
)

var (
	ErrComponent = errors.New("component")
	ErrGiver     = errors.New("giver")
)

func (e Runner) Run(ctx context.Context, routines Routines) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := e.constant.Run(gctx, routines.Giver); err != nil {
			return errors.Join(ErrGiver, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.operator.Run(gctx, routines.Component); err != nil {
			return errors.Join(ErrComponent, err)
		}
		return nil
	})

	return g.Wait()
}

/* === GIVER-RUNNER === */

type GiverRunnerImlp struct{}

func (e GiverRunnerImlp) Run(ctx context.Context, givers []GiverRoutine) error {
	wg := sync.WaitGroup{}
	wg.Add(len(givers))

	for i := range givers {
		giver := givers[i]
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					giver.OutPort <- giver.Msg
				}
			}
		}()
	}

	wg.Wait()

	return ctx.Err()
}

/* === COMPONENT-RUNNER === */

var (
	ErrRepo          = errors.New("repo")
	ErrComponentFunc = errors.New("operator func")
)

type ComponentRunnerImpl struct {
	repo map[ComponentRef]func(context.Context, IO) error
}

func (c ComponentRunnerImpl) Run(ctx context.Context, components []ComponentRoutine) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := range components {
		component := components[i]

		f, ok := c.repo[component.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrRepo, component.Ref)
		}

		g.Go(func() error {
			if err := f(gctx, component.IO); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrComponentFunc, err), component.Ref)
			}
			return nil
		})
	}

	return g.Wait()
}

/* === CONNECTOR === */

var (
	ErrBroadcast         = errors.New("broadcast")
	ErrDistribute        = errors.New("distribute")
	ErrSelectorSending   = errors.New("selector after sending")
	ErrSelectorReceiving = errors.New("selector before receiving")
)

type ConnectorImlp struct {
	interceptor Interceptor
}

type Interceptor interface {
	AfterSending(from ConnectionSideMeta, msg any) any
	BeforeReceiving(from, to ConnectionSideMeta, msg any) any
	AfterReceiving(from, to ConnectionSideMeta, msg any)
}

func (c ConnectorImlp) Connect(ctx context.Context, net []Connection) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := range net {
		conn := net[i]

		g.Go(func() error {
			if err := c.broadcast(gctx, conn); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrBroadcast, err), conn)
			}
			return nil
		})
	}

	return g.Wait()
}

func (c ConnectorImlp) broadcast(ctx context.Context, conn Connection) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-conn.Sender.Port:
			msg, err = c.applySelector(msg, conn.Sender.Meta.Selectors)
			if err != nil {
				return fmt.Errorf("%w: %v: %v", errors.Join(ErrSelectorSending, err), conn.Sender.Meta, msg)
			}

			msg = c.interceptor.AfterSending(conn.Sender.Meta, msg)

			if err := c.distribute(ctx, msg, conn.Sender.Meta, conn.Receivers); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrDistribute, err), msg)
			}
		}
	}
}

func (c ConnectorImlp) distribute(
	ctx context.Context,
	msg any,
	senderMeta ConnectionSideMeta,
	q []ConnectionSide,
) error {
	i := 0
	processedMessages := make(map[PortAddr]any, len(q)) // intercepted and selected

	for len(q) > 0 {
		recv := q[i]

		if _, ok := processedMessages[recv.Meta.PortAddr]; !ok {
			msg4 := c.interceptor.BeforeReceiving(senderMeta, recv.Meta, msg)
			msg5, err := c.applySelector(msg4, recv.Meta.Selectors)
			if err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrSelectorReceiving, err), recv.Meta)
			}
			processedMessages[recv.Meta.PortAddr] = msg5
		}

		msg5 := processedMessages[recv.Meta.PortAddr]

		select {
		case <-ctx.Done():
			return ctx.Err()
		case recv.Port <- msg5:
			c.interceptor.AfterReceiving(senderMeta, recv.Meta, msg5)
			q = append(q[:i], q[i+1:]...) // remove cur from q
		default: // cur is too busy to receive, it's buf is full
			if i < len(q) {
				i++ // so let's go ask next while it's busy and then return
			}
		}

		if i == len(q) { // end of q, last el was processed (maybe it was busy)
			i = 0 // start over
		}
	}

	return nil
}

func (c ConnectorImlp) applySelector(msg any, selectors []Selector) (any, error) {
	if len(selectors) == 0 {
		return msg, nil
	}

	first := selectors[0]
	if first.RecField != "" {
		// msg = msg.Rec()[first.RecField]
	} else {
		// msg = msg.Arr()[first.ArrIdx]
	}

	return c.applySelector(
		msg,
		selectors[1:],
	)
}
