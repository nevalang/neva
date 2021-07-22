package parser

type Parser interface {
	Parse([]byte) (Module, error)
}

func NewYAMLParser() Parser {
	return yamlParser{}
}

func NewJSONParser() Parser {
	return jsonParser{}
}
