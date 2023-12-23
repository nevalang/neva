package analyzer

import (
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

// Error is custom error interface implementation that allows to keep track of code location.
// TODO move this to compiler or src package
type Error struct {
	Err      error
	Location *src.Location
	Meta     *src.Meta
}

func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprint(e.Location, e.Meta)
}

func (e Error) Merge(prefer *Error) *Error {
	if prefer.Err != nil {
		if e.Err == nil {
			e.Err = prefer.Err
		} else {
			e.Err = fmt.Errorf("%w: %v", e.Err, prefer.Err)
		}
	}

	if prefer.Meta != nil {
		e.Meta = prefer.Meta
	}

	if prefer.Location != nil { //nolint:nestif
		if e.Location == nil {
			e.Location = prefer.Location
		} else {
			if prefer.Location.ModRef.Name != "" {
				e.Location.ModRef.Name = prefer.Location.ModRef.Name
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
