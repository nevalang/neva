package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	"github.com/nevalang/neva/internal/compiler/typesystem"
)

type foundInterface struct {
	iface    src.Interface
	location core.Location
}

func (a Analyzer) analyzeNodes(
	parentComponentName string,
	iface src.Interface, // resolved interface of the component that contains the nodes
	nodes map[string]src.Node, // nodes to analyze
	net []src.Connection, // network of the component that contains the nodes
	scope src.Scope, // scope of the component
) (
	map[string]src.Node, // resolved nodes
	map[string]foundInterface, // resolved nodes interfaces with locations
	bool, // one of the nodes has error guard
	*compiler.Error, // err
) {
	analyzedNodes := make(map[string]src.Node, len(nodes))
	nodesInterfaces := make(map[string]foundInterface, len(nodes))
	hasErrGuard := false

	for nodeName, node := range nodes {
		if node.ErrGuard {
			hasErrGuard = true
		}

		analyzedNode, nodeInterface, err := a.analyzeNode(
			nodeName,
			node,
			parentComponentName,
			scope,
			iface,
			nodes,
			net,
		)
		if err != nil {
			return nil, nil, false, compiler.Error{
				Meta: &node.Meta,
			}.Wrap(err)
		}

		nodesInterfaces[nodeName] = nodeInterface
		analyzedNodes[nodeName] = analyzedNode
	}

	return analyzedNodes, nodesInterfaces, hasErrGuard, nil
}

//nolint:gocyclo // Analyzer node handling is a high-branch routine.
func (a Analyzer) analyzeNode(
	name string, // name of the node
	node src.Node, // node to analyze
	parentComponentName string,
	scope src.Scope, // scope of the component that contains the node
	iface src.Interface, // interface of the component that contains the node
	nodes map[string]src.Node, // nodes of the component that contains the node
	net []src.Connection, // network of the component that contains the node
) (src.Node, foundInterface, *compiler.Error) {
	parentTypeParams := iface.TypeParams

	if node.EntityRef.Pkg == "" && node.EntityRef.Name == parentComponentName {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: fmt.Sprintf(
				"Recursive reference to component %q is not allowed. "+
					"If you meant the builtin component, explicitly import the builtin package and use builtin.%s.",
				parentComponentName,
				parentComponentName,
			),
			Meta: &node.Meta,
		}
	}

	nodeEntity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	if nodeEntity.Kind != src.ComponentEntity &&
		nodeEntity.Kind != src.InterfaceEntity {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: fmt.Sprintf("Node can only refer to flows or interfaces: %v", nodeEntity.Kind),
			Meta:    nodeEntity.Meta(),
		}
	}

	bindArg, hasBind := node.Directives[compiler.BindDirective]
	if hasBind && bindArg == "" {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Node with #bind directive must provide exactly one argument",
			Meta:    nodeEntity.Meta(),
		}
	}

	if hasBind && nodeEntity.Kind == src.InterfaceEntity {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Interface node cannot use #bind directive",
			Meta:    nodeEntity.Meta(),
		}
	}

	if nodeEntity.Kind == src.InterfaceEntity && node.DIArgs != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: "Only component node can have dependency injection",
			Meta:    nodeEntity.Meta(),
		}
	}

	// We need to get resolved frame from parent type parameters
	// in order to be able to resolve node's args
	// since they can refer to type parameter of the parent (interface)
	_, resolvedParentParamsFrame, err := a.resolver.ResolveParams(
		parentTypeParams.Params,
		scope,
	)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	// Now when we have frame made of parent type parameters constraints
	// we can resolve cases like `subnode SubFlow<T>`
	// where `T` refers to type parameter of the component/interface we're in.
	resolvedNodeArgs, err := a.resolver.ResolveExprsWithFrame(
		node.TypeArgs,
		resolvedParentParamsFrame,
		scope,
	)
	if err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &node.Meta,
		}
	}

	var (
		nodeIface     src.Interface
		overloadIndex *int // only for component nodes
	)
	if nodeEntity.Kind == src.InterfaceEntity {
		nodeIface = nodeEntity.Interface
	} else {
		var err *compiler.Error
		nodeIface, overloadIndex, err = a.getInterfaceAndOverloadingIndexForNode(
			name,
			nodeEntity,
			hasBind,
			node,
			scope,
			resolvedNodeArgs,
			iface,
			nodes,
			net,
		)
		if err != nil {
			return src.Node{}, foundInterface{}, err
		}
	}

	if node.ErrGuard {
		if _, ok := iface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message: "Error-guard operator '?' can only be used in components with ':err' outport to propagate errors",
				Meta:    &node.Meta,
			}
		}
		if _, ok := nodeIface.IO.Out["err"]; !ok {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message: "Error-guard operator '?' requires node to have ':err' outport to propagate errors",
				Meta:    &node.Meta,
			}
		}
	}

	// default any
	// TODO: Remove this! See github issue about implicit any in nodes.
	if len(resolvedNodeArgs) == 0 && len(nodeIface.TypeParams.Params) == 1 {
		resolvedNodeArgs = []typesystem.Expr{
			{
				Inst: &typesystem.InstExpr{
					Ref: core.EntityRef{Name: "any"},
				},
			},
		}
	}

	// Finally check that every argument is compatible
	// with corresponding parameter of the node's interface.
	if err = a.resolver.CheckArgsCompatibility(
		resolvedNodeArgs,
		nodeIface.TypeParams.Params,
		scope,
	); err != nil {
		return src.Node{}, foundInterface{}, &compiler.Error{
			Message: fmt.Sprintf("%s: %s", node.EntityRef.Name, err.Error()),
			Meta:    &node.Meta,
		}
	}

	if a.isUnionNode(node) {
		firstResolvedNodeArg := resolvedNodeArgs[0]
		if firstResolvedNodeArg.Lit == nil || firstResolvedNodeArg.Lit.Union == nil {
			return src.Node{}, foundInterface{}, &compiler.Error{
				Message: "Union<T> expects union type argument",
				Meta:    &node.Meta,
			}
		}
	}

	if node.DIArgs == nil {
		return src.Node{
				Directives:    node.Directives,
				EntityRef:     node.EntityRef,
				TypeArgs:      resolvedNodeArgs,
				Meta:          node.Meta,
				OverloadIndex: overloadIndex,
				ErrGuard:      node.ErrGuard,
			}, foundInterface{
				iface:    nodeIface,
				location: location,
			}, nil
	}

	resolvedFlowDI := make(map[string]src.Node, len(node.DIArgs))
	for depName, depNode := range node.DIArgs {
		// di arguments are not regular nodes in the network, so we generate a unique name
		// that won't be found in the network. This will cause the overloading logic to skip
		// network-based checks.
		uniqueName := "__di_" + name + "_" + depName
		resolvedDep, _, err := a.analyzeNode(
			uniqueName, // use a unique name that won't be found in the network
			depNode,
			parentComponentName,
			scope,
			iface,
			nodes,
			net, // pass the network so that type constraints can be collected
		)
		if err != nil {
			return src.Node{}, foundInterface{}, compiler.Error{
				Meta: &depNode.Meta,
			}.Wrap(err)
		}
		resolvedFlowDI[depName] = resolvedDep
	}

	return src.Node{
			Directives:    node.Directives,
			EntityRef:     node.EntityRef,
			TypeArgs:      resolvedNodeArgs,
			DIArgs:        resolvedFlowDI,
			ErrGuard:      node.ErrGuard,
			OverloadIndex: overloadIndex,
			Meta:          node.Meta,
		}, foundInterface{
			iface:    nodeIface,
			location: location,
		}, nil
}

