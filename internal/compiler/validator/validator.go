package validator

import "github.com/emil14/stream/internal/compiler/program"

type validator struct{}

func (v validator) Validate(mod program.Module) error {
	return nil
}

func New() validator {
	return validator{}
}

func MustNew() validator {
	return New()
}
