package funcs

import (
	"context"
	"errors"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type parseInt struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p parseInt) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	baseIn, err := io.In.Single("base")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	bitsIn, err := io.In.Single("bits")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, baseMsg, bitsMsg, ok := receive3(ctx, dataIn, baseIn, bitsIn)
			if !ok {
				return
			}

			parsedNum, err := messages.IntFromStringBase(dataMsg.Msg, baseMsg.Msg, bitsMsg.Msg)
			if err != nil {
				err = errors.New(strings.TrimPrefix(err.Error(), "strconv.Atoi: "))
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, parsedNum) {
				return
			}
		}
	}, nil
}
