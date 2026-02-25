package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type osGetenv struct{}

// Create creates runtime function for os.Getenv wrapper.
func (osGetenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	keyIn, err := rio.In.Single("key")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			keyMsg, ok := keyIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(os.Getenv(keyMsg.Str()))) {
				return
			}
		}
	}, nil
}

type osLookupEnv struct{}

// Create creates runtime function for os.LookupEnv wrapper.
func (osLookupEnv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	keyIn, err := rio.In.Single("key")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			keyMsg, ok := keyIn.Receive(ctx)
			if !ok {
				return
			}

			value, exists := os.LookupEnv(keyMsg.Str())
			if !resOut.Send(ctx, lookupEnvResultMsg(value, exists)) {
				return
			}
		}
	}, nil
}

type osSetenv struct{}

// Create creates runtime function for os.Setenv wrapper.
func (osSetenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	keyIn, err := rio.In.Single("key")
	if err != nil {
		return nil, err
	}

	valueIn, err := rio.In.Single("value")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			keyMsg, ok := keyIn.Receive(ctx)
			if !ok {
				return
			}

			valueMsg, ok := valueIn.Receive(ctx)
			if !ok {
				return
			}

			if err := os.Setenv(keyMsg.Str(), valueMsg.Str()); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}

type osUnsetenv struct{}

// Create creates runtime function for os.Unsetenv wrapper.
func (osUnsetenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	keyIn, err := rio.In.Single("key")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			keyMsg, ok := keyIn.Receive(ctx)
			if !ok {
				return
			}

			if err := os.Unsetenv(keyMsg.Str()); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}

type osClearenv struct{}

// Create creates runtime function for os.Clearenv wrapper.
func (osClearenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			os.Clearenv()
			if !resOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}

type osExpandEnv struct{}

// Create creates runtime function for os.ExpandEnv wrapper.
func (osExpandEnv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := rio.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			expanded := os.ExpandEnv(dataMsg.Str())
			if !resOut.Send(ctx, runtime.NewStringMsg(expanded)) {
				return
			}
		}
	}, nil
}

// lookupEnvResultMsg builds std/os.LookupEnvResult payload.
func lookupEnvResultMsg(value string, exists bool) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("value", runtime.NewStringMsg(value)),
		runtime.NewStructField("exists", runtime.NewBoolMsg(exists)),
	})
}
