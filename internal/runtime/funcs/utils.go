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

func receive2(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
) (runtime.Msg, runtime.Msg, bool) {
	var firstMsg, secondMsg runtime.Msg
	var firstOK, secondOK bool

	var wg sync.WaitGroup
	wg.Go(func() {
		firstMsg, firstOK = firstIn.Receive(ctx)
	})
	wg.Go(func() {
		secondMsg, secondOK = secondIn.Receive(ctx)
	})
	wg.Wait()

	return firstMsg, secondMsg, firstOK && secondOK
}

func receive3(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	thirdIn runtime.SingleInport,
) (runtime.Msg, runtime.Msg, runtime.Msg, bool) {
	var firstMsg, secondMsg, thirdMsg runtime.Msg
	var firstOK, secondOK, thirdOK bool

	var wg sync.WaitGroup
	wg.Go(func() {
		firstMsg, firstOK = firstIn.Receive(ctx)
	})
	wg.Go(func() {
		secondMsg, secondOK = secondIn.Receive(ctx)
	})
	wg.Go(func() {
		thirdMsg, thirdOK = thirdIn.Receive(ctx)
	})
	wg.Wait()

	return firstMsg, secondMsg, thirdMsg, firstOK && secondOK && thirdOK
}

func receive4(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	thirdIn runtime.SingleInport,
	fourthIn runtime.SingleInport,
) (runtime.Msg, runtime.Msg, runtime.Msg, runtime.Msg, bool) {
	var firstMsg, secondMsg, thirdMsg, fourthMsg runtime.Msg
	var firstOK, secondOK, thirdOK, fourthOK bool

	var wg sync.WaitGroup
	wg.Go(func() {
		firstMsg, firstOK = firstIn.Receive(ctx)
	})
	wg.Go(func() {
		secondMsg, secondOK = secondIn.Receive(ctx)
	})
	wg.Go(func() {
		thirdMsg, thirdOK = thirdIn.Receive(ctx)
	})
	wg.Go(func() {
		fourthMsg, fourthOK = fourthIn.Receive(ctx)
	})
	wg.Wait()

	return firstMsg, secondMsg, thirdMsg, fourthMsg, firstOK && secondOK && thirdOK && fourthOK
}
