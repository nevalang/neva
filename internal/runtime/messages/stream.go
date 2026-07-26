package messages

const (
	streamOpenTag  = "Open"
	streamDataTag  = "Data"
	streamCloseTag = "Close"
)

// NewStreamOpenMsg creates the opening event for a stream batch.
func NewStreamOpenMsg() UnionMsg {
	return NewUnionMsg(streamOpenTag, nil)
}

// NewStreamDataMsg creates a stream event carrying data.
func NewStreamDataMsg(data Msg) UnionMsg {
	return NewUnionMsg(streamDataTag, data)
}

// NewStreamCloseMsg creates the closing event for a stream batch.
func NewStreamCloseMsg() UnionMsg {
	return NewUnionMsg(streamCloseTag, nil)
}

// IsStreamOpen reports whether msg is an opening stream event.
func IsStreamOpen(msg Msg) bool {
	return hasStreamTag(msg, streamOpenTag)
}

// IsStreamData reports whether msg is a data-carrying stream event.
func IsStreamData(msg Msg) bool {
	return hasStreamTag(msg, streamDataTag)
}

// IsStreamClose reports whether msg is a closing stream event.
func IsStreamClose(msg Msg) bool {
	return hasStreamTag(msg, streamCloseTag)
}

// StreamDataValue returns the payload of a data-carrying stream event.
// It panics when msg is not a valid internal stream data event.
//
//nolint:ireturn // Stream payloads are Msg values by runtime contract.
func StreamDataValue(msg Msg) Msg {
	unionMsg, ok := AsUnion(msg)
	if !ok || unionMsg.Tag() != streamDataTag || unionMsg.Data() == nil {
		panic("messages: expected stream Data event with payload")
	}
	return unionMsg.Data()
}

func hasStreamTag(msg Msg, tag string) bool {
	unionMsg, ok := AsUnion(msg)
	return ok && unionMsg.Tag() == tag
}
