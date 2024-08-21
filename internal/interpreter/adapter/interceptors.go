package adapter

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type debugInterceptor struct{}

func (p debugInterceptor) Sent(sender runtime.InterceptorPortAddr, msg runtime.Msg) runtime.Msg {
	fmt.Println("sent from:", sender, msg)
	return msg
}

func (p debugInterceptor) Received(receiver runtime.InterceptorPortAddr, msg runtime.Msg) runtime.Msg {
	fmt.Println("received to:", receiver, msg)
	return msg
}

type prodInterceptor struct{}

func (prodInterceptor) Sent(sender runtime.InterceptorPortAddr, msg runtime.Msg) runtime.Msg {
	return msg
}

func (prodInterceptor) Received(receiver runtime.InterceptorPortAddr, msg runtime.Msg) runtime.Msg {
	return msg
}
