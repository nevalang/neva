package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/emil14/neva/internal/core"
)

func Println(ctx context.Context, io core.IO) error {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-dataIn:
			fmt.Println(msg)
			select {
			case dataOut <- msg:
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func Readln(ctx context.Context, io core.IO) error {
	sig, err := io.In.Port("sig")
	if err != nil {
		return err
	}

	data, err := io.Out.Port("data")
	if err != nil {
		return err
	}

	s := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sig:
			s.Scan()
			select {
			case data <- core.NewStrMsg(s.Text()):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
