package adapter

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type debugInterceptor struct{}

func (p debugInterceptor) Sent(sender, receiver runtime.InterceptorPortAddr, msg runtime.Msg) runtime.Msg {
	fmt.Println("sent:", sender, "->", receiver, msg)
	return nil
}

func (p debugInterceptor) Received(sender, receiver runtime.InterceptorPortAddr, msg runtime.Msg) {
	fmt.Println("received:", sender, "->", receiver, msg)
}

type prodInterceptor struct{}

func (prodInterceptor) Sent(sender, receiver runtime.InterceptorPortAddr, msg runtime.Msg) runtime.Msg {
	return nil
}

func (prodInterceptor) Received(sender, receiver runtime.InterceptorPortAddr, msg runtime.Msg) {}
