package parser

type (
	rawPkgDescriptor struct {
		Import rawPkgImports     `yaml:"import,required"`
		Scope  map[string]string `yaml:"scope,required"`
		Meta   rawPkgMeta        `yaml:"meta,required"`
		Exec   string            `yaml:"exec"`
		Export string            `yaml:"export"`
	}

	rawPkgImports struct {
		Std    map[string]string `yaml:"std"`
		Global map[string]string `yaml:"global"`
		Local  []string          `yaml:"local"`
	}

	rawPkgMeta struct {
		CompilerVersion string `yaml:"compiler"`
	}
)
