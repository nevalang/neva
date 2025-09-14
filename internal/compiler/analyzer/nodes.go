package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type foundInterface struct {
	iface    src.Interface
	location core.Location
}

func (a Analyzer) analyzeNodes(
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

func (a Analyzer) analyzeNode(
	name string, // name of the node
	node src.Node, // node to analyze
	scope src.Scope, // scope of the component that contains the node
	iface src.Interface, // interface of the component that contains the node
	nodes map[string]src.Node, // nodes of the component that contains the node
	net []src.Connection, // network of the component that contains the node
) (src.Node, foundInterface, *compiler.Error) {
	parentTypeParams := iface.TypeParams

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
	if hasBind && len(bindArg) != 1 {
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
		nodeIface, overloadIndex, err = a.getComponentNodeInterface(
			name,
			nodeEntity,
			hasBind,
			node,
			scope,
			resolvedNodeArgs,
			iface,
			// nodes,
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
		resolvedDep, _, err := a.analyzeNode(
			name, // TODO make sure DI works with overloading (example: Reduce{Add})
			depNode,
			scope,
			iface,
			nodes,
			net,
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

// getComponentNodeInterface returns interface of the component node.
// It also performs some validation.
// Overloading at the level of sourcecode is implemented here.
func (a Analyzer) getComponentNodeInterface(
	name string,
	entity src.Entity,
	hasBind bool,
	node src.Node,
	scope src.Scope,
	resolvedNodeArgs []typesystem.Expr,
	parentIface src.Interface, // resolved interface of the component that contains the node
	// nodes map[string]src.Node, // nodes of the component that contains the node
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
		version, overloadIndex, err = a.getNodeOverloadIndex(
			name,
			parentIface,
			// nodes,
			net,
			entity.Component,
			// scope,
		)
		if err != nil {
			return src.Interface{}, nil, &compiler.Error{
				Message: "Node can't use #bind if it isn't instantiated with the component that use #extern",
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

// getNodeOverloadIndex determines which overload of a component to use for a node.
// This is called when we know the node references an overloaded component with multiple implementations.
// It analyzes how the node is used in connections to determine the appropriate implementation.
func (a Analyzer) getNodeOverloadIndex(
	name string, // name of the node
	parentIface src.Interface, // resolved interface of the component that contains the node
	// nodes map[string]src.Node, // nodes of the component that contains the node
	net []src.Connection, // network of the component that contains the node
	versions []src.Component, // all versions of the component that node refers to
	// scope src.Scope,
) (src.Component, *int, *compiler.Error) {
	// We'll analyze all connections to find where this node is used
	nodeUsages := findNodeUsages(name, net)
	if len(nodeUsages) == 0 {
		return src.Component{}, nil, &compiler.Error{
			Message: fmt.Sprintf("no usages found for node %s", name),
			Meta:    &parentIface.Meta,
		}
	}

	// For each overload, check if it's compatible with all the node's usages
	for i, component := range versions {
		if isCompatibleWithAllUsages(component, nodeUsages) {
			return component, &i, nil
		}
	}

	return src.Component{}, nil, &compiler.Error{
		Message: fmt.Sprintf("no compatible overload found for node %s", name),
		Meta:    &parentIface.Meta,
	}
}

// findNodeUsages identifies all places where the specified node is used in connections
func findNodeUsages(nodeName string, connections []src.Connection) []nodeUsage {
	var usages []nodeUsage

	for _, conn := range connections {
		// Check array bypass connections
		if conn.ArrayBypass != nil {
			if conn.ArrayBypass.SenderOutport.Node == nodeName {
				usages = append(usages, nodeUsage{
					direction: outgoing,
					port:      conn.ArrayBypass.SenderOutport.Port,
					arrayIdx:  conn.ArrayBypass.SenderOutport.Idx,
				})
			}
			if conn.ArrayBypass.ReceiverInport.Node == nodeName {
				usages = append(usages, nodeUsage{
					direction: incoming,
					port:      conn.ArrayBypass.ReceiverInport.Port,
					arrayIdx:  conn.ArrayBypass.ReceiverInport.Idx,
				})
			}
			continue
		}

		// Check normal connections
		if conn.Normal != nil {
			// Check senders
			for _, sender := range conn.Normal.Senders {
				if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
					usages = append(usages, nodeUsage{
						direction: outgoing,
						port:      sender.PortAddr.Port,
						arrayIdx:  sender.PortAddr.Idx,
					})
				}
			}

			// Check receivers (more complex due to nesting)
			usages = append(usages, findNodeUsagesInReceivers(nodeName, conn.Normal.Receivers)...)
		}
	}

	return usages
}

type connectionDirection int

const (
	incoming connectionDirection = iota
	outgoing
)

type nodeUsage struct {
	direction connectionDirection
	port      string
	arrayIdx  *uint8
}

// findNodeUsagesInReceivers recursively checks receivers for node usages
func findNodeUsagesInReceivers(nodeName string, receivers []src.ConnectionReceiver) []nodeUsage {
	var usages []nodeUsage

	for _, receiver := range receivers {
		// Check direct port address
		if receiver.PortAddr != nil && receiver.PortAddr.Node == nodeName {
			usages = append(usages, nodeUsage{
				direction: incoming,
				port:      receiver.PortAddr.Port,
				arrayIdx:  receiver.PortAddr.Idx,
			})
		}

		// Check chained connection
		if receiver.ChainedConnection != nil {
			if receiver.ChainedConnection.Normal != nil {
				// Check senders in the chain
				for _, sender := range receiver.ChainedConnection.Normal.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						usages = append(usages, nodeUsage{
							direction: outgoing,
							port:      sender.PortAddr.Port,
							arrayIdx:  sender.PortAddr.Idx,
						})
					}
				}

				// Recursively check receivers in the chain
				usages = append(usages, findNodeUsagesInReceivers(nodeName, receiver.ChainedConnection.Normal.Receivers)...)
			}
		}

		// Check deferred connection
		if receiver.DeferredConnection != nil {
			// Similar logic to what we do with normal connections
			if receiver.DeferredConnection.Normal != nil {
				for _, sender := range receiver.DeferredConnection.Normal.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						usages = append(usages, nodeUsage{
							direction: outgoing,
							port:      sender.PortAddr.Port,
							arrayIdx:  sender.PortAddr.Idx,
						})
					}
				}

				usages = append(usages, findNodeUsagesInReceivers(nodeName, receiver.DeferredConnection.Normal.Receivers)...)
			}
		}

		// Check switch cases
		if receiver.Switch != nil {
			// Check each case in the switch
			for _, caseConn := range receiver.Switch.Cases {
				for _, sender := range caseConn.Senders {
					if sender.PortAddr != nil && sender.PortAddr.Node == nodeName {
						usages = append(usages, nodeUsage{
							direction: outgoing,
							port:      sender.PortAddr.Port,
							arrayIdx:  sender.PortAddr.Idx,
						})
					}
				}

				usages = append(usages, findNodeUsagesInReceivers(nodeName, caseConn.Receivers)...)
			}

			// Check default case
			if receiver.Switch.Default != nil {
				usages = append(usages, findNodeUsagesInReceivers(nodeName, receiver.Switch.Default)...)
			}
		}
	}

	return usages
}

// isCompatibleWithAllUsages checks if the given component is compatible with all the usages
func isCompatibleWithAllUsages(component src.Component, usages []nodeUsage) bool {
	for _, usage := range usages {
		var portExists bool
		var port src.Port

		if usage.direction == incoming {
			// Check if the inport exists in the component
			port, portExists = component.Interface.IO.In[usage.port]
		} else {
			// Check if the outport exists in the component
			port, portExists = component.Interface.IO.Out[usage.port]
		}

		// If the port doesn't exist or array usage doesn't match, not compatible
		if !portExists || port.IsArray != (usage.arrayIdx != nil) {
			return false
		}
	}
	return true
}
