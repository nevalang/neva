package main

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

// TypeDef with empty typeExpr and empty structFields is a native builtin TypeDef
type TypeDef struct { // l<t> = list<t> || l<t> = { foo t }
	Params []string // any type can have it
	Type   Type     // not for struct type (just "Type" in go spec)
}

// resolvable
type Type struct { // Instantiation
	Instantiation TypeInstantiation   // not for struct type (indirect recursion!)
	StructLit     map[string]Type     // only for struct type (direct recursion)
	EnumLit       map[string]struct{} // TODO implement
}

type TypeInstantiation struct { // list<list<int>>
	Ref  string
	Args []Type // can contain refs to generics, can be empty (indirect recursion!)
}

func resolve(expr Type, scope map[string]TypeDef) (Type, error) { // Add support for structs
	if expr.StructLit != nil {
		resolvedStruct := make(map[string]Type, len(expr.StructLit))
		for field, expr := range expr.StructLit {
			resolvedFieldExpr, err := resolve(expr, scope)
			if err != nil {
				return Type{}, errors.New("")
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Type{
			StructLit: resolvedStruct,
		}, nil
	}

	refType, ok := scope[expr.Instantiation.Ref] // check that reference type exists
	if !ok {
		return Type{}, errors.New("")
	}

	// check that generic args for every param is present
	if len(refType.Params) > len(expr.Instantiation.Args) { // compare equality? structural typing? linting?
		return Type{}, errors.New("")
	}

	resolvedArgs := make([]Type, 0, len(refType.Params))
	newScope := make(map[string]TypeDef, len(scope)+len(refType.Params)) // new scope contains resolved params (shadow)
	// optimized for concurrency (is there better way?)
	maps.Copy(newScope, scope)
	for i, param := range refType.Params {
		resolvedArg, err := resolve(expr.Instantiation.Args[i], scope)
		if err != nil {
			return Type{}, errors.New("")
		}
		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param] = TypeDef{
			Params: nil, // we don't refer generics with another generics inside!
			Type:   resolvedArg,
		}
	}

	if refType.Type.StructLit == nil { // reference type's body is an application, not a struct definition
		baseType, ok := scope[refType.Type.Instantiation.Ref] // FIXME not work structs
		if !ok {
			return Type{}, errors.New("")
		}
		if expr.Instantiation.Ref == baseType.Type.Instantiation.Ref {
			return Type{
				Instantiation: TypeInstantiation{
					Ref:  refType.Type.Instantiation.Ref,
					Args: resolvedArgs,
				},
				StructLit: nil, // todo
			}, nil
		}
	}

	return resolve(refType.Type, newScope) // if it's not a native type and not a struct, then do recursive
}

func main() {
	test3()
	// test2()
	// test1()
}

func test3() {
	scope := map[string]TypeDef{ // int = int, list<t> = list
		"int": {
			Type: Type{
				Instantiation: TypeInstantiation{Ref: "int"}, // native types references themselves
			},
		},
		"list": {
			Type: Type{
				Instantiation: TypeInstantiation{Ref: "list"}, // native types references themselves  (params?)
			},
			Params: []string{"t"},
		},
		"custom": { // custom<t> = { x: list<t> }
			Params: []string{"t"},
			Type: Type{
				StructLit: map[string]Type{
					"x": {
						Instantiation: TypeInstantiation{
							Ref: "list",
							Args: []Type{
								{
									Instantiation: TypeInstantiation{Ref: "t"}, // ref to param
								},
							},
						},
					},
				},
			},
		},
	}

	expr := Type{ // custom<int> -> {x: int}
		Instantiation: TypeInstantiation{
			Ref: "custom",
			Args: []Type{
				{Instantiation: TypeInstantiation{Ref: "int"}},
			},
		},
	}

	got, err := resolve(expr, scope)
	if err != nil {
		panic(err)
	}

	want := Type{
		StructLit: map[string]Type{
			"x": {
				Instantiation: TypeInstantiation{
					Ref: "list",
					Args: []Type{
						{
							Instantiation: TypeInstantiation{Ref: "int"},
						},
					},
				},
			},
		},
	}

	g, w := fmt.Sprint(got), fmt.Sprint(want)

	if fmt.Sprint(g) != fmt.Sprint(w) {
		panic("not equal")
	}

	fmt.Println(got)
}

// func test2() {
// 	scope := map[string]TypeDef{ // int = int, list<t> = list
// 		"int": {
// 			Type: Type{
// 				Instantiation: TypeInstantiation{Ref: "int"}, // native types references themselves
// 			},
// 		},
// 		"list": {
// 			Type: Type{
// 				Instantiation: TypeInstantiation{Ref: "list"}, // native types references themselves  (params?)
// 			},
// 			Params: []string{"t"},
// 		},
// 		"custom": { // custom<t> = { x: list<t> }
// 			Params: []string{"t"},
// 			Type: Type{
// 				StructLit: map[string]Type{
// 					"x": {
// 						Instantiation: TypeInstantiation{
// 							Ref: "list",
// 							Args: []Type{
// 								{
// 									Instantiation: TypeInstantiation{Ref: "t"}, // ref to param
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	expr := Type{ // custom<int> -> {x: int}
// 		Instantiation: TypeInstantiation{
// 			Ref: "custom",
// 			Args: []Type{
// 				{Instantiation: TypeInstantiation{Ref: "int"}},
// 			},
// 		},
// 	}

// 	got, err := resolve(expr, scope)
// 	if err != nil {
// 		panic(err)
// 	}

// 	want := Type{
// 		StructLit: map[string]Type{
// 			"x": {
// 				Instantiation: TypeInstantiation{
// 					Ref: "list",
// 					Args: []Type{
// 						{
// 							Instantiation: TypeInstantiation{Ref: "int"},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	g, w := fmt.Sprint(got), fmt.Sprint(want)

// 	if fmt.Sprint(g) != fmt.Sprint(w) {
// 		panic("not equal")
// 	}

// 	fmt.Println(got)
// }

// func test1() {
// 	scope := map[string]TypeDef{ // int = int, list<t> = list
// 		"int": {
// 			Type: Type{
// 				Instantiation: TypeInstantiation{Ref: "int"}, // native types references themselves
// 			},
// 		},
// 		"list": {
// 			Type: Type{
// 				Instantiation: TypeInstantiation{Ref: "list"}, // native types references themselves  (params?)
// 			},
// 			Params: []string{"t"},
// 		},
// 		"custom": { // custom<t> = list<list<t>>
// 			Params: []string{"t"},
// 			Type: Type{
// 				Instantiation: TypeInstantiation{
// 					Ref: "list",
// 					Args: []Type{
// 						{
// 							Instantiation: TypeInstantiation{
// 								Ref: "list",
// 								Args: []Type{
// 									{
// 										Instantiation: TypeInstantiation{Ref: "t"}, // from params
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	expr := Type{ // custom<int> -> list<list<int>>
// 		Instantiation: TypeInstantiation{
// 			Ref: "custom",
// 			Args: []Type{
// 				{Instantiation: TypeInstantiation{Ref: "int"}},
// 			},
// 		},
// 	}

// 	got, err := resolve(expr, scope)
// 	if err != nil {
// 		panic(err)
// 	}

// 	want := Type{
// 		Instantiation: TypeInstantiation{
// 			Ref: "list",
// 			Args: []Type{
// 				{
// 					Instantiation: TypeInstantiation{
// 						Ref: "list",
// 						Args: []Type{
// 							{
// 								Instantiation: TypeInstantiation{Ref: "int"},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	if fmt.Sprint(got.Instantiation) != fmt.Sprint(want.Instantiation) {
// 		panic("not equal")
// 	}

// 	fmt.Println("Got: ", got, "Want: ", want)
// }

func (expr Type) String() string {
	var s string

	if expr.StructLit != nil {
		s += "{"
		for fieldName, fieldExpr := range expr.StructLit {
			s += " " + fieldName + ": " + fieldExpr.String() + " "
		}
		s += "}"
		return s
	}

	if len(expr.Instantiation.Args) == 0 {
		return expr.Instantiation.Ref
	}

	s = expr.Instantiation.Ref + "<"
	for i, arg := range expr.Instantiation.Args {
		s += arg.String()
		if i < len(expr.Instantiation.Args)-1 {
			s += ", "
		}
	}
	s += ">"

	return s
}
