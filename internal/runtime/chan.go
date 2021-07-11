package runtime

type Chan map[string]struct {
	in  chan Msg
	out chan Msg
}
