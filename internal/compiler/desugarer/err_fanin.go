package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
)

// insertErrFanInIfNeeded inserts an intermediate FanIn for :err when multiple senders target it.
func (d *Desugarer) insertErrFanInIfNeeded(
	connections []src.Connection,
	nodes map[string]src.Node,
	componentMeta core.Meta,
) ([]src.Connection, error) {
	var errConns []src.Connection
	keptConns := make([]src.Connection, 0, len(connections))

	for _, conn := range connections {
		if conn.Normal == nil {
			keptConns = append(keptConns, conn)
			continue
		}

		if isErrOutportConnection(*conn.Normal) {
			errConns = append(errConns, conn)
			continue
		}

		keptConns = append(keptConns, conn)
	}

	totalIncoming := 0
	for _, conn := range errConns {
		totalIncoming += len(conn.Normal.Senders)
	}

	if totalIncoming <= 1 {
		return connections, nil
	}

	d.fanInCounter++
	fanInName := fmt.Sprintf("__err_fan_in__%d", d.fanInCounter)
	locOnlyMeta := componentMeta
	if len(errConns) > 0 {
		locOnlyMeta = core.Meta{Location: errConns[0].Meta.Location}
	}

	nodes[fanInName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "FanIn",
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	idx := uint8(0)
	for _, conn := range errConns {
		if len(conn.Normal.Receivers) != 1 {
			return nil, fmt.Errorf("err outport connection must have one receiver")
		}

		for _, sender := range conn.Normal.Senders {
			keptConns = append(keptConns, src.Connection{
				Normal: &src.NormalConnection{
					Senders: []src.ConnectionSender{sender},
					Receivers: []src.ConnectionReceiver{
						{
							PortAddr: &src.PortAddr{
								Node: fanInName,
								Port: "data",
								Idx:  compiler.Pointer(idx),
								Meta: locOnlyMeta,
							},
							Meta: locOnlyMeta,
						},
					},
					Meta: locOnlyMeta,
				},
				Meta: locOnlyMeta,
			})
			idx++
		}
	}

	keptConns = append(keptConns, src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: fanInName,
						Port: "res",
						Meta: locOnlyMeta,
					},
					Meta: locOnlyMeta,
				},
			},
			Receivers: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: "out",
						Port: "err",
						Meta: locOnlyMeta,
					},
					Meta: locOnlyMeta,
				},
			},
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	})

	return keptConns, nil
}

func isErrOutportConnection(conn src.NormalConnection) bool {
	if len(conn.Receivers) != 1 {
		return false
	}
	receiver := conn.Receivers[0]
	if receiver.PortAddr == nil {
		return false
	}
	return receiver.PortAddr.Node == "out" &&
		receiver.PortAddr.Port == "err" &&
		receiver.PortAddr.Idx == nil
}
