package funcs

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type osMkdir struct{}

// Create creates runtime function for os.Mkdir wrapper.
func (osMkdir) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
	if err != nil {
		return nil, err
	}

	permIn, err := rio.In.Single("perm")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			permMsg, ok := permIn.Receive(ctx)
			if !ok {
				return
			}

			mode, err := fileModeFromRuntimeMsg(permMsg)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if err := os.Mkdir(pathMsg.Str(), mode); err != nil {
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

type osMkdirAll struct{}

// Create creates runtime function for os.MkdirAll wrapper.
func (osMkdirAll) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
	if err != nil {
		return nil, err
	}

	permIn, err := rio.In.Single("perm")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			permMsg, ok := permIn.Receive(ctx)
			if !ok {
				return
			}

			mode, err := fileModeFromRuntimeMsg(permMsg)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if err := os.MkdirAll(pathMsg.Str(), mode); err != nil {
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

type osReadDir struct{}

// Create creates runtime function for os.ReadDir wrapper.
func (osReadDir) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			entries, err := os.ReadDir(pathMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, dirEntriesMsg(entries)) {
				return
			}
		}
	}, nil
}

type osRemove struct{}

// Create creates runtime function for os.Remove wrapper.
func (osRemove) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			if err := os.Remove(pathMsg.Str()); err != nil {
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

type osRemoveAll struct{}

// Create creates runtime function for os.RemoveAll wrapper.
func (osRemoveAll) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			if err := os.RemoveAll(pathMsg.Str()); err != nil {
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

type osRename struct{}

// Create creates runtime function for os.Rename wrapper.
func (osRename) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	oldPathIn, err := rio.In.Single("oldPath")
	if err != nil {
		return nil, err
	}

	newPathIn, err := rio.In.Single("newPath")
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
			oldPathMsg, ok := oldPathIn.Receive(ctx)
			if !ok {
				return
			}

			newPathMsg, ok := newPathIn.Receive(ctx)
			if !ok {
				return
			}

			if err := os.Rename(oldPathMsg.Str(), newPathMsg.Str()); err != nil {
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

type osStat struct{}

// Create creates runtime function for os.Stat wrapper.
func (osStat) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			info, err := os.Stat(pathMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, fileInfoMsg(info)) {
				return
			}
		}
	}, nil
}

type osLstat struct{}

// Create creates runtime function for os.Lstat wrapper.
func (osLstat) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			info, err := os.Lstat(pathMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, fileInfoMsg(info)) {
				return
			}
		}
	}, nil
}

type osTruncate struct{}

// Create creates runtime function for os.Truncate wrapper.
func (osTruncate) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
	if err != nil {
		return nil, err
	}

	sizeIn, err := rio.In.Single("size")
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
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			sizeMsg, ok := sizeIn.Receive(ctx)
			if !ok {
				return
			}

			if err := os.Truncate(pathMsg.Str(), sizeMsg.Int()); err != nil {
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

type osTempDir struct{}

// Create creates runtime function for os.TempDir wrapper.
func (osTempDir) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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

			if !resOut.Send(ctx, runtime.NewStringMsg(os.TempDir())) {
				return
			}
		}
	}, nil
}

type osMkdirTemp struct{}

// Create creates runtime function for os.MkdirTemp wrapper.
func (osMkdirTemp) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dirIn, err := rio.In.Single("dir")
	if err != nil {
		return nil, err
	}

	patternIn, err := rio.In.Single("pattern")
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
			dirMsg, ok := dirIn.Receive(ctx)
			if !ok {
				return
			}

			patternMsg, ok := patternIn.Receive(ctx)
			if !ok {
				return
			}

			path, err := os.MkdirTemp(dirMsg.Str(), patternMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(path)) {
				return
			}
		}
	}, nil
}

type osCreateTemp struct{}

// Create creates runtime function for os.CreateTemp wrapper.
func (osCreateTemp) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dirIn, err := rio.In.Single("dir")
	if err != nil {
		return nil, err
	}

	patternIn, err := rio.In.Single("pattern")
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
			dirMsg, ok := dirIn.Receive(ctx)
			if !ok {
				return
			}

			patternMsg, ok := patternIn.Receive(ctx)
			if !ok {
				return
			}

			file, err := os.CreateTemp(dirMsg.Str(), patternMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			fileName := file.Name()
			if err := file.Close(); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(fileName)) {
				return
			}
		}
	}, nil
}

const maxUint32AsInt64 = int64(^uint32(0))

// fileModeFromRuntimeMsg parses os.FileMode from a runtime integer message.
func fileModeFromRuntimeMsg(permMsg runtime.Msg) (os.FileMode, error) {
	perm := permMsg.Int()
	if perm < 0 || perm > maxUint32AsInt64 {
		return 0, fmt.Errorf("permission value out of range: %d", perm)
	}

	// #nosec G115 -- bounds checked above.
	return os.FileMode(uint32(perm)), nil
}

// dirEntriesMsg converts []os.DirEntry to list<DirEntry> runtime payload.
func dirEntriesMsg(entries []os.DirEntry) runtime.ListMsg {
	msgs := make([]runtime.Msg, len(entries))
	for i := range entries {
		msgs[i] = runtime.NewStructMsg([]runtime.StructField{
			runtime.NewStructField("name", runtime.NewStringMsg(entries[i].Name())),
			runtime.NewStructField("isDir", runtime.NewBoolMsg(entries[i].IsDir())),
		})
	}

	return runtime.NewListMsg(msgs)
}

// fileInfoMsg converts os.FileInfo to std/os.FileInfo runtime payload.
func fileInfoMsg(info os.FileInfo) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("name", runtime.NewStringMsg(info.Name())),
		runtime.NewStructField("size", runtime.NewIntMsg(info.Size())),
		runtime.NewStructField("mode", runtime.NewIntMsg(int64(info.Mode()))),
		runtime.NewStructField("modTimeUnix", runtime.NewIntMsg(info.ModTime().Unix())),
		runtime.NewStructField("isDir", runtime.NewBoolMsg(info.IsDir())),
	})
}
