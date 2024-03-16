package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type match struct{}

func (match) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	caseIn, ok := io.In["case"]
	if !ok {
		return nil, errors.New("port 'case' is required")
	}

	thenOut, ok := io.Out["then"]
	if !ok {
		return nil, errors.New("port 'then' is required")
	}

	elseOut, err := io.Out.Port("else")
	if err != nil {
		return nil, err
	}

	if len(caseIn) != len(thenOut) {
		return nil, errors.New("number of 'case' inports must match number of 'then' outports")
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case dataMsg := <-dataIn:
				select {
				case <-ctx.Done():
					return
				default:
					matchIdx := -1
					for i, slot := range caseIn {
						select {
						case <-ctx.Done():
							return
						case caseMsg := <-slot:
							if caseMsg == dataMsg {
								matchIdx = i
								break
							}
						}
					}
					if matchIdx != -1 {
						select {
						case <-ctx.Done():
							return
						case thenOut[matchIdx] <- dataMsg:
						}
					} else {
						select {
						case <-ctx.Done():
							return
						case elseOut <- dataMsg:
						}
					}
				}
			}
		}
	}, nil
}
