package flow

import (
	"context"

	"github.com/emil14/neva/internal/compiler/backend/golang/runtime"
)

func Void(ctx context.Context, io runtime.FuncIO) error {
	for {
		for _, portSlots := range io.In {
			for _, slot := range portSlots {
				select {
				case <-ctx.Done():
					return nil
				case <-slot:
				}
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
		for i := range sigs {
			select {
			case <-ctx.Done():
				return nil
			case <-sigs[i]:
			}
		}

		select {
		case <-ctx.Done():
			return nil
		case msg := <-vin:
			select {
			case <-ctx.Done():
				return nil
			case vout <- msg:
			}
		}
	}
}
