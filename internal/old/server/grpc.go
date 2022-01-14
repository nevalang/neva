package server

import (
	"context"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/sdk"
)

type Server struct {
	compiler compiler.Compiler
	runtime  runtime.Runtime

	workdir string
	q       chan runtime.Msg
}

func (s Server) ListPrograms(context.Context, *sdk.ListProgramsRequest) (*sdk.ListProgramsResponse, error) {
	return nil, nil
}

func (s Server) GetProgram(context.Context, *sdk.GetProgramRequest) (*sdk.GetProgramResponse, error) {
	return nil, nil
}

func (s Server) UpdateProgram(context.Context, *sdk.UpdateProgramRequest) (*sdk.UpdateProgramResponse, error) {
	return nil, nil
}

func (s Server) StartDebugger(*sdk.StartDebugRequest, sdk.Dev_StartDebuggerServer) error {
	return nil
}

func (s Server) SendDebugMessage(req *sdk.DebugRequest, stream sdk.Dev_SendDebugMessageServer) error {
	return stream.Send(&sdk.DebugResponse{})
}

func NewServer() Server {
	return Server{}
}
