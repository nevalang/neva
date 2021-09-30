package storage

type (
	Store interface {
		Module(ModuleParams) ([]byte, error)
	}

	ModuleParams struct {
		id                  string
		major, minor, patch uint64
	}
)
