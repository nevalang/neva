package funcs

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// streamTag identifies the union tag used by the stream protocol.
type streamTag uint8

const (
	streamTagOpen streamTag = iota
	streamTagData
	streamTagClose
)

func (tag streamTag) String() string {
	switch tag {
	case streamTagOpen:
		return "Open"
	case streamTagData:
		return "Data"
	case streamTagClose:
		return "Close"
	default:
		return fmt.Sprintf("streamTag(%d)", tag)
	}
}

// newStreamOpenMsg creates the Open union event for a stream.
func newStreamOpenMsg() messages.UnionMsg {
	return messages.NewUnionMsg(streamTagOpen.String(), nil)
}

// newStreamDataMsg creates the Data union event carrying a stream item payload.
func newStreamDataMsg(data messages.Msg) messages.UnionMsg {
	return messages.NewUnionMsg(streamTagData.String(), data)
}

// newStreamCloseMsg creates the Close union event for a stream.
func newStreamCloseMsg() messages.UnionMsg {
	return messages.NewUnionMsg(streamTagClose.String(), nil)
}

func isStreamOpen(msg messages.Msg) bool {
	return hasStreamTag(msg, streamTagOpen)
}

func isStreamData(msg messages.Msg) bool {
	return hasStreamTag(msg, streamTagData)
}

func isStreamClose(msg messages.Msg) bool {
	return hasStreamTag(msg, streamTagClose)
}

// streamDataValue returns the payload of a stream Data event.
//
//nolint:ireturn // Stream payloads are Msg values by runtime contract.
func streamDataValue(msg messages.Msg) messages.Msg {
	unionMsg, ok := messages.AsUnion(msg)
	if !ok {
		panic(fmt.Sprintf("runtime: expected stream union message, got %T", msg))
	}
	if unionMsg.Tag() != streamTagData.String() {
		panic("runtime: expected stream Data message, got " + unionMsg.Tag())
	}
	if unionMsg.Data() == nil {
		panic("runtime: expected stream Data message payload")
	}

	return unionMsg.Data()
}

func hasStreamTag(msg messages.Msg, tag streamTag) bool {
	unionMsg, ok := messages.AsUnion(msg)
	return ok && unionMsg.Tag() == tag.String()
}
