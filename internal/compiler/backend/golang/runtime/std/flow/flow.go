package flow

import (
	"context"

	"github.com/emil14/neva/internal/compiler/backend/golang/runtime"
)

func Void(ctx context.Context, io runtime.FuncIO) error {
	for {
		for _, inports := range io.In {
			for _, inport := range inports {
				<-inport
			}
		}
	}
}

func Trigger(ctx context.Context, io runtime.FuncIO) error {
	sigs, err := io.In.ArrPort("sigs")
	if err != nil {
		return err
	}

	vin, err := io.In.Port("v")
	if err != nil {
		return err
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			for i := range sigs {
				select {
				case <-ctx.Done():
					return nil
				case <-sigs[i]:
					msg := <-vin
					vout <- msg
				}
			}
		}
	}
}
