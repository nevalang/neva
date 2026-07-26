package messages

import "encoding/json"

// --- STRING ---
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type StringMsg struct {
	internalMsg
	v string
}

func (msg StringMsg) Str() string { return msg.v }

func (msg StringMsg) String() string { return msg.v }

func (msg StringMsg) MarshalJSON() ([]byte, error) {
	//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	return json.Marshal(msg.String())
}

func NewStringMsg(s string) StringMsg {
	return StringMsg{
		internalMsg: internalMsg{},
		v:           s,
	}
}

func equalString(left StringMsg, right Msg) bool {
	return equalStringValue(left.v, right)
}

func equalStringValue(left string, right Msg) bool {
	rightTyped, ok := right.(StringMsg)
	return ok && left == rightTyped.v
}
