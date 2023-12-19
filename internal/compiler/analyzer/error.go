package analyzer

import (
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Error struct {
	Err      error
	Location *src.Location
	Meta     *src.Meta
}

func (e Error) Error() string {
	return fmt.Sprintf("%v %v %v", e.Location, e.Meta, e.Err)
}

func (e Error) Wrap(err error) *Error {
	return &Error{
		Err:      fmt.Errorf("%w: %v", err, e.Err),
		Location: e.Location,
		Meta:     e.Meta,
	}
}

func (e Error) Merge(prefer *Error) *Error {
	if prefer.Err != nil {
		e.Err = prefer.Err
	}

	if prefer.Meta != nil {
		e.Meta = prefer.Meta
	}

	if prefer.Location != nil { //nolint:nestif
		if e.Location == nil {
			e.Location = prefer.Location
		} else {
			if prefer.Location.ModuleName != "" {
				e.Location.ModuleName = prefer.Location.ModuleName
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
