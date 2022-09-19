package runtime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type (
	Decoder interface {
		Decode([]byte) (Program, error)
	}
	PortGenerator interface {
		Ports([]AbsolutePortAddr) map[AbsolutePortAddr]chan core.Msg
	}
	// TODO maybe move mapping up to runtime?
	ConstSpawner interface {
		Spawn(context.Context, map[AbsolutePortAddr]Msg, map[AbsolutePortAddr]chan core.Msg) error
	}
	OperatorSpawner interface {
		Spawn(context.Context, []Operator, map[AbsolutePortAddr]chan core.Msg) error
	}
	Connector interface {
		Connect(context.Context, map[AbsolutePortAddr]chan core.Msg, []Connection) error
	}
)

var (
	ErrDecoder           = errors.New("program decoder")
	ErrOpSpawner         = errors.New("operator-node spawner")
	ErrConstSpawner      = errors.New("const spawner")
	ErrConnector         = errors.New("connector")
	ErrStartPortNotFound = errors.New("start port not found")
	ErrStartPortBlocked  = errors.New("start port blocked")
)

type Runtime struct {
	decoder      Decoder
	portGen      PortGenerator
	opSpawner    OperatorSpawner
	constSpawner ConstSpawner
	connector    Connector
}

// Run blocks until context is closed or error occurs
func (r Runtime) Run(ctx context.Context, bb []byte) error {
	prog, err := r.decoder.Decode(bb)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDecoder, err)
	}

	ports := r.portGen.Ports(prog.Ports)

	startPort, ok := ports[prog.StartPort]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPort)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := r.connector.Connect(ctx, ports, prog.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.opSpawner.Spawn(ctx, prog.Effects.Operators, ports); err != nil {
			fmt.Println(err) // FIXME get rid
			return fmt.Errorf("%w: %v", ErrOpSpawner, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.constSpawner.Spawn(ctx, prog.Effects.Constants, ports); err != nil {
			return fmt.Errorf("%w: %v", ErrConstSpawner, err)
		}
		return nil
	})

	g.Go(func() error {
		select {
		case startPort <- core.NewDictMsg(nil):
			return nil
		case <-time.After(time.Second): // FIXME deadline for all 3 jobs +  independent from program size
			return ErrStartPortBlocked
		}
	})

	if err := g.Wait(); err != nil { // fixme all goroutines must respect context cancelation
		return fmt.Errorf("wait group: %w", err)
	}

	return nil
}

func MustNew(
	decoder Decoder,
	portGen PortGenerator,
	opSpawner OperatorSpawner,
	constSpawner ConstSpawner,
	connector Connector,
) Runtime {
	utils.PanicOnNil(decoder, portGen, opSpawner, constSpawner, connector)

	return Runtime{
		decoder:      decoder,
		portGen:      portGen,
		opSpawner:    opSpawner,
		constSpawner: constSpawner,
		connector:    connector,
	}
}
