package funcs

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type osMkdir struct{}

// Create creates runtime function for os.Mkdir wrapper.
func (osMkdir) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "path", "perm", func(pathMsg, permMsg runtime.OrderedMsg) (messages.Msg, error) {
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
func (osMkdirAll) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "path", "perm", func(pathMsg, permMsg runtime.OrderedMsg) (messages.Msg, error) {
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
func (osReadDir) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (messages.Msg, error) {
		entries, err := os.ReadDir(pathMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.ReadDir: %w", err)
		}

		return messages.NewListMsg(dirEntries(entries)), nil
	})
}

type osRemove struct{}

// Create creates runtime function for os.Remove wrapper.
func (osRemove) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (messages.Msg, error) {
		if err := os.Remove(pathMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.Remove: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osRemoveAll struct{}

// Create creates runtime function for os.RemoveAll wrapper.
func (osRemoveAll) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (messages.Msg, error) {
		if err := os.RemoveAll(pathMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.RemoveAll: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osRename struct{}

// Create creates runtime function for os.Rename wrapper.
func (osRename) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "oldPath", "newPath", func(oldPathMsg, newPathMsg runtime.OrderedMsg) (messages.Msg, error) {
		if err := os.Rename(oldPathMsg.Str(), newPathMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.Rename: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osStat struct{}

// Create creates runtime function for os.Stat wrapper.
func (osStat) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (messages.Msg, error) {
		info, err := os.Stat(pathMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.Stat: %w", err)
		}

		return fileInfoMsg(info), nil
	})
}

type osLstat struct{}

// Create creates runtime function for os.Lstat wrapper.
func (osLstat) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (messages.Msg, error) {
		info, err := os.Lstat(pathMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.Lstat: %w", err)
		}

		return fileInfoMsg(info), nil
	})
}

type osTruncate struct{}

// Create creates runtime function for os.Truncate wrapper.
func (osTruncate) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "path", "size", func(pathMsg, sizeMsg runtime.OrderedMsg) (messages.Msg, error) {
		if err := os.Truncate(pathMsg.Str(), sizeMsg.Int()); err != nil {
			return nil, fmt.Errorf("os.Truncate: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osTempDir struct{}

// Create creates runtime function for os.TempDir wrapper.
func (osTempDir) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, false, func() (messages.Msg, error) {
		return messages.NewStringMsg(os.TempDir()), nil
	})
}

type osMkdirTemp struct{}

// Create creates runtime function for os.MkdirTemp wrapper.
func (osMkdirTemp) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "dir", "pattern", func(dirMsg, patternMsg runtime.OrderedMsg) (messages.Msg, error) {
		path, err := os.MkdirTemp(dirMsg.Str(), patternMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.MkdirTemp: %w", err)
		}

		return messages.NewStringMsg(path), nil
	})
}

type osCreateTemp struct{}

// Create creates runtime function for os.CreateTemp wrapper.
func (osCreateTemp) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	return createBinaryLoop(rio, "dir", "pattern", func(dirMsg, patternMsg runtime.OrderedMsg) (messages.Msg, error) {
		file, err := os.CreateTemp(dirMsg.Str(), patternMsg.Str())
		if err != nil {
			return nil, fmt.Errorf("os.CreateTemp: %w", err)
		}

		fileName := file.Name()
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("close temp file: %w", err)
		}

		return messages.NewStringMsg(fileName), nil
	})
}

const maxUint32AsInt64 = int64(^uint32(0))

// fileModeFromRuntimeMsg parses os.FileMode from a runtime integer message.
func fileModeFromRuntimeMsg(permMsg messages.Msg) (os.FileMode, error) {
	perm := permMsg.Int()
	if perm < 0 || perm > maxUint32AsInt64 {
		return 0, fmt.Errorf("permission value out of range: %d", perm)
	}

	// #nosec G115 -- bounds checked above.
	return os.FileMode(uint32(perm)), nil
}

// dirEntries converts []os.DirEntry into generic runtime struct messages.
func dirEntries(entries []os.DirEntry) []messages.Msg {
	msgs := make([]messages.Msg, len(entries))
	for i := range entries {
		msgs[i] = messages.NewStructMsg([]messages.StructField{
			messages.NewStructField("name", messages.NewStringMsg(entries[i].Name())),
			messages.NewStructField("isDir", messages.NewBoolMsg(entries[i].IsDir())),
		})
	}

	return msgs
}

// fileInfoMsg converts os.FileInfo to std/os.FileInfo runtime payload.
func fileInfoMsg(info os.FileInfo) messages.StructMsg {
	return messages.NewStructMsg([]messages.StructField{
		messages.NewStructField("name", messages.NewStringMsg(info.Name())),
		messages.NewStructField("size", messages.NewIntMsg(info.Size())),
		messages.NewStructField("mode", messages.NewIntMsg(int64(info.Mode()))),
		messages.NewStructField("modTimeUnix", messages.NewIntMsg(info.ModTime().Unix())),
		messages.NewStructField("isDir", messages.NewBoolMsg(info.IsDir())),
	})
}
