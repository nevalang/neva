package analyzer

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrRootComponentNodes              = errors.New("root component nodes")
	ErrRootComponentWithoutNodes       = errors.New("root component must have nodes")
	ErrRootComponentWithStaticInports  = errors.New("root component can't have static inports")
	ErrNodeRefPkgNotFound              = errors.New("node refers to not found pkg")
	ErrNodeRefEntityNotFound           = errors.New("node refers to not found entity")
	ErrRootCompWithNotCompNodes        = errors.New("root component nodes can only refer to other components")
	ErrRootComponentWithOutports       = errors.New("root component can't have outports")
	ErrRootComponentInportCount        = errors.New("root component must have 1 inport")
	ErrMainComponentWithoutStartInport = errors.New("root component must have 'sig' inport")
	ErrStartInportIsArray              = errors.New("sig inport of root component can't be array port")
	ErrSigLit                          = errors.New("sig inport of root component can't have literal as type")
	ErrSigType                         = errors.New("sig inport of root component must refer to type parameter")
)

// analyzeMainComponent checks root-component-specific requirements:
// Has only one type parameter without constraint;
// Doesn't have any outports;
// Has one not-array inport named `sig` that refers to that single type-parameter;
// Has at least one node (non-root components can have no nodes to implement just routing);
// All nodes has no static inports or references to anything but components (reason why pkg and pkgs needed);
func (a Analyzer) analyzeMainComponent(cmp compiler.Component, pkg compiler.Pkg, pkgs map[string]compiler.Pkg) error { //nolint:funlen,unparam,lll
	if len(cmp.TypeParams) != 0 {
		return errors.New("root component must not have type parameters")
	}

	if err := a.analyzeExecCmpIO(cmp.IO); err != nil {
		return err
	}

	if err := a.analyzeExecCmpNodes(cmp.Nodes, pkg, pkgs); err != nil {
		return fmt.Errorf("%w: %v", ErrRootComponentNodes, err)
	}

	return nil
}

func (a Analyzer) analyzeExecCmpIO(io compiler.IO) error {
	if len(io.Out) != 1 {
		return fmt.Errorf("%w: %v", ErrRootComponentWithOutports, io.Out)
	}

	if len(io.In) != 1 {
		return fmt.Errorf("%w: %v", ErrRootComponentInportCount, io.In)
	}

	startInport, ok := io.In["start"]
	if !ok {
		return ErrMainComponentWithoutStartInport
	}

	if startInport.IsArr {
		return ErrStartInportIsArray
	}

	// TODO replace subtype checking with straightforward check that type is empty record
	emptyRecType := h.Rec(map[string]ts.Expr{})
	if err := a.checker.Check(startInport.Type, emptyRecType, ts.TerminatorParams{}); err != nil {
		panic(err)
	}

	return nil
}

func (Analyzer) analyzeExecCmpNodes(nodes map[string]compiler.Node, pkg compiler.Pkg, pkgs map[string]compiler.Pkg) error {
	if len(nodes) == 0 {
		return ErrRootComponentWithoutNodes
	}

	for _, node := range nodes {
		if len(node.StaticInports) != 0 {
			return fmt.Errorf("%w: %v", ErrRootComponentWithStaticInports, node.StaticInports)
		}

		var pkgWithEntity compiler.Pkg
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

		if entity.Kind != compiler.ComponentEntity {
			return fmt.Errorf("%w: %v: %v", ErrRootCompWithNotCompNodes, node.Instance.Ref.Name, entity.Kind)
		}
	}

	return nil
}
