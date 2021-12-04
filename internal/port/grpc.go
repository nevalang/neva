package port

import (
	"github.com/emil14/neva/pkg/sdk"
)

type Server struct {
	sdk.UnimplementedDevServerServer
}

func NewServer() Server {
	return Server{}
}