// getInterfaceAndOverloadingIndexForNode returns interface and overload index for the given node.
// Overloading at the level of sourcecode is implemented here.
func (a Analyzer) getInterfaceAndOverloadingIndexForNode(
	nodeName string,
	entity src.Entity,
	hasBind bool,
	node src.Node,
	scope src.Scope,
	resolvedNodeArgs []typesystem.Expr,
	resolvedParentIface src.Interface, // resolved interface of the component that contains the node
	allParentNodes map[string]src.Node, // nodes of the component that contains the node
	net []src.Connection, // network of the component that contains the node
) (src.Interface, *int, *compiler.Error) {
	var (
		overloadIndex *int
		version       src.Component
	)
	if len(entity.Component) == 1 {
		version = entity.Component[0]
	} else {
		var err *compiler.Error
		version, overloadIndex, err = a.getNodeOverloadVersionAndIndex(
			nodeName,
			resolvedParentIface,
			scope,
			resolvedNodeArgs,
			allParentNodes,
			net,
			entity,
		)
		if err != nil {
			return src.Interface{}, nil, &compiler.Error{
				Message: err.Error(),
				Meta:    &node.Meta,
			}
		}
	}

	_, hasExtern := version.Directives[compiler.ExternDirective]
	if hasBind && !hasExtern {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Node can't use #bind if it isn't instantiated with the component that use #extern",
			Meta:    entity.Meta(),
		}
	}

	versionIface := version.Interface

	_, hasAutoPortsDirective := version.Directives[compiler.AutoportsDirective]
	if !hasAutoPortsDirective {
		return versionIface, overloadIndex, nil
	}

	// if we here then we have #autoports (only for structs)

	if len(versionIface.IO.In) != 0 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Component that uses struct inports directive must have no defined inports",
			Meta:    entity.Meta(),
		}
	}

	if len(versionIface.TypeParams.Params) != 1 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Exactly one type parameter expected",
			Meta:    entity.Meta(),
		}
	}

	resolvedTypeParamConstr, err := a.resolver.ResolveExpr(versionIface.TypeParams.Params[0].Constr, scope)
	if err != nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: err.Error(),
			Meta:    entity.Meta(),
		}
	}

	if resolvedTypeParamConstr.Lit == nil || resolvedTypeParamConstr.Lit.Struct == nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Struct type expected",
			Meta:    entity.Meta(),
		}
	}

	if len(resolvedNodeArgs) != 1 {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Exactly one type argument expected",
			Meta:    entity.Meta(),
		}
	}

	resolvedNodeArg, err := a.resolver.ResolveExpr(resolvedNodeArgs[0], scope)
	if err != nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: err.Error(),
			Meta:    entity.Meta(),
		}
	}

	if resolvedNodeArg.Lit == nil || resolvedNodeArg.Lit.Struct == nil {
		return src.Interface{}, nil, &compiler.Error{
			Message: "Struct argument expected",
			Meta:    entity.Meta(),
		}
	}

	structFields := resolvedNodeArg.Lit.Struct
	inports := make(map[string]src.Port, len(structFields))
	for fieldName, fieldTypeExpr := range structFields {
		inports[fieldName] = src.Port{
			TypeExpr: fieldTypeExpr,
		}
	}

	return src.Interface{
		TypeParams: versionIface.TypeParams,
		IO: src.IO{
			In: inports,
			Out: map[string]src.Port{
				// struct builder has exactly one outport - created structure
				"res": {
					TypeExpr: resolvedNodeArg,
					IsArray:  false,
					Meta:     versionIface.IO.Out["res"].Meta,
				},
			},
		},
		Meta: versionIface.Meta,
	}, overloadIndex, nil
}

