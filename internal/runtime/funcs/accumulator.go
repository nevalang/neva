package funcs

import (
	"context"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type accumulator struct{}

func (a accumulator) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	initIn, err := io.In.Single("init")
	if err != nil {
		return nil, err
	}

	updIn, err := io.In.Single("upd")
	if err != nil {
		return nil, err
	}

	flushIn, err := io.In.Single("flush")
	if err != nil {
		return nil, err
	}

	curOut, err := io.Out.Single("cur")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
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
