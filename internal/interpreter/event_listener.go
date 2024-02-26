package interpreter

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type DebugEventListener struct{}

func (e DebugEventListener) Send(event runtime.Event, msg runtime.Msg) runtime.Msg {
	if event.Type != runtime.MessagePendingEvent {
		fmt.Println(event, msg)
	}
	return msg
}
