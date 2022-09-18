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
		BeforeReceiving(
			sender, receiver runtime.AbsolutePortAddr,
			point runtime.ReceiverConnectionPoint,
			msg core.Msg,
		) core.Msg
		AfterReceiving(
			sender, receiver runtime.AbsolutePortAddr,
			point runtime.ReceiverConnectionPoint,
			msg core.Msg,
		)
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
			if err := c.connect(ctx, v); err != nil {
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

func (c Connector) connect(ctx context.Context, connection ConnectionWithChans) error {
	semaphore := make(chan struct{}, len(connection.receivers))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-connection.sender:
			msg = c.interceptor.AfterSending(connection.info, msg)
			if err := c.broadcast(connection, msg, semaphore); err != nil {
				return fmt.Errorf("broadcast: %w", err)
			}
		}
	}
}

// FIXME https://github.com/emil14/neva/issues/86
func (c Connector) broadcast(connection ConnectionWithChans, msg core.Msg, guard chan struct{}) error {
	for i := range connection.receivers {
		receiverPortChan := connection.receivers[i]
		receiverConnPoint := connection.info.ReceiversConnectionPoints[i] // we believe mapper

		if receiverConnPoint.Type == runtime.DictKeyReading {
			path := receiverConnPoint.DictReadingPath
			for _, part := range path[:len(path)-1] {
				var ok bool
				msg, ok = msg.Dict()[part]
				if !ok {
					return fmt.Errorf("%w: ", ErrDictKeyNotFound)
				}
			}
		}

		guard <- struct{}{}

		go func(m core.Msg) {
			receiverPortChan <- c.interceptor.BeforeReceiving(
				connection.info.SenderPortAddr, receiverConnPoint.PortAddr, receiverConnPoint, m,
			)
			c.interceptor.AfterReceiving(connection.info.SenderPortAddr, receiverConnPoint.PortAddr, receiverConnPoint, m)
			<-guard
		}(msg)
	}

	return nil
}

func MustNew(i Interceptor) Connector {
	utils.PanicOnNil(i)

	return Connector{
		interceptor: i,
		mapper:      mapper{},
	}
}
