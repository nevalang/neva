package analyzer

import (
	"fmt"

	src "github.com/nevalang/neva/pkg/sourcecode"
)

// Error is custom error interface implementation that allows to keep track of code location.
// TODO move this to compiler or src package
type Error struct {
	Err      error
	Location *src.Location
	Meta     *src.Meta
}

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

func (e Error) Merge(prefer *Error) *Error {
	if prefer.Err != nil {
		if e.Err == nil {
			e.Err = prefer.Err
		} else {
			e.Err = fmt.Errorf("%w: %v", e.Err, prefer.Err) // FIXME duplication of context
		}
	}

	if prefer.Meta != nil {
		e.Meta = prefer.Meta
	}

	if prefer.Location != nil { //nolint:nestif
		if e.Location == nil {
			e.Location = prefer.Location
		} else {
			if prefer.Location.ModRef.Path != "" {
				e.Location.ModRef.Path = prefer.Location.ModRef.Path
				e.Location.ModRef.Version = prefer.Location.ModRef.Version
			}
			if prefer.Location.PkgName != "" {
				e.Location.PkgName = prefer.Location.PkgName
			}
			if prefer.Location.FileName != "" {
				e.Location.FileName = prefer.Location.FileName
			}
		}
	}

	return &e
}