// getNodeOverloadVersionAndIndex determines which overload of a component to use for a node.
// This is called when we know the node references an overloaded component with multiple implementations.
// It analyzes how the node is used in connections to determine the appropriate implementation.
func (a Analyzer) getNodeOverloadVersionAndIndex(
	nodeName string,
	resolvedParentIface src.Interface,
	scope src.Scope,
	resolvedNodeArgs []typesystem.Expr,
	allParentNodes map[string]src.Node,
	net []src.Connection,
	entity src.Entity,
) (src.Component, *int, *compiler.Error) {
	nodeRefs := findNodeRefsInNet(nodeName, net)

	// for di arguments (which have no node references in the network),
	// we skip network-based overloading and rely only on interface compatibility.
	// di arguments have unique names that won't be found in the network.
	isDIArg := len(nodeRefs) == 0

	var nodeConstraints nodeUsageConstraints
	if isDIArg {
		// for di arguments, we need to get type constraints from the parent component's dependency declaration
		nodeConstraints = a.collectDITypeConstraintsFromParent(
			nodeName,
			scope,
			allParentNodes,
		)
	} else {
		// For overloaded components we need to lookup network in separate run to find how the node is used.
		nodeConstraints = a.deriveNodeConstraintsFromNetwork(
			nodeName,
			resolvedParentIface,
			scope,
			allParentNodes,
			net,
		)
	}

	remainingIdx := make([]int, 0, len(entity.Component))
	remainingComps := make([]src.Component, 0, len(entity.Component))
	for i, component := range entity.Component {
		// skip network-based compatibility check for di arguments
		if !isDIArg {
			if !isCandidateCompatibleWithAllNodeRefs(component, nodeRefs) {
				continue
			}
		}

		// for native components with multiple extern implementations, use type argument-based disambiguation
		// this handles cases like Dec(int) vs Dec(float) or Len(list) vs Len(string)
		if a.isNativeComponentWithMultipleExterns(component, entity) {
			// if we have type arguments, use them for disambiguation
			if len(resolvedNodeArgs) > 0 {
				if !a.doesNativeComponentMatchTypeArgs(component, resolvedNodeArgs, scope) {
					continue
				}
			} else {
				// if no type arguments, use network-based type constraints
				// this handles cases like Len without explicit type args but with network usage
				if !a.doesCandidateSatisfyTypeConstraints(
					component.Interface,
					resolvedNodeArgs,
					nodeConstraints,
					scope,
				) {
					continue
				}
			}
		} else {
			// for regular components, use the existing type constraint logic
			if !a.doesCandidateSatisfyTypeConstraints(
				component.Interface,
				resolvedNodeArgs,
				nodeConstraints,
				scope,
			) {
				continue
			}
		}

		remainingIdx = append(remainingIdx, i)
		remainingComps = append(remainingComps, component)
	}

	if len(remainingComps) == 0 {
		return src.Component{}, nil, &compiler.Error{
			Message: fmt.Sprintf(
				"no compatible overload found for node %s (total components: %d, remaining: %d)",
				nodeName,
				len(entity.Component),
				len(remainingComps),
			),
			Meta:    entity.Meta(),
		}
	}

	if l := len(remainingComps); l > 1 {
		return src.Component{}, nil, &compiler.Error{
			Message: fmt.Sprintf(
				"ambiguous overload for node %s: multiple candidates satisfy usage: %d",
				nodeName, l,
			),
			Meta: entity.Meta(),
		}
	}

	return remainingComps[0], &remainingIdx[0], nil
}

