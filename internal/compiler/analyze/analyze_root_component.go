package analyze

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrRootComponentNodes             = errors.New("root component nodes")
	ErrRootComponentWithoutNodes      = errors.New("root component must have nodes")
	ErrRootComponentWithStaticInports = errors.New("root component can't have static inports")
	ErrNodeRefPkgNotFound             = errors.New("node refers to not found pkg")
	ErrNodeRefEntityNotFound          = errors.New("node refers to not found entity")
	ErrRootCompWithNotCompNodes       = errors.New("root component nodes can only refer to other components")
)

// analyzeRootComponent checks root-component-specific requirements:
// Has only one type parameter without constraint;
// Doesn't have any outports;
// Has one not-array inport named `sig` that refers to that single type-parameter;
// Has at least one node (non-root components can have no nodes to implement just routing);
// All nodes has no static inports or references to anything but components (reason why pkg and pkgs needed);
func (a Analyzer) analyzeRootComponent(rootComp src.Component, pkg src.Pkg, pkgs map[string]src.Pkg) error { //nolint:funlen,unparam,lll
	if len(rootComp.TypeParams) != 1 {
		return errors.New("root component must have one type parameter")
	}

	typeParam := rootComp.TypeParams[0]

	if !typeParam.Constr.Empty() {
		return errors.New("type parameter of root component must not have constraint")
	}

	if err := a.analyzeRootComponentIO(rootComp.IO, typeParam); err != nil {
		return err
	}

	if err := a.analyzeRootComponentNodes(rootComp.Nodes, pkg, pkgs); err != nil {
		return fmt.Errorf("%w: %v", ErrRootComponentNodes, err)
	}

	return nil
}

func (Analyzer) analyzeRootComponentIO(io src.IO, typeParam ts.Param) error {
	if len(io.Out) != 0 {
		return errors.New("root component can't have outports")
	}

	if len(io.In) != 1 {
		return errors.New("root component must have 1 inport")
	}

	sigInport, ok := io.In["sig"]
	if !ok {
		return errors.New("root component must have 'sig' inport")
	}

	if sigInport.IsArr {
		return errors.New("sig inport of root component can't be array port")
	}

	if !sigInport.Type.Lit.Empty() {
		return errors.New("sig inport of root component can't have literal as type")
	}

	if sigInport.Type.Inst.Ref != typeParam.Name {
		return errors.New("sig inport of root component must refer to type parameter")
	}

	return nil
}

func (Analyzer) analyzeRootComponentNodes(nodes map[string]src.Node, pkg src.Pkg, pkgs map[string]src.Pkg) error {
	if len(nodes) == 0 {
		return ErrRootComponentWithoutNodes
	}

	for _, node := range nodes {
		if len(node.StaticInports) != 0 {
			return fmt.Errorf("%w: %v", ErrRootComponentWithStaticInports, node.StaticInports)
		}

		var pkgWithEntity src.Pkg
		if node.Instance.Ref.Pkg != "" {
			p, ok := pkgs[node.Instance.Ref.Pkg]
			if !ok {
				return fmt.Errorf("%w: %v", ErrNodeRefPkgNotFound, node.Instance.Ref.Pkg)
			}
			pkgWithEntity = p
		} else {
			pkgWithEntity = pkg
		}

		entity, ok := pkgWithEntity.Entities[node.Instance.Ref.Name]
		if !ok {
			return fmt.Errorf("%w: %v", ErrNodeRefEntityNotFound, node.Instance.Ref.Name)
		}

		if entity.Kind != src.ComponentEntity {
			return fmt.Errorf("%w: %v: %v", ErrRootCompWithNotCompNodes, node.Instance.Ref.Name, entity.Kind)
		}
	}

	return nil
}
