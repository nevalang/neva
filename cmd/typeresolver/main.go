package main

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

// Def is a type definition
type Def struct {
	Params []string // type can have parameters
	Expr   Expr     // type must have body in the form of expression
}

// Expr is something that can be resolved. It's either type Instantiation or a type literal.
type Expr struct {
	Instantiation InstantiationExpr // (indirect recursion)
	StructLit     map[string]Expr   // Only struct literals are supported in this version
}

// InstantiationExpr is an expression that always contains reference to some type
// and, if that type has parameters, list of arguments.
type InstantiationExpr struct {
	Ref  string
	Args []Expr // (indirect recursion)
}

// Resolve turns expression into resolved expression if that's possible, otherwise it returns an error.
// Resolved expression is one where all type references points to native types and there's nothing to substitute.
// It's is very much like a β-reduction in λ-calculus.
func Resolve(expr Expr, scope map[string]Def) (Expr, error) {
	if expr.StructLit != nil {
		resolvedStruct := make(map[string]Expr, len(expr.StructLit))
		for field, expr := range expr.StructLit {
			resolvedFieldExpr, err := Resolve(expr, scope)
			if err != nil {
				return Expr{}, errors.New("")
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			StructLit: resolvedStruct,
		}, nil
	}

	refType, ok := scope[expr.Instantiation.Ref] // check that reference type exists
	if !ok {
		return Expr{}, errors.New("")
	}

	if len(refType.Params) > len(expr.Instantiation.Args) { // check that generic args for every param is present
		return Expr{}, errors.New("")
	}

	newScope := make(map[string]Def, len(scope)+len(refType.Params)) // new scope contains resolved params (shadows)
	maps.Copy(newScope, scope)
	resolvedArgs := make([]Expr, 0, len(refType.Params))

	for i, param := range refType.Params {
		resolvedArg, err := Resolve(expr.Instantiation.Args[i], scope)
		if err != nil {
			return Expr{}, errors.New("")
		}
		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param] = Def{
			Expr: resolvedArg,
		}
	}

	if refType.Expr.StructLit == nil { // ref type's body is ot a struct literal
		baseType, ok := scope[refType.Expr.Instantiation.Ref]
		if !ok {
			return Expr{}, errors.New("")
		}
		if expr.Instantiation.Ref == baseType.Expr.Instantiation.Ref {
			return Expr{
				Instantiation: InstantiationExpr{
					Ref:  refType.Expr.Instantiation.Ref,
					Args: resolvedArgs,
				},
				StructLit: nil, // todo
			}, nil
		}
	}

	return Resolve(refType.Expr, newScope) // if it's not a native type and not a literal, then repeat
}

// String generates typescript-like type expressions like foo<bar<baz>>.
func (expr Expr) String() string {
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

func main() {
	scope := map[string]Def{
		"int": { // int = int
			Expr: Expr{
				Instantiation: InstantiationExpr{Ref: "int"}, // native types references themselves
			},
		},
		"list": { // list<t> = list
			Params: []string{"t"},
			Expr: Expr{
				Instantiation: InstantiationExpr{Ref: "list"}, // native types references themselves
			},
		},
		"custom": { // custom<t> = struct{ x list<t> }
			Params: []string{"t"},
			Expr: Expr{
				StructLit: map[string]Expr{
					"x": {
						Instantiation: InstantiationExpr{
							Ref: "list",
							Args: []Expr{
								{
									Instantiation: InstantiationExpr{Ref: "t"}, // ref to param
								},
							},
						},
					},
				},
			},
		},
	}

	expr := Expr{ // custom<int>
		Instantiation: InstantiationExpr{
			Ref: "custom",
			Args: []Expr{
				{Instantiation: InstantiationExpr{Ref: "int"}},
			},
		},
	}

	got, err := Resolve(expr, scope) // custom<int> -> struct{ x list<int> }
	if err != nil {
		panic(err)
	}

	want := Expr{
		StructLit: map[string]Expr{
			"x": {
				Instantiation: InstantiationExpr{
					Ref: "list",
					Args: []Expr{
						{
							Instantiation: InstantiationExpr{Ref: "int"},
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