// collectDITypeConstraintsFromParent collects type constraints for a DI argument
// from the parent component's dependency declaration.
//nolint:gocyclo // DI constraint derivation has multiple cases to check.
func (a Analyzer) collectDITypeConstraintsFromParent(
	nodeName string, // the unique name of the DI argument (e.g., "__di_reduce_reducer")
	scope src.Scope,
	allParentNodes map[string]src.Node,
) nodeUsageConstraints {
	// extract the dependency name from the unique node name
	// format: "__di_<parentNodeName>_<depName>"
	// we need to find the last underscore to get the depName
	lastUnderscore := -1
	for i := len(nodeName) - 1; i >= 0; i-- {
		if nodeName[i] == '_' {
			lastUnderscore = i
			break
		}
	}
	if lastUnderscore == -1 {
		return emptyConstraints()
	}

	depName := nodeName[lastUnderscore+1:]

	// find the parent node that contains this dependency
	// we need to look through all parent nodes to find the one that has this dependency
	var parentNodeName string
	var parentNode src.Node
	for nodeName, node := range allParentNodes {
		if node.DIArgs != nil {
			if _, hasDep := node.DIArgs[depName]; hasDep {
				parentNodeName = nodeName
				parentNode = node
				break
			}
		}
	}

	if parentNodeName == "" {
		return emptyConstraints()
	}

	// get the parent component to find the dependency declaration
	parentEntity, _, err := scope.Entity(parentNode.EntityRef)
	if err != nil {
		return emptyConstraints()
	}

	// find the component version that matches the parent node
	var parentComponent src.Component
	switch {
	case len(parentEntity.Component) == 1:
		parentComponent = parentEntity.Component[0]
	case parentNode.OverloadIndex != nil:
		parentComponent = parentEntity.Component[*parentNode.OverloadIndex]
	default:
		return emptyConstraints()
	}

	// we need to look up dependencies in the parent component's scope
	// because the nodes in parentComponent refer to entities in that scope
	parentScope := scope.Relocate(parentComponent.Interface.Meta.Location)

	// find the dependency declaration in the parent component's nodes
	var depNode src.Node
	var hasDep bool

	if depName == "" {
		// for anonymous dependencies, find the first interface node
		for _, node := range parentComponent.Nodes {
			entity, _, err := parentScope.Entity(node.EntityRef)
			if err == nil && entity.Kind == src.InterfaceEntity {
				depNode = node
				hasDep = true
				break
			}
		}
	} else {
		depNode, hasDep = parentComponent.Nodes[depName]
	}

	if !hasDep {
		return emptyConstraints()
	}

	// get the dependency interface
	depEntity, _, err := parentScope.Entity(depNode.EntityRef)
	if err != nil {
		return emptyConstraints()
	}

	if depEntity.Kind != src.InterfaceEntity {
		return emptyConstraints()
	}

	// resolve the dependency interface with the parent's type arguments
	// we need to create a frame from the parent's type arguments
	// the parent component's original interface has the type parameters
	parentTypeFrame := make(map[string]typesystem.Def)
	for i, param := range parentComponent.TypeParams.Params {
		if i < len(parentNode.TypeArgs) {
			parentTypeFrame[param.Name] = typesystem.Def{
				BodyExpr: &parentNode.TypeArgs[i],
				Meta:     param.Constr.Meta,
			}
		}
	}

	// convert the dependency interface to type constraints by resolving each port type
	constraints := nodeUsageConstraints{
		incoming: make(map[string][]typesystem.Expr),
		outgoing: make(map[string][]typesystem.Expr),
	}

	// add incoming constraints (input ports) - resolve each port type with the parent's type frame
	for portName, port := range depEntity.Interface.IO.In {
		resolvedType, err := a.resolver.ResolveExprWithFrame(port.TypeExpr, parentTypeFrame, scope)
		if err != nil {
			return emptyConstraints()
		}
		constraints.incoming[portName] = []typesystem.Expr{resolvedType}
	}

	// add outgoing constraints (output ports) - resolve each port type with the parent's type frame
	for portName, port := range depEntity.Interface.IO.Out {
		resolvedType, err := a.resolver.ResolveExprWithFrame(port.TypeExpr, parentTypeFrame, scope)
		if err != nil {
			return emptyConstraints()
		}
		constraints.outgoing[portName] = []typesystem.Expr{resolvedType}
	}

	// extract the dependency name from the unique node name
	// this handles cases where the interface has unnamed ports but the component has named ports
	if len(constraints.incoming) == 1 {
		for portName, types := range constraints.incoming {
			if portName == "" {
				// find the actual port name from the first candidate component
				// we need to look at the entity to find the actual port name
				// for now, let's use a common port name like "data"
				delete(constraints.incoming, "")
				constraints.incoming["data"] = types
				break
			}
		}
	}

	if len(constraints.outgoing) == 1 {
		for portName, types := range constraints.outgoing {
			if portName == "" {
				// find the actual port name from the first candidate component
				// we need to look at the entity to find the actual port name
				// for now, let's use a common port name like "res"
				delete(constraints.outgoing, "")
				constraints.outgoing["res"] = types
				break
			}
		}
	}

	return constraints
}

// isCandidateCompatibleWithAllNodeRefs checks if the given component is compatible
// with all the node references by checking that all the node references
// are compatible with the component's interface.
// Compatibility is checked by comparing the port types and array usage.
// The type of the port ignored for now, to be checked later.
func isCandidateCompatibleWithAllNodeRefs(
	component src.Component,
	nodeRefs []nodeRefInNet,
) bool {
	for _, nodeRef := range nodeRefs {
		var ports map[string]src.Port
		if !nodeRef.isOutgoing {
			ports = component.IO.In
		} else {
			ports = component.IO.Out
		}

		// handle empty port names by resolving to the actual port name
		portName := nodeRef.port
		if portName == "" {
			// if there's only one port, use that
			if len(ports) == 1 {
				for name := range ports {
					portName = name
					break
				}
			} else {
				// multiple ports but no port name specified - this is an error
				return false
			}
		}

		port, portExists := ports[portName]
		if !portExists {
			return false
		}

		isArrExpected := nodeRef.arrayIdx != nil
		if port.IsArray != isArrExpected {
			return false
		}
	}

	return true
}

// findNodeRefsInNet identifies all places where the specified node is referenced in a connection.
func findNodeRefsInNet(nodeName string, connections []src.Connection) []nodeRefInNet {
	var refs []nodeRefInNet

	for _, conn := range connections {
		if conn.ArrayBypass != nil {
			if conn.ArrayBypass.SenderOutport.Node == nodeName {
				refs = append(refs, nodeRefInNet{
					isOutgoing: true,
					port:       conn.ArrayBypass.SenderOutport.Port,
					arrayIdx:   conn.ArrayBypass.SenderOutport.Idx,
				})
			}
			if conn.ArrayBypass.ReceiverInport.Node == nodeName {
				refs = append(refs, nodeRefInNet{
					isOutgoing: false,
					port:       conn.ArrayBypass.ReceiverInport.Port,
					arrayIdx:   conn.ArrayBypass.ReceiverInport.Idx,
				})
			}
			continue
		}

		if conn.Normal != nil {
			for _, sender := range conn.Normal.Senders {
				if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
					refs = append(refs, nodeRefInNet{
						isOutgoing: true,
						port:       sender.PortAddr.Port,
						arrayIdx:   sender.PortAddr.Idx,
					})
				}
			}
			refs = append(refs, findNodeUsagesInReceivers(nodeName, conn.Normal.Receivers)...)
		}
	}

	return refs
}

