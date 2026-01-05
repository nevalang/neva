package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type state struct{}

func (s state) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	initIn, err := io.In.Single("init")
	if err != nil {
		return nil, err
	}
	updIn, err := io.In.Single("upd")
	if err != nil {
		return nil, err
	}
	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var (
			currentState runtime.Msg
			initialized  bool
			buffer       []runtime.Msg
		)

		updCh := make(chan runtime.Msg)
		initCh := make(chan runtime.Msg)

		// Forward 'upd' messages
		go func() {
			defer close(updCh)
			for {
				msg, ok := updIn.Receive(ctx)
				if !ok {
					return
				}
				select {
				case updCh <- msg:
				case <-ctx.Done():
					return
				}
			}
		}()

		// Forward 'init' messages
		go func() {
			defer close(initCh)
			for {
				msg, ok := initIn.Receive(ctx)
				if !ok {
					return
				}
				select {
				case initCh <- msg:
				case <-ctx.Done():
					return
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case initMsg, ok := <-initCh:
				if !ok {
					initCh = nil
					continue
				}

				currentState = initMsg
				initialized = true
				if !resOut.Send(ctx, currentState) {
					return
				}

				// Flush buffer
				if len(buffer) > 0 {
					for _, u := range buffer {
						currentState = u
						if !resOut.Send(ctx, currentState) {
							return
						}
					}
					buffer = nil // Clear buffer
				}

			case updMsg, ok := <-updCh:
				if !ok {
					updCh = nil
					continue
				}

				if !initialized {
					buffer = append(buffer, updMsg)
				} else {
					currentState = updMsg
					if !resOut.Send(ctx, currentState) {
						return
					}
				}
			}

			if initCh == nil && updCh == nil {
				return
			}
		}
	}, nil
}
