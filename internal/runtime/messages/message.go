package messages

type Msg interface {
	Bool() bool
	Int() int64
	Float() float64
	Str() string
	Bytes() []byte
	List() ListMsg
	Dict() DictMsg
	Struct() StructMsg
	Union() UnionMsg
}

// Internal

type internalMsg struct {
}

func (internalMsg) String() string { panic("unexpected String method call on internal message type") }
func (internalMsg) Bool() bool     { panic("unexpected Bool method call on internal message type") }
func (internalMsg) Int() int64     { panic("unexpected Int method call on internal message type") }
func (internalMsg) Float() float64 { panic("unexpected Float method call on internal message type") }
func (internalMsg) Str() string    { panic("unexpected Str method call on internal message type") }
func (internalMsg) Bytes() []byte  { panic("unexpected Bytes method call on internal message type") }

//nolint:ireturn // Msg contract uses interfaces.
func (internalMsg) List() ListMsg { panic("unexpected List method call on internal message type") }

//nolint:ireturn // Msg contract uses interfaces.
func (internalMsg) Dict() DictMsg {
	panic("unexpected Dict method call on internal message type")
}
func (internalMsg) Struct() StructMsg {
	panic("unexpected Struct method call on internal message type")
}
func (internalMsg) Union() UnionMsg { panic("unexpected Union method call on internal message type") }
func (internalMsg) MarshalJSON() ([]byte, error) {
	panic("unexpected MarshalJSON method call on internal message type")
}
