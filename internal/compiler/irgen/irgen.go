package irgen

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/pkg"
)

type Generator struct{}

type (
	nodeContext struct {
		path       []string
		node       src.Node
		portsUsage portsUsage
	}

	portsUsage struct {
		in  map[relPortAddr]struct{}
		out map[relPortAddr]struct{}
	}

	relPortAddr struct {
		Port string
		Idx  *uint8
	}
)

func (g Generator) Generate(
	build src.Build,
	mainPkgName string,
) (*ir.Program, error) {
	return g.GenerateForComponent(
		build,
		mainPkgName,
		"Main",
	)
}

// GenerateForComponent builds IR for a given exported component as the root.
// it maps all component inports to the program start payload and all outports
// to the program stop payload, enabling call/return semantics.
func (g Generator) GenerateForComponent(
	build src.Build,
	pkgName string,
	componentName string,
) (*ir.Program, error) {
	// prepare root node context for the target component with ports usage
	scope := src.NewScope(build, core.Location{
		ModRef:   build.EntryModRef,
		Package:  pkgName,
		Filename: "", // file is initially unknown
	})
	entity, loc, err := scope.Entity(core.EntityRef{
		Pkg:  "",
		Name: componentName,
	})
	if err != nil {
		return nil, err
	}

	// entry point is expected to be not overloaded
	version := entity.Component[0]

	// in and out ports of the root node are expected to be used by runtime
	inUsage := make(map[relPortAddr]struct{}, len(version.IO.In))
	for inName := range version.IO.In {
		inUsage[relPortAddr{Port: inName}] = struct{}{}
	}
	outUsage := make(map[relPortAddr]struct{}, len(version.IO.Out))
	for outName := range version.IO.Out {
		outUsage[relPortAddr{Port: outName}] = struct{}{}
	}

	rootNodeCtx := nodeContext{
		path: []string{},
		node: src.Node{
			EntityRef: core.EntityRef{
				Pkg:  "", // package in reference is empty because entity is local
				Name: componentName,
			},
			Meta: core.Meta{Location: loc},
		},
		portsUsage: portsUsage{
			in:  inUsage,
			out: outUsage,
		},
	}

	result := &ir.Program{
		Connections: map[ir.PortAddr]ir.PortAddr{},
		Funcs:       []ir.FuncCall{},
		Comment: buildProgramComment(
			build.EntryModRef.Path,
			build.EntryModRef.Version,
			pkgName,
		),
	}

	g.processNode(rootNodeCtx, scope, result)

	return &ir.Program{
		Connections: result.Connections,
		Funcs:       result.Funcs,
		Comment:     result.Comment,
	}, nil
}

func buildProgramComment(modulePath, moduleVersion, mainPackage string) string {
	return fmt.Sprintf(
		"// module=%s@%s main=%s compiler=%s",
		modulePath,
		moduleVersion,
		mainPackage,
		pkg.Version,
	)
}

