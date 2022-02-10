package connector

import (
	"testing"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

func TestConnector_Connect(t *testing.T) {
	t.Parallel()

	c := Connector{} // TODO

	tests := []struct {
		name        string
		interceptor Interceptor
		net         []runtime.Connection
		nodesIO     map[string]core.IO
	}{
		{
			name:        "",
			interceptor: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c.Connect(tt.nodesIO, tt.net)
		})
	}
}
