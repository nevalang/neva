package runtime

type traceCarrier interface {
	traceMeta() (uint64, uint64)
}

func mustTraceMeta(msg Msg) (uint64, uint64) {
	carrier, ok := msg.(traceCarrier)
	if !ok {
		panic("runtime invariant: message does not carry trace metadata")
	}
	return carrier.traceMeta()
}

// TraceIDFromMsg extracts runtime Dataflow Trace identity from message payload.
func TraceIDFromMsg(msg Msg) (uint64, bool) {
	carrier, ok := msg.(traceCarrier)
	if !ok {
		return 0, false
	}
	traceID, _ := carrier.traceMeta()
	return traceID, traceID != 0
}

func mustTraceIDFromMsg(msg Msg) uint64 {
	traceID, _ := mustTraceMeta(msg)
	if traceID == 0 {
		panic("runtime invariant: trace id is not initialized")
	}
	return traceID
}

func parentTraceIDFromMsg(msg Msg) uint64 {
	carrier, ok := msg.(traceCarrier)
	if !ok {
		return 0
	}
	_, parentTraceID := carrier.traceMeta()
	return parentTraceID
}

//nolint:ireturn // Msg is runtime contract type.
func withTrace(msg Msg, traceID, parentTraceID uint64) Msg {
	switch typedMsg := msg.(type) {
	case BoolMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case IntMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case FloatMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case StringMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case BytesMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case ListMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case DictMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case StructMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case UnionMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	default:
		return msg
	}
}

func AsUnion(msg Msg) (UnionMsg, bool) {
	unionMsg, ok := msg.(UnionMsg)
	return unionMsg, ok
}
