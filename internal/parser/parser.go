package parser

type Parser interface {
	Parse([]byte) (Module, error)
}

func NewYAMLParser(v Validator) Parser {
	return yamlParser{v}
}

func NewJSONParser(v Validator) Parser {
	return jsonParser{v}
}