// findNodeUsagesInReceivers recursively checks receivers for node usages
func findNodeUsagesInReceivers(nodeName string, receivers []src.ConnectionReceiver) []nodeRefInNet {
	var nodeRefs []nodeRefInNet

	for _, receiver := range receivers {
		// Check direct port address
		if receiver.PortAddr != nil && receiver.PortAddr.Node == nodeName {
			nodeRefs = append(nodeRefs, nodeRefInNet{
				isOutgoing: false,
				port:       receiver.PortAddr.Port,
				arrayIdx:   receiver.PortAddr.Idx,
			})
		}

		// Check chained connection
		if receiver.ChainedConnection != nil {
			if receiver.ChainedConnection.Normal != nil {
				// Check senders in the chain
				for _, sender := range receiver.ChainedConnection.Normal.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						nodeRefs = append(nodeRefs, nodeRefInNet{
							isOutgoing: true,
							port:       sender.PortAddr.Port,
							arrayIdx:   sender.PortAddr.Idx,
						})
					}
				}

				// Recursively check receivers in the chain
				nodeRefs = append(nodeRefs, findNodeUsagesInReceivers(nodeName, receiver.ChainedConnection.Normal.Receivers)...)
			}
		}

		// Check deferred connection
		if receiver.DeferredConnection != nil {
			// Similar logic to what we do with normal connections
			if receiver.DeferredConnection.Normal != nil {
				for _, sender := range receiver.DeferredConnection.Normal.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						nodeRefs = append(nodeRefs, nodeRefInNet{
							isOutgoing: true,
							port:       sender.PortAddr.Port,
							arrayIdx:   sender.PortAddr.Idx,
						})
					}
				}

				nodeRefs = append(nodeRefs, findNodeUsagesInReceivers(nodeName, receiver.DeferredConnection.Normal.Receivers)...)
			}
		}
	}

	return nodeRefs
}

//nolint:govet // fieldalignment: keep semantic grouping.
type nodeRefInNet struct {
	isOutgoing bool
	port       string
	arrayIdx   *uint8
}

// nodeUsageConstraints captures incoming produced types and outgoing expected types per port.
type nodeUsageConstraints struct {
	incoming map[string][]typesystem.Expr // for inports: types produced by connected senders
	outgoing map[string][]typesystem.Expr // for outports: types expected by connected receivers
}

// emptyConstraints returns an empty nodeUsageConstraints with initialized maps
func emptyConstraints() nodeUsageConstraints {
	return nodeUsageConstraints{
		incoming: make(map[string][]typesystem.Expr),
		outgoing: make(map[string][]typesystem.Expr),
	}
}

