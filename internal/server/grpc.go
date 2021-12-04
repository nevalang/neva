package port

import (
	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/sdk"
)

type Server struct {
	compiler compiler.Compiler
	runtime  runtime.Runtime

	sdk.UnimplementedDevServer
}

func NewServer() Server {
	return Server{}
}
