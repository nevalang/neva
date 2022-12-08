package analyze

import (
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
)

type Comparator struct{}

func (c Comparator) Compare(a, b src.TypeExpr) error {
	isAStruct := a.Struct != nil
	isBStruct := b.Struct != nil

	if isAStruct && !isBStruct {
		return errors.New("first type is struct")
	}
	if isBStruct && !isAStruct {
		return errors.New("second type is struct")
	}

	if !isAStruct {
		if a.Ref != b.Ref || len(a.RefArgs) != len(b.RefArgs) {
			return errors.New("")
		}
		for i := range a.RefArgs {
			if err := c.Compare(
				a.RefArgs[i],
				b.RefArgs[i],
			); err != nil {
				return err
			}
		}
		return nil
	}

	if len(a.Struct) != len(b.Struct) {
		return errors.New("")
	}

	for field := range a.Struct {
		if err := c.Compare(a.Struct[field], b.Struct[field]); err != nil {
			return errors.New("")
		}
	}

	return nil
}
