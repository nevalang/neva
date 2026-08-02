package messages

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReflectTypeMessageRoundTrip(t *testing.T) {
	t.Parallel()

	self := int64(0)
	typeValue := ReflectType{Nodes: []ReflectTypeNode{
		{
			Kind: ReflectTypeStruct,
			Fields: []ReflectStructField{
				{Name: "text", Node: 1},
				{Name: "child", Node: 2},
			},
		},
		{Kind: ReflectTypeString},
		{
			Kind: ReflectTypeUnion,
			Cases: []ReflectUnionCase{
				{Tag: "Some", Data: &self},
				{Tag: "None"},
			},
		},
	}}

	message, err := ReflectTypeToMessage(typeValue)
	require.NoError(t, err)

	got, err := ReflectTypeFromMessage(message)
	require.NoError(t, err)
	require.Equal(t, typeValue, got)
}

func TestReflectTypeFromMessageRejectsInvalidDescriptor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message Msg
		wantErr string
	}{
		{
			name:    "not a list",
			message: NewIntMsg(0),
			wantErr: "must be a list",
		},
		{
			name:    "empty",
			message: NewUntypedListMsg(nil),
			wantErr: "root node at index 0",
		},
		{
			name: "unknown variant",
			message: NewUntypedListMsg([]Msg{
				NewUnionMsg("Unknown", nil),
			}),
			wantErr: `unknown type node variant "Unknown"`,
		},
		{
			name: "out of bounds",
			message: NewUntypedListMsg([]Msg{
				NewUnionMsg("List", NewIntMsg(1)),
			}),
			wantErr: "type node index 1 out of bounds for 1 nodes",
		},
		{
			name: "negative index",
			message: NewUntypedListMsg([]Msg{
				NewUnionMsg("Dict", NewIntMsg(-1)),
			}),
			wantErr: "type node index -1 out of bounds for 1 nodes",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := ReflectTypeFromMessage(test.message)
			require.ErrorContains(t, err, test.wantErr)
		})
	}
}

func TestReflectTypeToMessageRejectsDuplicateNames(t *testing.T) {
	t.Parallel()

	for _, typeValue := range []ReflectType{
		{Nodes: []ReflectTypeNode{{
			Kind: ReflectTypeStruct,
			Fields: []ReflectStructField{
				{Name: "value", Node: 0},
				{Name: "value", Node: 0},
			},
		}}},
		{Nodes: []ReflectTypeNode{{
			Kind: ReflectTypeUnion,
			Cases: []ReflectUnionCase{
				{Tag: "Value"},
				{Tag: "Value"},
			},
		}}},
	} {
		_, err := ReflectTypeToMessage(typeValue)
		require.ErrorContains(t, err, "duplicate")
	}
}
