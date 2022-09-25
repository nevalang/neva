package connector

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"golang.org/x/sync/errgroup"
)

type (
	Interceptor interface {
		AfterSending(runtime.Connection, core.Msg) core.Msg
		BeforeReceiving(saddr runtime.AbsolutePortAddr, point runtime.ReceiverConnectionPoint, msg core.Msg) core.Msg
		AfterReceiving(saddr runtime.AbsolutePortAddr, point runtime.ReceiverConnectionPoint, msg core.Msg)
	}

	Mapper interface {
		MapPortsToConnections(map[runtime.AbsolutePortAddr]chan core.Msg, []runtime.Connection) ([]ConnectionWithChans, error)
	}

	ConnectionWithChans struct {
		sender    chan core.Msg
		receivers []chan core.Msg
		info      runtime.Connection
	}
)

var (
	ErrMapper          = errors.New("mapper")
	ErrDictKeyNotFound = errors.New("dict key not found")
)

type Connector struct {
	mapper      Mapper
	interceptor Interceptor
}

func (c Connector) Connect(
	ctx context.Context,
	ports map[runtime.AbsolutePortAddr]chan core.Msg,
	connections []runtime.Connection,
) error {
	connectionsWithChans, err := c.mapper.MapPortsToConnections(ports, connections)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMapper, err)
	}

	g := errgroup.Group{}
	for i := range connectionsWithChans {
		v := connectionsWithChans[i]
		g.Go(func() error {
			if err := c.handleConnection(ctx, v); err != nil {
				return fmt.Errorf("connect: %w", err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("errgroup wait: %w", err)
	}

	return nil
}

func (c Connector) handleConnection(ctx context.Context, connection ConnectionWithChans) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-connection.sender:
			msg = c.interceptor.AfterSending(connection.info, msg)
			if err := c.broadcast(connection, msg); err != nil {
				return fmt.Errorf("broadcast: %w", err)
			}
		}
	}
}

// FIXME https://github.com/emil14/neva/issues/86
func (c Connector) broadcast(connection ConnectionWithChans, msg core.Msg) error {
	for i := range connection.receivers {
		receiverPort := connection.receivers[i]
		receiverPoint := connection.info.ReceiversConnectionPoints[i] // we believe mapper

		if receiverPoint.Type == runtime.DictKeyReading {
			path := receiverPoint.DictReadingPath
			for _, part := range path[:len(path)-1] {
				var ok bool
				msg, ok = msg.Dict()[part]
				if !ok {
					return fmt.Errorf("%w: ", ErrDictKeyNotFound)
				}
			}
		}

		receiverPort <- c.interceptor.BeforeReceiving(
			connection.info.SenderPortAddr, receiverPoint.PortAddr, receiverPoint, msg,
		)
		c.interceptor.AfterReceiving(connection.info.SenderPortAddr, receiverPoint.PortAddr, receiverPoint, msg)
	}

	return nil
}

func MustNew(i Interceptor) Connector {
	utils.NilPanic(i)

	return Connector{
		interceptor: i,
		mapper:      mapper{},
	}
}