func (g Generator) processNode(
	nodeCtx nodeContext,
	scope src.Scope,
	result *ir.Program,
) {
	entity, location, err := scope.
		Relocate(nodeCtx.node.Meta.Location).
		Entity(nodeCtx.node.EntityRef)
	if err != nil {
		panic(err)
	}

	components := entity.Component
	inportAddrs := g.insertAndReturnInports(nodeCtx)   // for inports we only use parent context because all inports are used
	outportAddrs := g.insertAndReturnOutports(nodeCtx) //  for outports we use both parent context and component's interface

	runtimeFuncRef, version, err := g.getFuncRef(components, nodeCtx.node)
	if err != nil {
		panic(err)
	}

	if runtimeFuncRef != "" {
		cfgMsg, err := getConfigMsg(nodeCtx.node, scope)
		if err != nil {
			panic(err)
		}
		result.Funcs = append(result.Funcs, ir.FuncCall{
			Ref: runtimeFuncRef,
			IO: ir.FuncIO{
				In:  inportAddrs,
				Out: outportAddrs,
			},
			Msg: cfgMsg,
		})
		return
	}

	// We use network as a source of true about how subnodes ports instead subnodes interface definitions.
	// We cannot rely on them because there's no information about how many array slots are used (in case of array ports).
	// On the other hand, we believe network has everything we need because program' correctness is verified by analyzer.
	subnodesPortsUsage, err := g.processNetwork(
		version.Net,
		&scope,
		nodeCtx,
		result,
	)
	if err != nil {
		panic(err)
	}

	for subnodeName, subnode := range version.Nodes {
		nodePortsUsage, ok := subnodesPortsUsage[subnodeName]
		if !ok {
			panic(fmt.Errorf("node usage not found: %v", subnodeName))
		}

		// TODO e2e test
		// sometimes DI nodes are drilled down
		// example: `handler Pass<T>{handler IHandler<T>}`
		// our component is used like this `Parent{handler FilterOdd<T>}`
		// Parent.handler is not interface, but its component has interface
		// It needs our DI nodes, so we merge our DI with node's DI
		if len(nodeCtx.node.DIArgs) > 0 {
			if subnode.DIArgs == nil {
				subnode.DIArgs = make(map[string]src.Node)
			}
			for k, ourDIarg := range nodeCtx.node.DIArgs {
				// FIXME HOC with drilled di arg can't resolve ref to it
				// because it resolves it with its own location
				// rather than with the location, where drilled DI arg was first passed

				// FIXME handle case when DI args drilled anonymously
				// e.g. when we pass Filter{Predicate} so Split{Predicate} works
				// and predicate is k="" figure out by desugarer
				// to do so, we need to take first Split's DI param name and use it instead of k
				// Analyzer probably must check where it's possible to use anonymous DI args and where not
				// without explicit names we would have to traverse all component tree down to the leaf
				// that uses actual dependency, to get its DI name
				// so anonymous DI args are only possible without DI drilling.

				// if sub-node doesn't have DI arg, we just add it
				existing, exists := subnode.DIArgs[k]
				if !exists {
					subnode.DIArgs[k] = ourDIarg
					continue
				}

				// if sub-node has DI arg, we check if it's interface
				kind, err := scope.GetEntityKind(existing.EntityRef)
				if err != nil {
					panic(err)
				}

				// if it's interface, we replace it with our DI arg
				// that's how we can drill DI arguments down to composite components
				if kind == src.InterfaceEntity {
					subnode.DIArgs[k] = ourDIarg
				}
			}
		}

		subNodeCtx := nodeContext{
			path:       append(nodeCtx.path, subnodeName),
			portsUsage: nodePortsUsage,
			node:       subnode,
		}

		var scopeToUse src.Scope
		if injectedSubnode, isDISubnode := nodeCtx.node.DIArgs[subnodeName]; isDISubnode {
			subNodeCtx.node = injectedSubnode
			scopeToUse = scope
		} else {
			scopeToUse = scope.Relocate(location)
		}

		g.processNode(subNodeCtx, scopeToUse, result)
	}
}

func (Generator) insertAndReturnInports(nodeCtx nodeContext) []ir.PortAddr {
	inports := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.in))

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use component's IO
	// actually we can't use IO because we need to know how many slots are used
	for relAddr := range nodeCtx.portsUsage.in {
		absAddr := ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "in"),
			Port: relAddr.Port,
		}
		if relAddr.Idx != nil {
			absAddr.IsArray = true
			absAddr.Idx = *relAddr.Idx
		}
		inports = append(inports, absAddr)
	}

	sortPortAddrs(inports)

	return inports
}

func (Generator) insertAndReturnOutports(nodeCtx nodeContext) []ir.PortAddr {
	outports := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.out))

	// In a valid (desugared) program all outports are used so it's safe to depend on nodeCtx and not use component's IO.
	// Actually we can't use IO because we need to know how many slots are used.
	for addr := range nodeCtx.portsUsage.out {
		irAddr := ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "out"),
			Port: addr.Port,
		}
		if addr.Idx != nil {
			irAddr.IsArray = true
			irAddr.Idx = *addr.Idx
		}
		outports = append(outports, irAddr)
	}

	sortPortAddrs(outports)

	return outports
}

func New() Generator {
	return Generator{}
}
