package compiler

import (
	"context"
	"fmt"

	neva "github.com/nevalang/neva/internal/compiler/utils/generated"
	"github.com/nevalang/neva/pkg/core"
)

//go:generate neva build --target=go --target-go-mode=pkg --target-go-runtime-path=../runtime --output=utils/generated utils

// Pointer allows to avoid creating of temporary variables just to take pointers.
func Pointer[T any](v T) *T {
	return &v
}

// ParseEntityRef calls Neva and marshals result into core.EntityRef.
func ParseEntityRef(ctx context.Context, ref string) (core.EntityRef, error) {
	// Call the generated Neva function
	out, err := neva.ParseEntityRef(ctx, neva.ParseEntityRefInput{Ref: ref})
	if err != nil {
		return core.EntityRef{}, err
	}

	// Unmarshal the result
	if !out.Res.IsStruct() {
		return core.EntityRef{}, fmt.Errorf("expected struct msg, got %T", out.Res)
	}
	msg := out.Res.Struct()

	return core.EntityRef{
		Pkg:  msg.Get("pkg").Str(),
		Name: msg.Get("name").Str(),
		Meta: core.Meta{Text: msg.Get("metaText").Str()},
	}, nil
}
