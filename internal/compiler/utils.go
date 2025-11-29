package compiler

//go:generate neva build --target=go --target-go-mode=pkg --output=utils/generated utils

// Pointer allows to avoid creating of temporary variables just to take pointers.
func Pointer[T any](v T) *T {
	return &v
}
