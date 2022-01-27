package platform

type (
	Compiler interface{}
	Runtime  interface{}
)

type Platform struct {
	compiler Compiler
	runtime  Runtime
}
