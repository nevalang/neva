package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringJoinList struct{}

func (stringJoinList) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	sepIn, err := singleInport(runtimeIO, "sep")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, dataReceived := dataIn.Receive(ctx)
			if !dataReceived {
				return
			}

			sepMsg, sepReceived := sepIn.Receive(ctx)
			if !sepReceived {
				return
			}

			builder := strings.Builder{}
			sep := sepMsg.Str()
			writeJoinedList(&builder, dataMsg.List(), sep)

			if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
				return
			}
		}
	}, nil
}

type stringJoinStream struct{}

func (stringJoinStream) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	sepIn, err := singleInport(runtimeIO, "sep")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		runStringJoinStream(ctx, dataIn, sepIn, resOut)
	}, nil
}

func writeJoinedList(builder *strings.Builder, list []runtime.Msg, sep string) {
	for idx := range list {
		appendStreamItem(builder, list[idx].Str(), sep)
	}
}

func appendStreamItem(builder *strings.Builder, item, sep string) {
	if builder.Len() > 0 {
		builder.WriteString(sep)
	}

	builder.WriteString(item)
}

func receiveJoinSeparator(ctx context.Context, sepIn runtime.SingleInport) (string, bool) {
	sepMsg, sepReceived := sepIn.Receive(ctx)
	if !sepReceived {
		return "", false
	}

	return sepMsg.Str(), true
}

func handleJoinedStreamMessage(
	ctx context.Context,
	builder *strings.Builder,
	resOut runtime.SingleOutport,
	msg runtime.Msg,
	sep string,
	hasSep bool,
) (bool, bool) {
	switch {
	case isStreamOpen(msg):
		builder.Reset()
		return hasSep, true
	case isStreamData(msg):
		appendStreamItem(builder, streamDataValue(msg).Str(), sep)
		return hasSep, true
	case isStreamClose(msg):
		if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
			return false, false
		}

		builder.Reset()
		return false, true
	default:
		panic("strings_join_stream: unexpected stream tag")
	}
}

func runStringJoinStream(
	ctx context.Context,
	dataIn, sepIn runtime.SingleInport,
	resOut runtime.SingleOutport,
) {
	builder := strings.Builder{}
	var (
		sep    string
		hasSep bool
	)

	for {
		msg, dataReceived := dataIn.Receive(ctx)
		if !dataReceived {
			return
		}

		if !hasSep {
			nextSep, sepReceived := receiveJoinSeparator(ctx, sepIn)
			if !sepReceived {
				return
			}

			sep = nextSep
			hasSep = true
		}

		nextHasSep, keepRunning := handleJoinedStreamMessage(ctx, &builder, resOut, msg, sep, hasSep)
		if !keepRunning {
			return
		}

		hasSep = nextHasSep
	}
}
