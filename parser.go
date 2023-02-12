package commandlinetoolkit

// here we build and store the pseudo tree for the parser
type parser struct {
	Executeable string
}

func newparser() *parser {
	return &parser{}
}

type StringParseable interface {
	GetParseArgs() []string
}

type KeyParseable interface {
	GetParseKeys() []Key
}
