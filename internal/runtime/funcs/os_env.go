package funcs

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type osGetenv struct{}

// Create creates runtime function for os.Getenv wrapper.
func (osGetenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "key", false, func(keyMsg runtime.OrderedMsg) (runtime.Msg, error) {
		return runtime.NewStringMsg(os.Getenv(keyMsg.Str())), nil
	})
}

type osLookupEnv struct{}

// Create creates runtime function for os.LookupEnv wrapper.
func (osLookupEnv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "key", false, func(keyMsg runtime.OrderedMsg) (runtime.Msg, error) {
		value, exists := os.LookupEnv(keyMsg.Str())
		return lookupEnvResultMsg(value, exists), nil
	})
}

type osSetenv struct{}

// Create creates runtime function for os.Setenv wrapper.
func (osSetenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "key", "value", func(keyMsg, valueMsg runtime.OrderedMsg) (runtime.Msg, error) {
		if err := os.Setenv(keyMsg.Str(), valueMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.Setenv: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osUnsetenv struct{}

// Create creates runtime function for os.Unsetenv wrapper.
func (osUnsetenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "key", true, func(keyMsg runtime.OrderedMsg) (runtime.Msg, error) {
		if err := os.Unsetenv(keyMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.Unsetenv: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osClearenv struct{}

// Create creates runtime function for os.Clearenv wrapper.
func (osClearenv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, false, func() (runtime.Msg, error) {
		os.Clearenv()
		return emptyStruct(), nil
	})
}

type osExpandEnv struct{}

// Create creates runtime function for os.ExpandEnv wrapper.
func (osExpandEnv) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "data", false, func(dataMsg runtime.OrderedMsg) (runtime.Msg, error) {
		return runtime.NewStringMsg(os.ExpandEnv(dataMsg.Str())), nil
	})
}

// lookupEnvResultMsg builds std/os.LookupEnvResult payload.
func lookupEnvResultMsg(value string, exists bool) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("value", runtime.NewStringMsg(value)),
		runtime.NewStructField("exists", runtime.NewBoolMsg(exists)),
	})
}