// deriveNodeConstraintsFromNetwork inspects the network and extracts type constraints for the given node.
// It only uses available information (parent iface, literals, neighbor node interfaces). it does not depend on nodesIfaces.
// It is needed only to select correct version of the overloaded component.
//nolint:gocyclo // Node constraint derivation handles many network patterns.
func (a Analyzer) deriveNodeConstraintsFromNetwork(
	nodeName string,
	resolvedParentIface src.Interface,
	scope src.Scope,
	nodes map[string]src.Node,
	net []src.Connection,
) nodeUsageConstraints {
	c := nodeUsageConstraints{
		incoming: map[string][]typesystem.Expr{},
		outgoing: map[string][]typesystem.Expr{},
	}

	// helper to append unique types (by String) to a slice
	appendUnique := func(dst *[]typesystem.Expr, t typesystem.Expr) {
		s := t.String()
		for _, e := range *dst {
			if e.String() == s {
				return
			}
		}
		*dst = append(*dst, t)
	}

	// resolve parent param frame once
	_, parentFrame, err := a.resolver.ResolveParams(
		resolvedParentIface.TypeParams.Params, scope,
	)
	if err != nil {
		return c
	}

	// walk all connections
	for _, conn := range net {
		if conn.Normal == nil {
			continue
		}

		// check if our node is a sender in this connection (including chained connections)
		for _, sender := range conn.Normal.Senders {
			if sender.PortAddr == nil || sender.PortAddr.Node != nodeName {
				continue
			}
			port := a.resolvePortName(nodeName, nodes, scope, false, sender.PortAddr.Port)
			// derive expected types from all receivers
			recvPortAddrs := a.flattenReceiversPortAddrs(conn.Normal.Receivers)
			for _, rpa := range recvPortAddrs {
				// parent out
				if rpa.Node == "out" {
					if p, ok := resolvedParentIface.IO.Out[rpa.Port]; ok {
						if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, parentFrame, scope); err == nil {
							list := c.outgoing[port]
							appendUnique(&list, resolved)
							c.outgoing[port] = list
						}
					}
					continue
				}
				// other node inport
				recvNode, ok := nodes[rpa.Node]
				if !ok {
					continue
				}
				// get possible inport types across overloads
				for _, t := range a.getPossibleNodePortTypes(scope, parentFrame, recvNode, true, rpa.Port) {
					list := c.outgoing[port]
					appendUnique(&list, t)
					c.outgoing[port] = list
				}
			}

			// also check chained connections for outgoing constraints
			for _, receiver := range conn.Normal.Receivers {
				if receiver.ChainedConnection != nil && receiver.ChainedConnection.Normal != nil {
					// this is a chained connection, look at the receivers within the chain
					chainedRecvPortAddrs := a.flattenReceiversPortAddrs(receiver.ChainedConnection.Normal.Receivers)
					for _, rpa := range chainedRecvPortAddrs {
						// parent out
						if rpa.Node == "out" {
							if p, ok := resolvedParentIface.IO.Out[rpa.Port]; ok {
								if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, parentFrame, scope); err == nil {
									list := c.outgoing[port]
									appendUnique(&list, resolved)
									c.outgoing[port] = list
								}
							}
							continue
						}
						// other node inport
						recvNode, ok := nodes[rpa.Node]
						if !ok {
							continue
						}
						// get possible inport types across overloads
						for _, t := range a.getPossibleNodePortTypes(scope, parentFrame, recvNode, true, rpa.Port) {
							list := c.outgoing[port]
							appendUnique(&list, t)
							c.outgoing[port] = list
						}
					}
				}
			}
		}

		// also check if our node is a sender within chained connections
		// in a connection like "a -> b -> c", b is a sender in a chained connection
		// and should receive incoming constraints from a
		// this needs to be recursive to handle nested chains like "a -> b -> c -> d"
		var checkChainedConnections func(outerSenders []src.ConnectionSender, receivers []src.ConnectionReceiver)
		checkChainedConnections = func(outerSenders []src.ConnectionSender, receivers []src.ConnectionReceiver) {
			for _, receiver := range receivers {
				if receiver.ChainedConnection != nil && receiver.ChainedConnection.Normal != nil {
					for _, sender := range receiver.ChainedConnection.Normal.Senders {
						if sender.PortAddr == nil || sender.PortAddr.Node != nodeName {
							continue
						}

						// this node is a sender in a chained connection
						// collect incoming constraints from the outer senders
						inPort := a.resolvePortName(nodeName, nodes, scope, true, "")

						// collect types from the outer senders
						for _, outerSender := range outerSenders {
							types := a.getPossibleSenderTypes(scope, parentFrame, resolvedParentIface, nodes, outerSender, net)
							for _, t := range types {
								list := c.incoming[inPort]
								appendUnique(&list, t)
								c.incoming[inPort] = list
							}
						}

						// collect outgoing constraints from the chained connection's receivers
						port := a.resolvePortName(nodeName, nodes, scope, false, sender.PortAddr.Port)
						chainedRecvPortAddrs := a.flattenReceiversPortAddrs(receiver.ChainedConnection.Normal.Receivers)
						for _, rpa := range chainedRecvPortAddrs {
							if rpa.Node == "out" {
								if p, ok := resolvedParentIface.IO.Out[rpa.Port]; ok {
									if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, parentFrame, scope); err == nil {
										list := c.outgoing[port]
										appendUnique(&list, resolved)
										c.outgoing[port] = list
									}
								}
								continue
							}
							recvNode, ok := nodes[rpa.Node]
							if !ok {
								continue
							}
							for _, t := range a.getPossibleNodePortTypes(scope, parentFrame, recvNode, true, rpa.Port) {
								list := c.outgoing[port]
								appendUnique(&list, t)
								c.outgoing[port] = list
							}
						}
					}

					// recursively check nested chained connections
					// in a chain like "a -> b -> c -> d", we need to check if c contains our node
					checkChainedConnections(receiver.ChainedConnection.Normal.Senders, receiver.ChainedConnection.Normal.Receivers)
				}
			}
		}
		checkChainedConnections(conn.Normal.Senders, conn.Normal.Receivers)

		// Check if our node is a receiver in this connection.
		recvPairs := a.collectReceiverSenderPairs(conn.Normal.Receivers, conn.Normal.Senders)
		for _, pair := range recvPairs {
			if pair.portAddr.Node != nodeName {
				continue
			}
			port := a.resolvePortName(nodeName, nodes, scope, true, pair.portAddr.Port)
			for _, sender := range pair.senders {
				types := a.getPossibleSenderTypes(scope, parentFrame, resolvedParentIface, nodes, sender, net)
				for _, t := range types {
					list := c.incoming[port]
					appendUnique(&list, t)
					c.incoming[port] = list
				}
			}
		}
	}

	return c
}

// flattenReceiversPortAddrs extracts all receiver port addresses from nested receiver trees,
// ignoring which sender feeds each receiver. It is useful for coarse "who can receive" scans.
//
// Examples:
//   :a -> b -> :c
//     => returns [:c] (not paired with b)
//
//   :x -> [y, z]
//     => returns [y, z]
//
// Note: use collectReceiverSenderPairs when sender/receiver pairing matters.
func (a Analyzer) flattenReceiversPortAddrs(receivers []src.ConnectionReceiver) []src.PortAddr {
	var res []src.PortAddr
	var visit func(recs []src.ConnectionReceiver)
	visit = func(recs []src.ConnectionReceiver) {
		for _, r := range recs {
			if r.PortAddr != nil {
				res = append(res, *r.PortAddr)
			}
			if r.ChainedConnection != nil && r.ChainedConnection.Normal != nil {
				// only the head receiver is relevant as a consumer of our sender
				visit(r.ChainedConnection.Normal.Receivers)
			}
			if r.DeferredConnection != nil && r.DeferredConnection.Normal != nil {
				visit(r.DeferredConnection.Normal.Receivers)
			}
		}
	}
	visit(receivers)
	return res
}

type receiverSenderPair struct {
	portAddr src.PortAddr
	senders  []src.ConnectionSender
}

