package main

type (
	Module struct {
		GenericParams []string // TypeInstance?
		IO            struct {
			In  map[string]Port
			Out map[string]Port
		}
		Nodes struct {
			StructBuilders struct {
				MsgType TypeInstance // Should be struct
			}
			Consts map[string]struct {
				MsgType  TypeInstance
				MsgValue MsgValue
			}
			// Workers        map[string]???
		}
	}

	Port struct {
		IsArray bool
		MsgType TypeInstance
	}

	TypeInstance struct {
		TypeRef     string
		GenericArgs []TypeInstance
	}

	MsgValue struct {
		BaseType    BaseType
		BoolValue   bool
		IntValue    int64
		FloatValue  float64
		StrValue    string
		ListValue   []MsgValue
		DictValue   map[string]MsgValue
		StructValue map[string]MsgValue
	}
)

type (
	BaseType uint8 // bool, int, float, string, list, struct

	BoolType  struct{}
	IntType   struct{}
	FloatType struct{}
	StrType   struct{}

	ListType   struct{ ValueType TypeInstance }
	Dict       struct{ ValueType TypeInstance }
	StructType map[string]TypeInstance
)

func main() {
}
