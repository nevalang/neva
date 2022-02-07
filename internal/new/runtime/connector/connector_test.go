package connector

import (
	"testing"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type dummyInterceptor struct{}

func (d dummyInterceptor) Event(e EventType, c runtime.Connection, m core.Msg) core.Msg {
	return m
}

func TestConnector_Connect(t *testing.T) {
	t.Parallel()

	c := Connector{dummyInterceptor{}}

	type args struct {
		net     []runtime.Connection
		nodesIO map[string]core.IO
		stop    chan<- error
	}

	tests := []struct {
		name        string
		interceptor Interceptor
		args        args
	}{
		{
			name:        "",
			interceptor: nil,
			args:        args{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c.Connect(tt.args.net, tt.args.nodesIO, tt.args.stop)
		})
	}
}