// collectReceiverSenderPairs maps each receiver port to the senders that feed it.
// It preserves sender/receiver pairing across chained and deferred connections.
//
// Examples:
//
//	:a -> b -> :c
//	  => pairs: (:c <- b), and b's inport gets (:a <- in) via recursion
//
//	:start -> U::A -> switch:case[0]
//	  => pairs: (switch:case[0] <- U::A)
//
//	:x -> { :y -> :z }
//	  => pairs: (:z <- :y) for the deferred connection
func (a Analyzer) collectReceiverSenderPairs(
	receivers []src.ConnectionReceiver,
	senders []src.ConnectionSender,
) []receiverSenderPair {
	var pairs []receiverSenderPair
	// Inline recursion keeps the accumulator local and avoids extra allocations/signatures.
	var visit func(recs []src.ConnectionReceiver, snd []src.ConnectionSender)
	visit = func(recs []src.ConnectionReceiver, snd []src.ConnectionSender) {
		for _, r := range recs {
			if r.PortAddr != nil {
				pairs = append(pairs, receiverSenderPair{
					portAddr: *r.PortAddr,
					senders:  snd,
				})
				continue
			}
			if r.ChainedConnection != nil && r.ChainedConnection.Normal != nil {
				visit(r.ChainedConnection.Normal.Receivers, r.ChainedConnection.Normal.Senders)
			}
			if r.DeferredConnection != nil && r.DeferredConnection.Normal != nil {
				visit(r.DeferredConnection.Normal.Receivers, r.DeferredConnection.Normal.Senders)
			}
		}
	}
	visit(receivers, senders)
	return pairs
}

// getPossibleSenderTypes is needed to derive node constraints from the network.
// It's part of the overloading implementation.
// It returns a set of possible types produced by a given sender without requiring resolved node interfaces.
func (a Analyzer) getPossibleSenderTypes(
	scope src.Scope,
	parentFrame map[string]typesystem.Def,
	parentIface src.Interface,
	nodes map[string]src.Node,
	sender src.ConnectionSender,
	net []src.Connection,
) []typesystem.Expr {
	// FIXME: looks like we ignore errors here (and in some lower-level functions we call)

	// const sender
	if sender.Const != nil {
		// for type constraint collection, we need to get the resolved type without validation
		// since we're just collecting constraints, not validating the sender
		if sender.Const.Value.Ref != nil {
			if t, err := a.getResolvedConstTypeByRef(*sender.Const.Value.Ref, scope); err == nil {
				return []typesystem.Expr{t}
			}
		} else if sender.Const.TypeExpr.Inst != nil || sender.Const.TypeExpr.Lit != nil {
			// for literal constants, resolve the type expression directly
			if t, err := a.resolver.ResolveExpr(sender.Const.TypeExpr, scope); err == nil {
				return []typesystem.Expr{t}
			}
		}
	}

	// port-addr
	if sender.PortAddr != nil {
		// Switch:case[i] array outport slot is a special case because of possible pattern matching.
		// When T in Switch<T> is (resolves to) union, and corresponding union member has type expression body,
		// It means there is an unboxing process happening, so the output type of the Switch must be type of that member,
		// And not the type of the union itself (not T). This is kind of "ad-hoc type inference" that we must handle
		// In several different places here in analyzer, and this is one of them.
		// Please note that even though the Switch itself is NOT overloaded, it's required to cover Switch:case[i] here
		// Because it might be connected to a node that needs overloading resolution.
		// Example: `switch:case[0] -> add:left`. Without type inference (this compiler magic) output type could be union and not int,
		// Which means compiler won't be able to resolve overloading for `Add` and will throw an error that no overloading is found.
		if isSwitchCasePort(*sender.PortAddr, nodes) {
			typeExpr, err := a.getSwitchCaseOutportType(*sender.PortAddr, nodes, scope, net)
			if err != nil {
				panic(err)
			}
			return []typesystem.Expr{*typeExpr}
		}

		// If sender is input port of the current component, possible sender type can be resolved from component interface (plus frame).
		portAddr := *sender.PortAddr
		if portAddr.Node == "in" {
			if p, ok := parentIface.IO.In[portAddr.Port]; ok {
				if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, parentFrame, scope); err == nil {
					return []typesystem.Expr{resolved}
				}
			}
			return nil
		}

		// Sender is an outport of another (sub) node
		other, ok := nodes[portAddr.Node]
		if !ok {
			return nil
		}

		return a.getPossibleNodePortTypes(scope, parentFrame, other, false, portAddr.Port)
	}

	return nil
}

// getPossibleNodePortTypes collects possible port types across all overloads of a node's component.
// isInput determines whether we look at inports (true) or outports (false).
func (a Analyzer) getPossibleNodePortTypes(
	scope src.Scope,
	parentFrame map[string]typesystem.Def,
	node src.Node,
	isInput bool,
	portName string,
) []typesystem.Expr {
	var out []typesystem.Expr
	// resolve node's entity
	entity, _, err := scope.Entity(node.EntityRef)
	if err != nil || entity.Kind != src.ComponentEntity {
		return out
	}
	// resolve node type args against parent frame
	resolvedArgs, err2 := a.resolver.ResolveExprsWithFrame(node.TypeArgs, parentFrame, scope)
	if err2 != nil {
		return out
	}
	for _, comp := range entity.Component {
		iface := comp.Interface
		// choose port
		var p src.Port
		var ok bool
		if isInput {
			if portName == "" {
				if len(iface.IO.In) == 1 {
					for _, v := range iface.IO.In {
						p = v
						ok = true
						break
					}
				}
			} else {
				p, ok = iface.IO.In[portName]
			}
		} else {
			if portName == "" {
				// if multiple, prefer first non-err
				if len(iface.IO.Out) == 1 {
					for _, v := range iface.IO.Out {
						p = v
						ok = true
						break
					}
				} else {
					for name, v := range iface.IO.Out {
						if name != "err" {
							p = v
							ok = true
							break
						}
					}
				}
			} else {
				p, ok = iface.IO.Out[portName]
			}
		}
		if !ok {
			continue
		}
		// substitute node args into port type
		frame := make(map[string]typesystem.Def, len(iface.TypeParams.Params))
		for i, param := range iface.TypeParams.Params {
			if i < len(resolvedArgs) {
				frame[param.Name] = typesystem.Def{BodyExpr: &resolvedArgs[i]}
			}
		}
		if resolved, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, frame, scope); err == nil {
			out = append(out, resolved)
		}
	}
	return out
}

