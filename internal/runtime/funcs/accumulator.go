package funcs

import (
	"context"
	"fmt"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type accumulator struct{}

//nolint:cyclop,funlen,gocognit,gocyclo // The stream cycle and its cancellation paths must remain together.
func (a accumulator) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	initIn, err := runtimeIO.In.Single("init")
	if err != nil {
		return nil, fmt.Errorf("get init inport: %w", err)
	}

	updIn, err := runtimeIO.In.Single("upd")
	if err != nil {
		return nil, fmt.Errorf("get upd inport: %w", err)
	}

	flushIn, err := runtimeIO.In.Single("flush")
	if err != nil {
		return nil, fmt.Errorf("get flush inport: %w", err)
	}

	curOut, err := runtimeIO.Out.Single("cur")
	if err != nil {
		return nil, fmt.Errorf("get cur outport: %w", err)
	}

	resOut, err := runtimeIO.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("get res outport: %w", err)
	}

	return func(ctx context.Context) {
		for {
			var acc runtime.Msg

			initMsg, initOk := initIn.Receive(ctx)
			if !initOk {
				return
			}

			if !curOut.Send(ctx, initMsg) {
				return
			}

			acc = initMsg

			cycleCtx, cancel := context.WithCancel(ctx)
			updCh := make(chan runtime.Msg)
			flushCh := make(chan runtime.Msg)

			go func() {
				defer close(updCh)
				for {
					msg, ok := updIn.Receive(cycleCtx)
					if !ok {
						return
					}
					select {
					case <-cycleCtx.Done():
						return
					case updCh <- msg:
					}
				}
			}()

			go func() {
				defer close(flushCh)
				for {
					msg, ok := flushIn.Receive(cycleCtx)
					if !ok {
						return
					}
					select {
					case <-cycleCtx.Done():
						return
					case flushCh <- msg:
					}
				}
			}()

			cycleDone := false
			for !cycleDone {
				select {
				case <-ctx.Done():
					cancel()
					return
				case dataMsg, ok := <-updCh:
					if !ok {
						cancel()
						return
					}
					if !curOut.Send(ctx, dataMsg) {
						cancel()
						return
					}
					acc = dataMsg
				case flushMsg, ok := <-flushCh:
					if !ok {
						cancel()
						return
					}
					if !flushMsg.Bool() {
						continue
					}

					// Close can arrive before the last reduced value is delivered to upd.
					// Wait for a short quiet period and consume trailing upd messages.
					drainTimer := time.NewTimer(time.Millisecond)
					draining := true
					for draining {
						select {
						case <-ctx.Done():
							if !drainTimer.Stop() {
								<-drainTimer.C
							}
							cancel()
							return
						case dataMsg, ok := <-updCh:
							if !ok {
								draining = false
								break
							}
							if !curOut.Send(ctx, dataMsg) {
								if !drainTimer.Stop() {
									<-drainTimer.C
								}
								cancel()
								return
							}
							acc = dataMsg
							if !drainTimer.Stop() {
								<-drainTimer.C
							}
							drainTimer.Reset(time.Millisecond)
						case <-drainTimer.C:
							draining = false
						}
					}

					if !resOut.Send(ctx, acc) {
						cancel()
						return
					}
					cycleDone = true
				}
			}
			cancel()
		}
	}, nil
}
