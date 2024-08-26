package adapter

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type debugInterceptor struct{}

func (d debugInterceptor) Sent(sender runtime.PortSlotAddr, msg runtime.Msg) runtime.Msg {
	fmt.Println("sent from:", d.formatPortSlotAddr(sender), d.formatMsg(msg))
	return msg
}

func (d debugInterceptor) Received(receiver runtime.PortSlotAddr, msg runtime.Msg) runtime.Msg {
	fmt.Println("received to:", d.formatPortSlotAddr(receiver), d.formatMsg(msg))
	return msg
}

func (d debugInterceptor) formatMsg(msg runtime.Msg) string {
	if s, ok := msg.(runtime.StrMsg); ok {
		return fmt.Sprintf(`"%s"`, s)
	}
	return fmt.Sprint(msg)
}

func (d debugInterceptor) formatPortSlotAddr(slotAddr runtime.PortSlotAddr) string {
	s := fmt.Sprintf("%v:%v", slotAddr.Path, slotAddr.Port)
	if slotAddr.Index != nil {
		s = fmt.Sprintf("%v[%v]", s, *slotAddr.Index)
	}
	return s
}

type prodInterceptor struct{}

func (prodInterceptor) Sent(sender runtime.PortSlotAddr, msg runtime.Msg) runtime.Msg {
	return msg
}

func (prodInterceptor) Received(receiver runtime.PortSlotAddr, msg runtime.Msg) runtime.Msg {
	return msg
}
