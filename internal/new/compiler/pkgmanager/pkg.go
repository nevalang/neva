package pkgmanager

type (
	Pkg struct {
		Imports             Imports
		Scope               map[string]ImportRef
		RootModule          string
		Exports             []string
		WantCompilerVersion string
	}

	Imports struct {
		Local  []string
		StdLib map[string]string
		Global map[string]GlobalImport
	}

	ImportRef struct {
		NameSpace NameSpace
		Pkg       string
		Component string
	}

	GlobalImport struct {
		Pkg     string
		Version string
	}

	NameSpace uint8
)

const (
	StdLib NameSpace = iota + 1
	Local
	Global
)
