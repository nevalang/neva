package compiler

// Pointer allows to avoid creating of temporary variables just to take pointers.
func Pointer[T any](v T) *T {
	return &v
}
