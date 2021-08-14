package server

import (
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/program"
)

type Server struct {
	runtime runtime.Runtime
}

func (srv Server) RunProgram(p program.Program, in map[program.PortAddr]runtime.Msg) (runtime.Msg, error) {
	// io, err := srv.runtime.Run(p)
	// if err != nil {
	// 	return runtime.Msg{}, err
	// }

	// for addr, msg := range in {
	// 	io.In
	// }

	return nil, nil // TODO
}

func New(srv Server)
