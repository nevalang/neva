package compiler

import (
	"fmt"

	src "github.com/nevalang/neva/pkg/sourcecode"
)

// Error represents end-user error that preserves not only message but also location in sourcecode where error occurred.
// It must be properly used by all compiler's dependencies so user face human-readable error.
type Error struct {
	Err      error
	Location *src.Location
	Meta     *src.Meta
}

// Error simply implements error interface.
func (e Error) Error() string {
	switch {
	case e.Location != nil && e.Meta != nil:
		return fmt.Sprintf("%v:%v:%v %v", *e.Location, e.Meta.Start.Column, e.Meta.Start.Column, e.Err)
	case e.Meta != nil:
		return fmt.Sprintf("%v:%v %v", e.Meta.Start.Column, e.Meta.Start.Column, e.Err)
	case e.Location != nil:
		return fmt.Sprintf("%v: %v", *e.Location, e.Err)
	case e.Err != nil:
		return e.Err.Error()
	}
	return ""
}

// Merge merges e (parent) Error with the child Error, every data child has is preferred to the parent's.
// Parent provides broad location while child provide more accurate.
// Ideally there's no need to wrap errors the old fashioned way by imitating stack-trace. Location should provide that.
func (e Error) Merge(child *Error) *Error {
	if child.Err != nil {
		if e.Err == nil {
			e.Err = child.Err
		} else {
			e.Err = fmt.Errorf("%w: %v", e.Err, child.Err) // FIXME duplication of context
		}
	}

	if child.Meta != nil {
		e.Meta = child.Meta
	}

	if child.Location != nil { //nolint:nestif
		if e.Location == nil {
			e.Location = child.Location
		} else {
			if child.Location.ModRef.Path != "" {
				e.Location.ModRef.Path = child.Location.ModRef.Path
				e.Location.ModRef.Version = child.Location.ModRef.Version
			}
			if child.Location.PkgName != "" {
				e.Location.PkgName = child.Location.PkgName
			}
			if child.Location.FileName != "" {
				e.Location.FileName = child.Location.FileName
			}
		}
	}

	return &e
}
