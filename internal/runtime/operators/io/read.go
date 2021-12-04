package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/emil14/neva/internal/runtime"
)

var ErrRead = errors.New("multiplication")

func Read(io runtime.IO) error {
	out, err := io.Out.Port(runtime.PortAddr{Port: "out"})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWrite, err)
	}

	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			out <- runtime.NewStrMsg(s.Text())
		}
	}()

	return nil
}
