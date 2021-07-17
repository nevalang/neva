package runtime

type NativeModule struct {
	inportsInterface  InportsInterface
	outportsInterface OutportsInterface
	implementation    func(in, out map[string]chan Msg)
}

func (a NativeModule) Interface() (InportsInterface, OutportsInterface) {
	return a.inportsInterface, a.outportsInterface
}

func (a NativeModule) Run(in, out map[string]chan Msg) {
	go a.implementation(in, out)
}

type NewNativeModuleParams struct {
	inportsInterface  InportsInterface
	outportsInterface OutportsInterface
	implementation    func(in, out map[string]chan Msg)
}

func NewNativeModule(p NewNativeModuleParams) NativeModule {
	return NativeModule{
		inportsInterface:  p.inportsInterface,
		outportsInterface: p.outportsInterface,
		implementation:    p.implementation,
	}
}
