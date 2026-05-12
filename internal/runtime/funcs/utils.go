package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

func errFromErr(err error) runtime.StructMsg {
	return errFromString(err.Error())
}

func errFromString(s string) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("text", runtime.NewStringMsg(s)),
		runtime.NewStructField("child", runtime.NewUnionMsg("None", nil)),
	})
}

func streamItem(data runtime.Msg, idx int64, last bool) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("data", data),
		runtime.NewStructField("idx", runtime.NewIntMsg(idx)),
		runtime.NewStructField("last", runtime.NewBoolMsg(last)),
	})
}

func emptyStruct() runtime.StructMsg {
	return runtime.NewStructMsg(nil)
}

func receive2Ordered(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
) (runtime.OrderedMsg, runtime.OrderedMsg, bool) {
	var firstMsg, secondMsg runtime.OrderedMsg
	var firstOK, secondOK bool

	var waitGroup sync.WaitGroup
	waitGroup.Go(func() {
		firstMsg, firstOK = firstIn.ReceiveOrdered(ctx)
	})
	waitGroup.Go(func() {
		secondMsg, secondOK = secondIn.ReceiveOrdered(ctx)
	})
	waitGroup.Wait()

	return firstMsg, secondMsg, firstOK && secondOK
}

//nolint:ireturn // runtime.Msg is the runtime contract type for function ports.
func receive2(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
) (runtime.Msg, runtime.Msg, bool) {
	firstMsg, secondMsg, ok := receive2Ordered(ctx, firstIn, secondIn)
	return firstMsg.Msg, secondMsg.Msg, ok
}

//nolint:ireturn // runtime.Msg is the runtime contract type for function ports.
func receive3(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	thirdIn runtime.SingleInport,
) (runtime.Msg, runtime.Msg, runtime.Msg, bool) {
	var firstMsg, secondMsg, thirdMsg runtime.OrderedMsg
	var firstOK, secondOK, thirdOK bool

	var waitGroup sync.WaitGroup
	waitGroup.Go(func() {
		firstMsg, firstOK = firstIn.ReceiveOrdered(ctx)
	})
	waitGroup.Go(func() {
		secondMsg, secondOK = secondIn.ReceiveOrdered(ctx)
	})
	waitGroup.Go(func() {
		thirdMsg, thirdOK = thirdIn.ReceiveOrdered(ctx)
	})
	waitGroup.Wait()

	return firstMsg.Msg, secondMsg.Msg, thirdMsg.Msg, firstOK && secondOK && thirdOK
}

//nolint:ireturn // runtime.Msg is the runtime contract type for function ports.
func receive4(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	thirdIn runtime.SingleInport,
	fourthIn runtime.SingleInport,
) (runtime.Msg, runtime.Msg, runtime.Msg, runtime.Msg, bool) {
	var firstMsg, secondMsg, thirdMsg, fourthMsg runtime.OrderedMsg
	var firstOK, secondOK, thirdOK, fourthOK bool

	var waitGroup sync.WaitGroup
	waitGroup.Go(func() {
		firstMsg, firstOK = firstIn.ReceiveOrdered(ctx)
	})
	waitGroup.Go(func() {
		secondMsg, secondOK = secondIn.ReceiveOrdered(ctx)
	})
	waitGroup.Go(func() {
		thirdMsg, thirdOK = thirdIn.ReceiveOrdered(ctx)
	})
	waitGroup.Go(func() {
		fourthMsg, fourthOK = fourthIn.ReceiveOrdered(ctx)
	})
	waitGroup.Wait()

	return firstMsg.Msg, secondMsg.Msg, thirdMsg.Msg, fourthMsg.Msg, firstOK && secondOK && thirdOK && fourthOK
}
