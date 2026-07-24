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
			dataMsg, sepMsg, received := receive2(ctx, dataIn, sepIn)
			if !received {
				return
			}

			builder := strings.Builder{}
			sep := sepMsg.Str()

			list := dataMsg.List()
			if stringsList, ok := runtime.AsListStrings(list); ok {
				for i := range stringsList {
					if i > 0 {
						builder.WriteString(sep)
					}
					builder.WriteString(stringsList[i])
				}
			} else {
				writeJoinedList(&builder, list.Msgs(), sep)
			}

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
		var msg runtime.OrderedMsg
		if !hasSep {
			sepMsg, dataMsg, received := receive2(ctx, sepIn, dataIn)
			if !received {
				return
			}

			sep = sepMsg.Str()
			msg = dataMsg
			hasSep = true
		} else {
			dataMsg, dataReceived := dataIn.Receive(ctx)
			if !dataReceived {
				return
			}
			msg = dataMsg
		}

		nextHasSep, keepRunning := handleJoinedStreamMessage(ctx, &builder, resOut, msg.Msg, sep, hasSep)
		if !keepRunning {
			return
		}

		hasSep = nextHasSep
	}
}
