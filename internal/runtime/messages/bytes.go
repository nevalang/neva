package messages

import (
	"bytes"
	"encoding/json"
)

// --- BYTES ---
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type BytesMsg struct {
	internalMsg
	v []byte
}

func (msg BytesMsg) Bytes() []byte { return msg.v }

func (msg BytesMsg) String() string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg BytesMsg) MarshalJSON() ([]byte, error) {
	//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	return json.Marshal(msg.v)
}

func NewBytesMsg(v []byte) BytesMsg {
	return BytesMsg{
		internalMsg: internalMsg{},
		v:           v,
	}
}

func equalBytes(left BytesMsg, right Msg) bool {
	rightTyped, ok := right.(BytesMsg)
	return ok && bytes.Equal(left.v, rightTyped.v)
}