// doesCandidateSatisfyTypeConstraints checks that a candidate interface matches all collected constraints.
func (a Analyzer) doesCandidateSatisfyTypeConstraints(
	candidate src.Interface,
	resolvedNodeArgs []typesystem.Expr,
	nodeUsageConstr nodeUsageConstraints,
	scope src.Scope,
) bool {
	// build frame for candidate from node's resolved args
	frame := make(map[string]typesystem.Def, len(candidate.TypeParams.Params))
	for i, param := range candidate.TypeParams.Params {
		if i < len(resolvedNodeArgs) {
			frame[param.Name] = typesystem.Def{BodyExpr: &resolvedNodeArgs[i]}
		}
	}

	// check incoming constraints
	for port, types := range nodeUsageConstr.incoming {
		p, ok := candidate.IO.In[port]
		if !ok {
			return false
		}
		candType, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, frame, scope)
		if err != nil {
			return false
		}
		for _, t := range types {
			if a.resolver.IsSubtypeOf(t, candType, scope) != nil {
				return false
			}
		}
	}

	// check outgoing constraints
	for port, types := range nodeUsageConstr.outgoing {
		p, ok := candidate.IO.Out[port]
		if !ok {
			return false
		}
		candType, err := a.resolver.ResolveExprWithFrame(p.TypeExpr, frame, scope)
		if err != nil {
			return false
		}
		for _, exp := range types {
			if err := a.resolver.IsSubtypeOf(candType, exp, scope); err != nil {
				return false
			}
		}
	}

	return true
}

// isNativeComponentWithMultipleExterns checks if this is a native component with multiple extern implementations
// (like Dec with int_dec and float_dec, or Len with list_len, map_len, and string_len)
func (a Analyzer) isNativeComponentWithMultipleExterns(component src.Component, entity src.Entity) bool {
	// check if this component has extern directive
	_, hasExtern := component.Directives[compiler.ExternDirective]
	if !hasExtern {
		return false
	}

	// check if the entity has multiple components (overloaded)
	return len(entity.Component) > 1
}

// doesNativeComponentMatchTypeArgs checks if a native component's interface matches the type arguments
// this is used for disambiguating between overloaded native components like Dec(int) vs Dec(float)
func (a Analyzer) doesNativeComponentMatchTypeArgs(component src.Component, resolvedNodeArgs []typesystem.Expr, scope src.Scope) bool {
	// for native components, we need to check if the type arguments match the component's interface
	// for example, if we have Dec<int>, we need to find the component with Dec(data int) (res int)

	// if no type arguments provided, we can't disambiguate
	if len(resolvedNodeArgs) == 0 {
		return true // let other logic handle this
	}

	// resolve the first type argument to get the concrete type
	if len(resolvedNodeArgs) > 0 {
		resolvedType, err := a.resolver.ResolveExpr(resolvedNodeArgs[0], scope)
		if err != nil {
			return true // fallback to other logic
		}

		// check if this component's interface matches the resolved type
		// for native components, we need to check the input port type
		if len(component.IO.In) == 1 {
			for _, port := range component.IO.In {
				// resolve the port type
				portType, err := a.resolver.ResolveExpr(port.TypeExpr, scope)
				if err != nil {
					continue
				}

				// check if the resolved type matches the port type
				// for native components, we need exact type matching
				if a.typesMatchExactly(resolvedType, portType) {
					return true
				}
			}
		}
	}

	return false
}

// typesMatchExactly checks if two types match exactly (for native component disambiguation)
// uses string comparison as a simple and reliable equality check
func (a Analyzer) typesMatchExactly(type1, type2 typesystem.Expr) bool {
	return type1.String() == type2.String()
}

// resolvePortName resolves an empty port name to the actual port name by checking
// if the node has only one port of the given direction. returns the port name if found,
// otherwise returns the original portName (which may be empty).
func (a Analyzer) resolvePortName(
	nodeName string,
	nodes map[string]src.Node,
	scope src.Scope,
	isInput bool,
	portName string,
) string {
	// TODO figure out if this func should return error or panic, not just empty string

	if portName != "" {
		return portName
	}

	node, ok := nodes[nodeName]
	if !ok {
		return portName
	}

	entity, _, err := scope.Entity(node.EntityRef)
	if err != nil || entity.Kind != src.ComponentEntity || len(entity.Component) == 0 {
		return portName
	}

	// use first component to determine port name (works for single-port components)
	iface := entity.Component[0].Interface
	ports := iface.IO.Out
	if isInput {
		ports = iface.IO.In
	}

	// only resolve if there's exactly one port
	if len(ports) == 1 {
		for name := range ports {
			return name
		}
	}

	return portName
}
