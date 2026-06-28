package funcs

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

func TestLookupEnvResultMsg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		value      string
		wantValue  string
		exists     bool
		wantExists bool
	}{
		{name: "exists", value: "value", exists: true, wantValue: "value", wantExists: true},
		{name: "missing", exists: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			msg := lookupEnvResultMsg(tt.value, tt.exists)
			if got := msg.Get("value").Str(); got != tt.wantValue {
				t.Fatalf("value = %q, want %q", got, tt.wantValue)
			}
			if got := msg.Get("exists").Bool(); got != tt.wantExists {
				t.Fatalf("exists = %t, want %t", got, tt.wantExists)
			}
		})
	}
}

func TestOSGetenvRuntimeFunc(t *testing.T) {
	t.Setenv("NEVA_OS_GETENV", "value")

	got := runUnaryRuntimeFunc(t, osGetenv{}, "key", runtime.NewStringMsg("NEVA_OS_GETENV"))
	if got.Str() != "value" {
		t.Fatalf("getenv = %q, want value", got.Str())
	}
}

func TestOSLookupEnvRuntimeFuncReportsMissing(t *testing.T) {
	got := runUnaryRuntimeFunc(t, osLookupEnv{}, "key", runtime.NewStringMsg("NEVA_OS_MISSING"))
	fields := got.Struct()
	if fields.Get("value").Str() != "" || fields.Get("exists").Bool() {
		t.Fatalf("lookup missing = %v, want empty/false", got)
	}
}

func TestOSSetenvRuntimeFuncStoresValue(t *testing.T) {
	runBinaryRuntimeFunc(
		t,
		osSetenv{},
		"key",
		"value",
		runtime.NewStringMsg("NEVA_OS_SETENV"),
		runtime.NewStringMsg("stored"),
	)
	if got := getenv(t, "NEVA_OS_SETENV"); got != "stored" {
		t.Fatalf("setenv stored %q, want stored", got)
	}
}

func TestOSSetenvRuntimeFuncReportsInvalidKey(t *testing.T) {
	got := runBinaryRuntimeFuncErr(
		t,
		osSetenv{},
		"key",
		"value",
		runtime.NewStringMsg("BAD=KEY"),
		runtime.NewStringMsg("value"),
	)
	if !strings.Contains(got.Struct().Get("text").Str(), "os.Setenv") {
		t.Fatalf("setenv error = %v, want os.Setenv", got)
	}
}

func TestOSUnsetenvRuntimeFuncRemovesValue(t *testing.T) {
	t.Setenv("NEVA_OS_UNSETENV", "value")

	runUnaryRuntimeFunc(t, osUnsetenv{}, "key", runtime.NewStringMsg("NEVA_OS_UNSETENV"))
	if _, ok := lookupEnv(t, "NEVA_OS_UNSETENV"); ok {
		t.Fatal("unsetenv left variable set")
	}
}

func TestOSExpandEnvRuntimeFuncSubstitutesValue(t *testing.T) {
	t.Setenv("NEVA_OS_EXPAND", "expanded")

	got := runUnaryRuntimeFunc(t, osExpandEnv{}, "data", runtime.NewStringMsg("$NEVA_OS_EXPAND"))
	if got.Str() != "expanded" {
		t.Fatalf("expand_env = %q, want expanded", got.Str())
	}
}

func TestOSClearenvRuntimeFuncClearsProcessEnv(t *testing.T) {
	env := os.Environ()
	t.Cleanup(func() {
		restoreEnv(t, env)
	})
	setenv(t, "NEVA_OS_CLEARENV", "value")

	runSignalRuntimeFunc(t, osClearenv{})
	if _, ok := lookupEnv(t, "NEVA_OS_CLEARENV"); ok {
		t.Fatal("clearenv left variable set")
	}
}

func runSignalRuntimeFunc(t *testing.T, creator runtime.FuncCreator) runtime.OrderedMsg {
	t.Helper()

	io, inChans, outChans := newIO([]string{"sig"}, []string{"res", "err"})
	handler, err := creator.Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans["sig"] <- runtime.OrderedMsg{Msg: emptyStruct()}
	return receiveRuntimeMsg(t, outChans, "res")
}

func runUnaryRuntimeFunc(
	t *testing.T,
	creator runtime.FuncCreator,
	inName string,
	input runtime.Msg,
) runtime.OrderedMsg {
	t.Helper()

	io, inChans, outChans := newIO([]string{inName}, []string{"res", "err"})
	handler, err := creator.Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans[inName] <- runtime.OrderedMsg{Msg: input}
	return receiveRuntimeMsg(t, outChans, "res")
}

func runUnaryRuntimeFuncErr(
	t *testing.T,
	creator runtime.FuncCreator,
	inName string,
	input runtime.Msg,
) runtime.OrderedMsg {
	t.Helper()

	io, inChans, outChans := newIO([]string{inName}, []string{"res", "err"})
	handler, err := creator.Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans[inName] <- runtime.OrderedMsg{Msg: input}
	return receiveRuntimeMsg(t, outChans, "err")
}

func runBinaryRuntimeFunc(
	t *testing.T,
	creator runtime.FuncCreator,
	firstName string,
	secondName string,
	first runtime.Msg,
	second runtime.Msg,
) runtime.OrderedMsg {
	t.Helper()

	io, inChans, outChans := newIO([]string{firstName, secondName}, []string{"res", "err"})
	handler, err := creator.Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	sendInOrder(t, inChans, []string{secondName, firstName}, map[string]runtime.Msg{
		firstName:  first,
		secondName: second,
	})
	return receiveRuntimeMsg(t, outChans, "res")
}

func runBinaryRuntimeFuncErr(
	t *testing.T,
	creator runtime.FuncCreator,
	firstName string,
	secondName string,
	first runtime.Msg,
	second runtime.Msg,
) runtime.OrderedMsg {
	t.Helper()

	io, inChans, outChans := newIO([]string{firstName, secondName}, []string{"res", "err"})
	handler, err := creator.Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	sendInOrder(t, inChans, []string{secondName, firstName}, map[string]runtime.Msg{
		firstName:  first,
		secondName: second,
	})
	return receiveRuntimeMsg(t, outChans, "err")
}

func receiveRuntimeMsg(
	t *testing.T,
	outChans map[string]chan runtime.OrderedMsg,
	outName string,
) runtime.OrderedMsg {
	t.Helper()

	select {
	case got := <-outChans[outName]:
		return got
	case <-time.After(time.Second):
		t.Fatalf("no output on %s", outName)
		return runtime.OrderedMsg{}
	}
}

func setenv(t *testing.T, key string, value string) {
	t.Helper()

	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Setenv(%q): %v", key, err)
	}
	t.Cleanup(func() {
		_ = os.Unsetenv(key)
	})
}

func getenv(t *testing.T, key string) string {
	t.Helper()

	value, ok := os.LookupEnv(key)
	if !ok {
		t.Fatalf("LookupEnv(%q) missing", key)
	}
	return value
}

func lookupEnv(t *testing.T, key string) (string, bool) {
	t.Helper()

	return os.LookupEnv(key)
}

func restoreEnv(t *testing.T, env []string) {
	t.Helper()

	os.Clearenv()
	for _, item := range env {
		key, value, ok := strings.Cut(item, "=")
		if !ok {
			t.Fatalf("invalid environment entry %q", item)
		}
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("restore env %q: %v", key, err)
		}
	}
}
