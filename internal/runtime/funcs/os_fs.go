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
	return createBinaryLoop(rio, "path", "perm", func(pathMsg, permMsg runtime.OrderedMsg) (runtime.Msg, error) {
		mode, err := fileModeFromRuntimeMsg(permMsg)
		if err != nil {
			return nil, err
		}

		if err := os.Mkdir(pathMsg.Str(), mode); err != nil {
			return nil, fmt.Errorf("os.Mkdir: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osMkdirAll struct{}

// Create creates runtime function for os.MkdirAll wrapper.
func (osMkdirAll) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "path", "perm", func(pathMsg, permMsg runtime.OrderedMsg) (runtime.Msg, error) {
		mode, err := fileModeFromRuntimeMsg(permMsg)
		if err != nil {
			return nil, err
		}

		if err := os.MkdirAll(pathMsg.Str(), mode); err != nil {
			return nil, fmt.Errorf("os.MkdirAll: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osReadDir struct{}

// Create creates runtime function for os.ReadDir wrapper.
func (osReadDir) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (runtime.Msg, error) {
		entries, err := os.ReadDir(pathMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.ReadDir: %w", err)
		}

		return dirEntriesMsg(entries), nil
	})
}

type osRemove struct{}

// Create creates runtime function for os.Remove wrapper.
func (osRemove) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (runtime.Msg, error) {
		if err := os.Remove(pathMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.Remove: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osRemoveAll struct{}

// Create creates runtime function for os.RemoveAll wrapper.
func (osRemoveAll) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (runtime.Msg, error) {
		if err := os.RemoveAll(pathMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.RemoveAll: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osRename struct{}

// Create creates runtime function for os.Rename wrapper.
func (osRename) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "oldPath", "newPath", func(oldPathMsg, newPathMsg runtime.OrderedMsg) (runtime.Msg, error) {
		if err := os.Rename(oldPathMsg.Str(), newPathMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.Rename: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osStat struct{}

// Create creates runtime function for os.Stat wrapper.
func (osStat) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (runtime.Msg, error) {
		info, err := os.Stat(pathMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.Stat: %w", err)
		}

		return fileInfoMsg(info), nil
	})
}

type osLstat struct{}

// Create creates runtime function for os.Lstat wrapper.
func (osLstat) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (runtime.Msg, error) {
		info, err := os.Lstat(pathMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.Lstat: %w", err)
		}

		return fileInfoMsg(info), nil
	})
}

type osTruncate struct{}

// Create creates runtime function for os.Truncate wrapper.
func (osTruncate) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "path", "size", func(pathMsg, sizeMsg runtime.OrderedMsg) (runtime.Msg, error) {
		if err := os.Truncate(pathMsg.Str(), sizeMsg.Int()); err != nil {
			return nil, fmt.Errorf("os.Truncate: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osTempDir struct{}

// Create creates runtime function for os.TempDir wrapper.
func (osTempDir) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, false, func() (runtime.Msg, error) {
		return runtime.NewStringMsg(os.TempDir()), nil
	})
}

type osMkdirTemp struct{}

// Create creates runtime function for os.MkdirTemp wrapper.
func (osMkdirTemp) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "dir", "pattern", func(dirMsg, patternMsg runtime.OrderedMsg) (runtime.Msg, error) {
		path, err := os.MkdirTemp(dirMsg.Str(), patternMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.MkdirTemp: %w", err)
		}

		return runtime.NewStringMsg(path), nil
	})
}

type osCreateTemp struct{}

// Create creates runtime function for os.CreateTemp wrapper.
func (osCreateTemp) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "dir", "pattern", func(dirMsg, patternMsg runtime.OrderedMsg) (runtime.Msg, error) {
		file, err := os.CreateTemp(dirMsg.Str(), patternMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.CreateTemp: %w", err)
		}

		fileName := file.Name()
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("close temp file: %w", err)
		}

		return runtime.NewStringMsg(fileName), nil
	})
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
