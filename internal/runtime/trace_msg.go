package runtime

import "fmt"

type tracedMsg struct {
	msg           Msg
	traceID       uint64
	parentTraceID uint64
}

func (m tracedMsg) String() string {
	return fmt.Sprint(m.msg)
}

func (m tracedMsg) Bool() bool {
	return m.msg.Bool()
}

func (m tracedMsg) Int() int64 {
	return m.msg.Int()
}

func (m tracedMsg) Float() float64 {
	return m.msg.Float()
}

func (m tracedMsg) Str() string {
	return m.msg.Str()
}

func (m tracedMsg) Bytes() []byte {
	return m.msg.Bytes()
}

func (m tracedMsg) List() []Msg {
	return m.msg.List()
}

func (m tracedMsg) Dict() map[string]Msg {
	return m.msg.Dict()
}

func (m tracedMsg) Struct() StructMsg {
	return m.msg.Struct()
}

func (m tracedMsg) Union() UnionMsg {
	return m.msg.Union()
}

func (m tracedMsg) Equal(other Msg) bool {
	return m.msg.Equal(UnwrapTraceMsg(other))
}

//nolint:ireturn // Msg is runtime contract type.
func withTrace(msg Msg, traceID, parentTraceID uint64) Msg {
	return tracedMsg{
		msg:           UnwrapTraceMsg(msg),
		traceID:       traceID,
		parentTraceID: parentTraceID,
	}
}

func traceMeta(msg Msg) (uint64, uint64, bool) {
	tm, ok := msg.(tracedMsg)
	if !ok {
		return 0, 0, false
	}
	return tm.traceID, tm.parentTraceID, true
}

// UnwrapTraceMsg returns the original runtime message payload.
//
//nolint:ireturn // Msg is runtime contract type.
func UnwrapTraceMsg(msg Msg) Msg {
	cur := msg
	for {
		tm, ok := cur.(tracedMsg)
		if !ok {
			return cur
		}
		cur = tm.msg
	}
}

// TraceIDFromMsg extracts runtime Dataflow Trace identity from a message.
func TraceIDFromMsg(msg Msg) (uint64, bool) {
	traceID, _, ok := traceMeta(msg)
	return traceID, ok
}

func parentTraceIDFromMsg(msg Msg) uint64 {
	_, parentTraceID, ok := traceMeta(msg)
	if !ok {
		return 0
	}
	return parentTraceID
}

// AsUnion unwraps trace metadata and checks whether message is a union.
func AsUnion(msg Msg) (UnionMsg, bool) {
	unionMsg, ok := UnwrapTraceMsg(msg).(UnionMsg)
	return unionMsg, ok
}
