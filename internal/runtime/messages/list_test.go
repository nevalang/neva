package messages

import "testing"

func TestListLen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		list Msg
		name string
		want int
	}{
		{name: "untyped", list: NewListMsg([]Msg{NewIntMsg(1), NewIntMsg(2)}), want: 2},
		{name: "bool", list: NewListBoolMsg([]bool{true}), want: 1},
		{name: "int", list: NewListIntMsg([]int64{1, 2, 3}), want: 3},
		{name: "float", list: NewListFloatMsg([]float64{}), want: 0},
		{name: "string", list: NewListStringMsg([]string{"one", "two"}), want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ListLen(tt.list.List()); got != tt.want {
				t.Fatalf("ListLen() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestListToMessageSlice(t *testing.T) {
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

			got := ListToMessageSlice(tt.list.List())
			if len(got) != len(tt.want) {
				t.Fatalf("ListToMessageSlice length = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if !Equal(got[i], tt.want[i]) {
					t.Fatalf("ListToMessageSlice[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestListToMessageSliceTypedResultDoesNotShareStorage(t *testing.T) {
	t.Parallel()

	values := []int64{1, 2}
	boxed := ListToMessageSlice(NewListIntMsg(values).List())
	values[0] = 99

	if got := boxed[0].Int(); got != 1 {
		t.Fatalf("boxed value = %d, want 1", got)
	}
}

func TestListToMessageSliceUntypedReturnsExistingStorage(t *testing.T) {
	t.Parallel()

	values := []Msg{NewIntMsg(1)}
	boxed := ListToMessageSlice(NewListMsg(values).List())
	boxed[0] = NewIntMsg(2)

	if got := values[0].Int(); got != 2 {
		t.Fatalf("untyped storage value = %d, want 2", got)
	}
}

func TestListAt(t *testing.T) {
	t.Parallel()

	list := NewListIntMsg([]int64{10, 20, 30}).List()
	item, found := ListAt(list, -2)
	if !found || item.Int() != 20 {
		t.Fatalf("ListAt(-2) = (%v, %t), want (20, true)", item, found)
	}

	_, found = ListAt(list, 3)
	if found {
		t.Fatal("ListAt(3) found = true, want false")
	}
}

func TestListSlicePreservesRepresentationAndCopiesStorage(t *testing.T) {
	t.Parallel()

	values := []int64{1, 2, 3, 4, 5}
	sliced := ListSlice(NewListIntMsg(values).List(), -3, 99)
	result, ok := ListAsInts(sliced.List())
	if !ok {
		t.Fatal("ListSlice() did not preserve typed integer storage")
	}
	if got, want := len(result), 3; got != want {
		t.Fatalf("ListSlice length = %d, want %d", got, want)
	}
	if result[0] != 3 || result[2] != 5 {
		t.Fatalf("ListSlice values = %v, want [3 4 5]", result)
	}

	values[2] = 99
	if result[0] != 3 {
		t.Fatal("ListSlice result shares backing storage with source list")
	}

	empty := ListSlice(NewListMsg([]Msg{NewIntMsg(1), NewIntMsg(2)}).List(), 2, 1)
	if got := ListLen(empty.List()); got != 0 {
		t.Fatalf("ListSlice reversed range length = %d, want 0", got)
	}
}

func TestListAppendPreservesCompatibleTypedStorage(t *testing.T) {
	t.Parallel()

	typed := ListAppend(NewListIntMsg([]int64{1, 2}).List(), NewIntMsg(3))
	values, ok := ListAsInts(typed.List())
	if !ok || len(values) != 3 || values[2] != 3 {
		t.Fatalf("ListAppend typed result = %v, want typed [1 2 3]", typed)
	}

	mixed := ListAppend(NewListIntMsg([]int64{1}).List(), NewStringMsg("two"))
	if _, ok := ListAsUntyped(mixed.List()); !ok {
		t.Fatal("ListAppend incompatible value did not produce untyped storage")
	}
}

func TestListPrependPreservesCompatibleTypedStorage(t *testing.T) {
	t.Parallel()

	values := []int64{2, 3}
	typed := ListPrepend(NewListIntMsg(values).List(), NewIntMsg(1))
	result, ok := ListAsInts(typed.List())
	if !ok || len(result) != 3 || result[0] != 1 || result[1] != 2 || result[2] != 3 {
		t.Fatalf("ListPrepend typed result = %v, want typed [1 2 3]", typed)
	}

	values[0] = 99
	if result[1] != 2 {
		t.Fatal("ListPrepend result shares backing storage with source list")
	}

	mixed := ListPrepend(NewListIntMsg([]int64{1}).List(), NewStringMsg("zero"))
	if _, ok := ListAsUntyped(mixed.List()); !ok {
		t.Fatal("ListPrepend incompatible value did not produce untyped storage")
	}
}
