package types_test

import (
	"reflect"
	"testing"

	"github.com/emil14/neva/pkg/types"
)

func TestResolver_Resolve(t *testing.T) {
	tests := []struct {
		name    string
		expr    types.Expr
		scope   map[string]types.Def
		want    types.Expr
		wantErr bool
	}{
		{},
	}

	r := types.NewResolver(nil, nil) // TODO mocks

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Resolve(tt.expr, tt.scope)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
