package types

import "testing"

func TestExpr_IsSubType(t *testing.T) {
	tests := []struct {
		name    string
		expr Expr
		constraint Expr
		wantErr bool
	}{
		{
			""
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			if err := expr.IsSubType(tt.args.constraint); (err != nil) != tt.wantErr {
				t.Errorf("Expr.IsSubType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
