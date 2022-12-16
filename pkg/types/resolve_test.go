package types_test

import (
	"reflect"
	"testing"

	"github.com/emil14/neva/pkg/types"
)

func Test_resolve(t *testing.T) {
	type args struct {
		expr  types.Expr
		scope map[string]types.Def
	}

	tests := []struct {
		name    string
		args    args
		want    types.Expr
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := types.Resolve(tt.args.expr, tt.args.scope)
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
