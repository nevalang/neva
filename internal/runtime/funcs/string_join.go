package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type stringJoinList struct{}

func (stringJoinList) Create(runtimeIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	dataIn, sepIn, resOut, err := resolveStringJoinPorts(runtimeIO)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		runStringJoinList(ctx, dataIn, sepIn, resOut)
	}, nil
}

func resolveStringJoinPorts(runtimeIO runtime.IO) (runtime.SingleInport, runtime.SingleInport, runtime.SingleOutport, error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	sepIn, err := singleInport(runtimeIO, "sep")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	return dataIn, sepIn, resOut, nil
}

func runStringJoinList(
	ctx context.Context,
	dataIn, sepIn runtime.SingleInport,
	resOut runtime.SingleOutport,
) {
	for {
		dataMsg, sepMsg, received := receive2(ctx, dataIn, sepIn)
		if !received {
			return
		}

		result := joinList(dataMsg.List(), sepMsg.Str())
		if !resOut.Send(ctx, messages.NewStringMsg(result)) {
			return
		}
	}
}

func joinList(list messages.ListMsg, sep string) string {
	builder := strings.Builder{}
	if stringsList, ok := messages.ListAsStrings(list); ok {
		for i := range stringsList {
			if i > 0 {
				builder.WriteString(sep)
			}
			builder.WriteString(stringsList[i])
		}
		return builder.String()
	}

	writeJoinedList(&builder, list.Untyped(), sep)
	return builder.String()
}

type stringJoinStream struct{}

func (stringJoinStream) Create(runtimeIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	dataIn, sepIn, resOut, err := resolveStringJoinPorts(runtimeIO)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		runStringJoinStream(ctx, dataIn, sepIn, resOut)
	}, nil
}

func writeJoinedList(builder *strings.Builder, list []messages.Msg, sep string) {
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
	msg messages.Msg,
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
		if !resOut.Send(ctx, messages.NewStringMsg(builder.String())) {
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
