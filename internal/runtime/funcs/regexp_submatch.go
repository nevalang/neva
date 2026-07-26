package funcs

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type regexpSubmatch struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (r regexpSubmatch) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	regexpIn, err := io.In.Single("regexp")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	dataIn, err := io.In.Single("data")
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
			regexpMsg, dataMsg, ok := receive2(ctx, regexpIn, dataIn)
			if !ok {
				return
			}

			regex, err := regexp.Compile(regexpMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, messages.NewStringMsg(err.Error())) {
					return
				}
				continue
			}

			if !resOut.Send(
				ctx,
				stringsToList(
					regex.FindStringSubmatch(
						fmt.Sprint(dataMsg),
					),
				),
			) {
				return
			}
		}
	}, nil
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func stringsToList(ss []string) messages.Msg {
	msgs := make([]messages.Msg, 0, len(ss))
	for _, s := range ss {
		msgs = append(msgs, messages.NewStringMsg(s))
	}
	return messages.NewListMsg(msgs)
}
