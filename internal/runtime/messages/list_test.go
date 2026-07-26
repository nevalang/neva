package messages

import "testing"

func TestListToMsgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		list Msg
		want []Msg
	}{
		{
			name: "untyped",
			list: NewListMsg([]Msg{NewStringMsg("one"), NewIntMsg(2)}),
			want: []Msg{NewStringMsg("one"), NewIntMsg(2)},
		},
		{
			name: "bool",
			list: NewListBoolMsg([]bool{true, false}),
			want: []Msg{NewBoolMsg(true), NewBoolMsg(false)},
		},
		{
			name: "int",
			list: NewListIntMsg([]int64{1, 2}),
			want: []Msg{NewIntMsg(1), NewIntMsg(2)},
		},
		{
			name: "float",
			list: NewListFloatMsg([]float64{1.5, 2.5}),
			want: []Msg{NewFloatMsg(1.5), NewFloatMsg(2.5)},
		},
		{
			name: "string",
			list: NewListStringMsg([]string{"one", "two"}),
			want: []Msg{NewStringMsg("one"), NewStringMsg("two")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ListToMsgs(tt.list.List())
			if len(got) != len(tt.want) {
				t.Fatalf("ListToMsgs length = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if !Equal(got[i], tt.want[i]) {
					t.Fatalf("ListToMsgs[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestListToMsgsTypedResultDoesNotShareStorage(t *testing.T) {
	t.Parallel()

	values := []int64{1, 2}
	boxed := ListToMsgs(NewListIntMsg(values).List())
	values[0] = 99

	if got := boxed[0].Int(); got != 1 {
		t.Fatalf("boxed value = %d, want 1", got)
	}
}

func TestListToMsgsUntypedReturnsExistingStorage(t *testing.T) {
	t.Parallel()

	values := []Msg{NewIntMsg(1)}
	boxed := ListToMsgs(NewListMsg(values).List())
	boxed[0] = NewIntMsg(2)

	if got := values[0].Int(); got != 2 {
		t.Fatalf("untyped storage value = %d, want 2", got)
	}
}
