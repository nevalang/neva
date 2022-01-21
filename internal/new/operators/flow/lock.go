package main

import (
	"fmt"

	"github.com/emil14/neva/internal/new/core"
)

func Lock(io core.IO) error {
	sig, ok := io.In[core.PortAddr{Port: "sig"}]
	if !ok {
		return fmt.Errorf("%w: in: %v", core.ErrPortNotFound, "sig")
	}

	data, ok := io.In[core.PortAddr{Port: "data"}]
	if !ok {
		return fmt.Errorf("%w: in: %v", core.ErrPortNotFound, "data")
	}

	out, ok := io.Out[core.PortAddr{Port: "out"}]
	if !ok {
		return fmt.Errorf("%w: out: %v", core.ErrPortNotFound, "out")
	}

	go func() {
		for msg := range data {
			<-sig
			out <- msg
		}
	}()

	return nil
}
