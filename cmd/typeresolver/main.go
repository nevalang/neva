package main

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

// TypeDef with empty typeExpr and empty structFields is a native builtin type
type TypeDef struct { // l<t> = list<t> || l<t> = { foo t }
	params       []string            // any type can have it
	typeExpr     TypeExpr            // not for struct type
	structFields map[string]TypeExpr // only for struct type
}

type TypeExpr struct { // list<list<int>>
	ref  string
	args []TypeExpr // can contain refs to generics
}

func resolve(expr TypeExpr, scope map[string]TypeDef) (TypeExpr, error) { // Add support for structs
	refType, ok := scope[expr.ref] // check that reference type exists
	if !ok {
		return TypeExpr{}, errors.New("")
	}

	// check that generic args for every param is present
	if len(refType.params) > len(expr.args) { // compare equality? structural typing? linting?
		return TypeExpr{}, errors.New("")
	}

	resolvedArgs := make([]TypeExpr, 0, len(refType.params))
	newScope := make(map[string]TypeDef, len(scope)+len(refType.params)) // new scope contains resolved params (shadow)
	// optimized for concurrency (is there better way?)
	maps.Copy(newScope, scope)
	for i, param := range refType.params {
		resolvedArg, err := resolve(expr.args[i], scope)
		if err != nil {
			return TypeExpr{}, errors.New("")
		}
		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param] = TypeDef{
			params:       nil, // we don't refer generics with another generics inside!
			typeExpr:     resolvedArg,
			structFields: nil, // type argument can only be an expression
		}
	}

	baseType, ok := scope[refType.typeExpr.ref]
	if !ok {
		return TypeExpr{}, errors.New("")
	}
	if expr.ref == baseType.typeExpr.ref {
		return TypeExpr{
			ref:  refType.typeExpr.ref,
			args: resolvedArgs,
		}, nil
	}

	return resolve(refType.typeExpr, newScope)
}

func main() {
	resolved, err := resolve(TypeExpr{ // custom<int> -> list<list<int>>
		ref: "custom",
		args: []TypeExpr{
			{ref: "int"},
		},
	}, map[string]TypeDef{
		"int": {
			typeExpr: TypeExpr{ref: "int"}, // base types references themselves
		},
		"list": {
			typeExpr: TypeExpr{ref: "list"}, // base types references themselves (params?)
			params:   []string{"t"},
		},
		"custom": { // custom<t> = list<list<t>>
			params: []string{"t"},
			typeExpr: TypeExpr{
				ref: "list",
				args: []TypeExpr{
					{
						ref: "list",
						args: []TypeExpr{
							{ref: "t"}, // from params
						},
					},
				},
			},
			structFields: map[string]TypeExpr{},
		},
	})
	if err != nil {
		panic(err)
	}

	expected := TypeExpr{
		ref: "list",
		args: []TypeExpr{
			{
				ref: "list",
				args: []TypeExpr{
					{ref: "int"},
				},
			},
		},
	}

	fmt.Println("GOT", resolved)
	fmt.Println("WANT", expected)
}
