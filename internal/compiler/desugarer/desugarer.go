package desugarer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/analyzer"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrConstSenderEntityKind = errors.New("Entity that is used as a const reference in component's network must be of kind constant") //nolint:lll

type Desugarer struct{}

func (d Desugarer) Desugar(mod src.Module) (src.Module, error) {
	return mod, nil
}

// if senderSide.ConstRef != nil {
// 	constTypeExpr, err := d.getConstType(*senderSide.ConstRef, scope)
// 	if err != nil {
// 		return ts.Expr{}, Error{
// 			Location: &scope.Location,
// 			Meta:     &senderSide.ConstRef.Meta,
// 		}.Merge(err)
// 	}
// 	return constTypeExpr, nil
// }

func (d Desugarer) getConstType(ref src.EntityRef, scope src.Scope) (ts.Expr, *analyzer.Error) {
	entity, _, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &analyzer.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &analyzer.Error{
			Err:      fmt.Errorf("%w: %v", ErrConstSenderEntityKind, entity.Kind),
			Location: &scope.Location,
			Meta:     entity.Meta(),
		}
	}

	if entity.Const.Ref != nil {
		expr, err := d.getConstType(*entity.Const.Ref, scope)
		if err != nil {
			return ts.Expr{}, analyzer.Error{
				Location: &scope.Location,
				Meta:     entity.Meta(),
			}.Merge(err)
		}

		return expr, nil
	}

	return entity.Const.Value.TypeExpr, nil
}
