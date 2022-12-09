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

// resolvable
type TypeExpr struct {
	Application TypeApplication     // not for struct type (indirect recursion!)
	StructDef   map[string]TypeExpr // only for struct type (direct recursion)
}

// Applyes type to its arguments
type TypeApplication struct { // list<list<int>>
	ref  string
	args []TypeExpr // can contain refs to generics, can be empty (indirect recursion!)
}

func resolve(expr TypeExpr, scope map[string]TypeDef) (TypeExpr, error) { // Add support for structs
	refType, ok := scope[expr.Application.ref] // check that reference type exists
	if !ok {
		return TypeExpr{}, errors.New("")
	}

	// check that generic args for every param is present
	if len(refType.params) > len(expr.Application.args) { // compare equality? structural typing? linting?
		return TypeExpr{}, errors.New("")
	}

	resolvedArgs := make([]TypeExpr, 0, len(refType.params))
	newScope := make(map[string]TypeDef, len(scope)+len(refType.params)) // new scope contains resolved params (shadow)
	// optimized for concurrency (is there better way?)
	maps.Copy(newScope, scope)
	for i, param := range refType.params {
		resolvedArg, err := resolve(expr.Application.args[i], scope)
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

	baseType, ok := scope[refType.typeExpr.Application.ref]
	if !ok {
		return TypeExpr{}, errors.New("")
	}
	if expr.Application.ref == baseType.typeExpr.Application.ref {
		return TypeExpr{
			Application: TypeApplication{
				ref:  refType.typeExpr.Application.ref,
				args: resolvedArgs,
			},
			StructDef: nil, // todo
		}, nil
	}

	return resolve(refType.typeExpr, newScope)
}

func main() {
	// test1()
	test2()
}

func test2() {
	scope := map[string]TypeDef{ // int = int, list<t> = list
		"int": {
			typeExpr: TypeExpr{
				Application: TypeApplication{ref: "int"}, // native types references themselves
			},
		},
		"list": {
			typeExpr: TypeExpr{
				Application: TypeApplication{ref: "list"}, // native types references themselves  (params?)
			},
			params: []string{"t"},
		},
		"custom": { // custom<t> = { x: list<t> }
			params: []string{"t"},
			typeExpr: TypeExpr{
				StructDef: map[string]TypeExpr{
					"x": {
						Application: TypeApplication{
							ref: "list",
							args: []TypeExpr{
								{
									Application: TypeApplication{ref: "t"}, // ref to param
								},
							},
						},
					},
				},
			},
		},
	}

	expr := TypeExpr{ // custom<int> -> list<list<int>>
		Application: TypeApplication{
			ref: "custom",
			args: []TypeExpr{
				{Application: TypeApplication{ref: "int"}},
			},
		},
	}

	got, err := resolve(expr, scope)
	if err != nil {
		panic(err)
	}

	want := TypeExpr{
		Application: TypeApplication{
			ref: "list",
			args: []TypeExpr{
				{
					Application: TypeApplication{
						ref: "list",
						args: []TypeExpr{
							{
								Application: TypeApplication{ref: "int"},
							},
						},
					},
				},
			},
		},
	}

	if fmt.Sprint(got.Application) != fmt.Sprint(want.Application) {
		panic("not equal")
	}

	fmt.Println("Got: ", got, "Want: ", want)
}

func test1() {
	scope := map[string]TypeDef{ // int = int, list<t> = list
		"int": {
			typeExpr: TypeExpr{
				Application: TypeApplication{ref: "int"}, // native types references themselves
			},
		},
		"list": {
			typeExpr: TypeExpr{
				Application: TypeApplication{ref: "list"}, // native types references themselves  (params?)
			},
			params: []string{"t"},
		},
		"custom": { // custom<t> = list<list<t>>
			params: []string{"t"},
			typeExpr: TypeExpr{
				Application: TypeApplication{
					ref: "list",
					args: []TypeExpr{
						{
							Application: TypeApplication{
								ref: "list",
								args: []TypeExpr{
									{
										Application: TypeApplication{ref: "t"}, // from params
									},
								},
							},
						},
					},
				},
			},
		},
	}

	expr := TypeExpr{ // custom<int> -> list<list<int>>
		Application: TypeApplication{
			ref: "custom",
			args: []TypeExpr{
				{Application: TypeApplication{ref: "int"}},
			},
		},
	}

	got, err := resolve(expr, scope)
	if err != nil {
		panic(err)
	}

	want := TypeExpr{
		Application: TypeApplication{
			ref: "list",
			args: []TypeExpr{
				{
					Application: TypeApplication{
						ref: "list",
						args: []TypeExpr{
							{
								Application: TypeApplication{ref: "int"},
							},
						},
					},
				},
			},
		},
	}

	if fmt.Sprint(got.Application) != fmt.Sprint(want.Application) {
		panic("not equal")
	}

	fmt.Println("Got: ", got, "Want: ", want)
}
