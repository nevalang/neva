package analyze

import (
	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

// analyzeRootComponent checks root-component-specific requirements:
// Has only one type parameter without constraint;
// Doesn't have any outports;
// Has one not-array inport named `sig` that refers to that single type-parameter;
// Has at least one node (non-root components can have no nodes to implement just routing);
// All nodes has no static inports or references to anything but components (reason why pkg and pkgs needed);
func (a Analyzer) analyzeRootComponent(rootComp src.Component, pkg src.Pkg, pkgs map[string]src.Pkg) error { //nolint:funlen,unparam,lll
	if len(rootComp.TypeParams) != 1 {
		panic("root component must have one type parameter")
	}

	typeParam := rootComp.TypeParams[0]

	if !typeParam.Constr.Empty() {
		panic("type parameter of root component must not have constraint")
	}

	if err := a.analyzeRootComponentIO(rootComp.IO, typeParam); err != nil {
		return err
	}

	if err := a.analyzeRootComponentNodes(rootComp.Nodes, pkg, pkgs); err != nil {
		return err
	}

	return nil
}

func (Analyzer) analyzeRootComponentIO(io src.IO, typeParam ts.Param) error {
	if len(io.Out) != 0 {
		panic("root component can't have outports")
	}

	if len(io.In) != 1 {
		panic("root component must have 1 inport")
	}

	sigInport, ok := io.In["sig"]
	if !ok {
		panic("root component must have 'sig' inport")
	}

	if sigInport.IsArr {
		panic("sig inport of root component can't be array port")
	}

	if !sigInport.Type.Lit.Empty() {
		panic("sig inport of root component can't have literal as type")
	}

	if sigInport.Type.Inst.Ref != typeParam.Name {
		panic("sig inport of root component must refer to type parameter")
	}

	return nil
}

func (Analyzer) analyzeRootComponentNodes(nodes map[string]src.Node, pkg src.Pkg, pkgs map[string]src.Pkg) error {
	if len(nodes) == 0 {
		panic("root component must have nodes")
	}

	for _, node := range nodes {
		if len(node.StaticInports) != 0 {
			panic("root component can't have static inports")
		}

		var pkgWithEntity src.Pkg
		if node.Instance.Ref.Pkg != "" {
			p, ok := pkgs[node.Instance.Ref.Pkg]
			if !ok {
				panic("pkg not found")
			}
			pkgWithEntity = p
		} else {
			pkgWithEntity = pkg
		}

		entity, ok := pkgWithEntity.Entities[node.Instance.Ref.Name]
		if !ok {
			panic("entity not found") 
		}

		if entity.Kind != src.ComponentEntity {
			panic("root component nodes can only refer to other components")
		}
	}

	return nil
}
