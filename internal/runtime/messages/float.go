package messages

import "fmt"

// Float

type FloatMsg struct {
	internalMsg
	v float64
}

func (msg FloatMsg) Float() float64               { return msg.v }
func (msg FloatMsg) String() string               { return fmt.Sprint(msg.v) }
func (msg FloatMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewFloatMsg(n float64) FloatMsg {
	return FloatMsg{
		internalMsg: internalMsg{},
		v:           n,
	}
}

func equalFloat(left FloatMsg, right Msg) bool {
	rightTyped, ok := right.(FloatMsg)
	return ok && left.v == rightTyped.v
}
