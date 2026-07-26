package runtime

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// OrderedMsg is a transport envelope with payload and runtime ordering metadata.
type OrderedMsg struct {
	messages.Msg
	index uint64
}

// String formats the payload without exposing transport metadata.
func (o OrderedMsg) String() string { return fmt.Sprint(o.Msg) }
