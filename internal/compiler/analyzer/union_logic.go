package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

type unionActiveTagInfo struct {
	tag         string
	tagTypeExpr ts.Expr
}

// buildUnionActiveTagBindings inspects the network wiring to determine which concrete tag
// is active for each Union<T> node (via its :tag port). It returns a map from Union node name
// to the active tag and its payload type, so Union:data connections can be validated.
func (a Analyzer) buildUnionActiveTagBindings(
	net []src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (map[string]unionActiveTagInfo, *compiler.Error) {
	infos := map[string]unionActiveTagInfo{}

	for nodeName, node := range nodes {
		// Only Union<T> nodes participate in tag/data validation.
		if !a.isUnionNode(node) {
			continue
		}

		// Type args are already resolved in analyzeNode; we only need to confirm it's a union.
		resolvedUnionTypeArg := node.TypeArgs[0]
		if resolvedUnionTypeArg.Lit == nil || resolvedUnionTypeArg.Lit.Union == nil {
			return nil, &compiler.Error{
				Message: "Union<T> expects union type argument",
				Meta:    &node.Meta,
			}
		}

		// Gather all senders connected to Union:tag (including chained connections).
		tagSenders := a.collectPortSenders(net, nodeName, "tag")
		if len(tagSenders) == 0 {
			return nil, &compiler.Error{
				Message: fmt.Sprintf("Union tag input is missing for node %q", nodeName),
				Meta:    &node.Meta,
			}
		}

		// Artificial restriction: disallow fan-in on Union:tag to keep semantics simple.
		// In practice, Union:tag is a single constant sender (often a literal).
		// This covers the vast majority of valid cases while avoiding ambiguous wiring.
		if len(tagSenders) != 1 {
			return nil, &compiler.Error{
				Message: fmt.Sprintf("Union tag input must have exactly one sender for node %q", nodeName),
				Meta:    &node.Meta,
			}
		}

		// Validate the single tag tagSender for this node and record its tag + payload type.
		// Tag tagSender must be a union literal (direct or const ref).
		tagSender := tagSenders[0]
		unionLiteral, err := a.resolveUnionConstSender(tagSender, scope)
		if err != nil {
			return nil, compiler.Error{Meta: &tagSender.Meta}.Wrap(err)
		}

		// Ensure literal's union type matches Union<T>.
		_, err = a.resolveUnionTypeFromLiteral(unionLiteral, scope)
		if err != nil {
			return nil, compiler.Error{Meta: &tagSender.Meta}.Wrap(err)
		}

		// Validate tag exists on the union and has a payload type.
		memberTypeExpr, ok := resolvedUnionTypeArg.Lit.Union[unionLiteral.Tag]
		if !ok {
			return nil, &compiler.Error{
				Message: fmt.Sprintf(
					"tag %q not found in %v",
					unionLiteral.Tag,
					resolvedUnionTypeArg,
				),
				Meta: &tagSender.Meta,
			}
		}
		if memberTypeExpr == nil {
			return nil, &compiler.Error{
				Message: fmt.Sprintf(
					"tag %q requires a value of a defined type",
					unionLiteral.Tag,
				),
				Meta: &tagSender.Meta,
			}
		}

		infos[nodeName] = unionActiveTagInfo{
			tag:         unionLiteral.Tag,
			tagTypeExpr: *memberTypeExpr,
		}
	}

	return infos, nil
}

// validateUnionDataReceiverPort checks a single port address and applies the Union:data constraint
// if the receiver is a Union node with a resolved active tag type.
func (a Analyzer) validateUnionDataReceiverPort(
	portAddr src.PortAddr,
	senders []src.ConnectionSender,
	resolvedSenderTypes []*ts.Expr,
	unionActiveTags map[string]unionActiveTagInfo,
	scope src.Scope,
) *compiler.Error {
	tagInfo, ok := unionActiveTags[portAddr.Node]
	if !ok {
		// Not a Union<T> node or no tag info recorded.
		return nil
	}

	// All senders must be subtype-compatible with the selected tag payload type.
	for i, senderType := range resolvedSenderTypes {
		if err := a.resolver.IsSubtypeOf(*senderType, tagInfo.tagTypeExpr, scope); err != nil {
			return &compiler.Error{
				Message: fmt.Sprintf(
					"Union data type %v is not compatible with tag %q (%v): %v",
					senderType,
					tagInfo.tag,
					tagInfo.tagTypeExpr,
					err,
				),
				Meta: &senders[i].Meta,
			}
		}
	}

	return nil
}

// collectPortSenders finds all senders that target nodeName:port in the network.
// This includes direct receivers and nested receivers inside chained connections.
func (a Analyzer) collectPortSenders(
	net []src.Connection,
	nodeName string,
	port string,
) []src.ConnectionSender {
	var out []src.ConnectionSender

	for _, conn := range net {
		a.collectPortSendersInReceivers(
			conn.Senders,
			conn.Receivers,
			nodeName,
			port,
			&out,
		)
	}

	return out
}

// collectPortSendersInReceivers is the recursive worker that walks receiver trees
// and appends matching senders into out.
func (a Analyzer) collectPortSendersInReceivers(
	senders []src.ConnectionSender,
	receivers []src.ConnectionReceiver,
	nodeName string,
	port string,
	out *[]src.ConnectionSender,
) {
	for _, receiver := range receivers {
		// Direct receiver: outer senders feed nodeName:port.
		if receiver.PortAddr != nil {
			if receiver.PortAddr.Node == nodeName && receiver.PortAddr.Port == port {
				*out = append(*out, senders...)
			}
		}

		// Chained receiver: chain head is the receiver for outer senders.
		if receiver.ChainedConnection != nil {
			for _, chainHead := range receiver.ChainedConnection.Senders {
				if chainHead.PortAddr != nil &&
					chainHead.PortAddr.Node == nodeName &&
					chainHead.PortAddr.Port == port {
					*out = append(*out, senders...)
				}
			}

			// Recurse into the chain to catch deeper receivers.
			a.collectPortSendersInReceivers(
				receiver.ChainedConnection.Senders,
				receiver.ChainedConnection.Receivers,
				nodeName,
				port,
				out,
			)
		}
	}
}

// isUnionNode identifies the builtin Union<T> node used for wrapping tags/payloads.
func (a Analyzer) isUnionNode(node src.Node) bool {
	return node.EntityRef.Name == "Union" && (node.EntityRef.Pkg == "" || node.EntityRef.Pkg == "builtin")
}

// resolveUnionConstSender extracts a union literal from a sender.
// Valid senders are union literal constants or references to constants that resolve to union literals.
func (a Analyzer) resolveUnionConstSender(
	sender src.ConnectionSender,
	scope src.Scope,
) (*src.UnionLiteral, *compiler.Error) {
	if sender.Const == nil {
		return nil, &compiler.Error{
			Message: "Union tag sender must be a union constant",
			Meta:    &sender.Meta,
		}
	}

	if sender.Const.Value.Message != nil && sender.Const.Value.Message.Union != nil {
		return sender.Const.Value.Message.Union, nil
	}

	if sender.Const.Value.Ref != nil {
		return a.resolveUnionConstRef(*sender.Const.Value.Ref, scope)
	}

	return nil, &compiler.Error{
		Message: "Union tag sender must be a union constant",
		Meta:    &sender.Meta,
	}
}

// resolveUnionTypeFromLiteral resolves the union type definition referenced by a union literal.
// It returns the fully analyzed union type expression or an error if the reference isn't a union.
func (a Analyzer) resolveUnionTypeFromLiteral(
	unionLiteral *src.UnionLiteral,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	typeDef, _, err := scope.GetType(unionLiteral.EntityRef)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf("failed to resolve union type: %v", err),
			Meta:    &unionLiteral.Meta,
		}
	}

	unionTypeExpr, analyzeExprErr := a.analyzeTypeExpr(*typeDef.BodyExpr, scope)
	if analyzeExprErr != nil {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf("failed to resolve union type: %v", analyzeExprErr),
			Meta:    &unionLiteral.Meta,
		}
	}

	if unionTypeExpr.Lit == nil || unionTypeExpr.Lit.Union == nil {
		return ts.Expr{}, &compiler.Error{
			Message: fmt.Sprintf("type %v is not a union", unionLiteral.EntityRef),
			Meta:    &unionLiteral.Meta,
		}
	}

	return unionTypeExpr, nil
}
