package runtime

import "fmt"

// StreamTag identifies the runtime representation of a stream event.
type StreamTag uint8

const (
	StreamTagOpen StreamTag = iota
	StreamTagData
	StreamTagClose
)

// String returns the Neva union tag name for a stream event tag.
func (tag StreamTag) String() string {
	switch tag {
	case StreamTagOpen:
		return "Open"
	case StreamTagData:
		return "Data"
	case StreamTagClose:
		return "Close"
	default:
		return fmt.Sprintf("StreamTag(%d)", tag)
	}
}

// NewStreamOpenMsg creates the Open event for a stream.
func NewStreamOpenMsg() UnionMsg {
	return NewUnionMsg(StreamTagOpen.String(), nil)
}

// NewStreamDataMsg creates a Data event carrying a stream item payload.
func NewStreamDataMsg(data Msg) UnionMsg {
	return NewUnionMsg(StreamTagData.String(), data)
}

// NewStreamCloseMsg creates the Close event for a stream.
func NewStreamCloseMsg() UnionMsg {
	return NewUnionMsg(StreamTagClose.String(), nil)
}

// IsStreamOpen reports whether msg is a stream Open event.
func IsStreamOpen(msg Msg) bool {
	return hasStreamTag(msg, StreamTagOpen)
}

// IsStreamData reports whether msg is a stream Data event.
func IsStreamData(msg Msg) bool {
	return hasStreamTag(msg, StreamTagData)
}

// IsStreamClose reports whether msg is a stream Close event.
func IsStreamClose(msg Msg) bool {
	return hasStreamTag(msg, StreamTagClose)
}

// StreamDataValue returns the payload of a stream Data event.
//
//nolint:ireturn // Stream payloads are Msg values by runtime contract.
func StreamDataValue(msg Msg) Msg {
	unionMsg, ok := AsUnion(msg)
	if !ok {
		panic(fmt.Sprintf("runtime: expected stream union message, got %T", msg))
	}
	if unionMsg.Tag() != StreamTagData.String() {
		panic("runtime: expected stream Data message, got " + unionMsg.Tag())
	}
	if unionMsg.Data() == nil {
		panic("runtime: expected stream Data message payload")
	}

	return unionMsg.Data()
}

func hasStreamTag(msg Msg, tag StreamTag) bool {
	unionMsg, ok := AsUnion(msg)
	return ok && unionMsg.Tag() == tag.String()
}
