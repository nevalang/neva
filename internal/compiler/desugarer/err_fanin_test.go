package desugarer

import (
	"testing"

	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	"github.com/stretchr/testify/require"
)

func TestInsertErrFanInIfNeeded(t *testing.T) {
	t.Run("single_err_connection_no_fanin", func(t *testing.T) {
		d := Desugarer{}
		nodes := map[string]src.Node{}
		conns := []src.Connection{
			errOutConn("node1", "err"),
		}

		got, err := d.insertErrFanInIfNeeded(conns, nodes, core.Meta{})
		require.NoError(t, err)
		require.Len(t, nodes, 0)
		require.Len(t, got, 1)
		require.True(t, isErrOutportConnection(*got[0].Normal))
	})

	t.Run("multiple_err_connections_insert_fanin", func(t *testing.T) {
		d := Desugarer{}
		nodes := map[string]src.Node{}
		conns := []src.Connection{
			errOutConn("node1", "err"),
			errOutConn("node2", "err"),
		}

		got, err := d.insertErrFanInIfNeeded(conns, nodes, core.Meta{})
		require.NoError(t, err)
		require.Contains(t, nodes, "__err_fan_in__1")

		fanInToErr := 0
		fanInInputs := 0
		for _, conn := range got {
			if conn.Normal == nil {
				continue
			}
			receiver := conn.Normal.Receivers[0].PortAddr
			if receiver.Node == "out" && receiver.Port == "err" {
				fanInToErr++
				require.Equal(t, "__err_fan_in__1", conn.Normal.Senders[0].PortAddr.Node)
				continue
			}
			if receiver.Node == "__err_fan_in__1" && receiver.Port == "data" {
				fanInInputs++
			}
		}

		require.Equal(t, 1, fanInToErr)
		require.Equal(t, 2, fanInInputs)
	})
}

func errOutConn(senderNode, senderPort string) src.Connection {
	return src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: senderNode,
						Port: senderPort,
					},
				},
			},
			Receivers: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: "out",
						Port: "err",
					},
				},
			},
		},
	}
}
