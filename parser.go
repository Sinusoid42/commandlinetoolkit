package commandlinetoolkit

import "errors"

// here we build and store the pseudo tree for the parser
type parser struct {

	//the tree of all the nodes from the configuration
	_parseTree *parsetree

	Executeable string
}

func newparser() *parser {
	return &parser{_parseTree: newparsetree()}
}

// from here we also generate the JSON or read the JSON
// builds the commandline tree dynamicall from the json
func (p *parser) parseProgram(prgm *program) (string, error) {

	//after checking we most definately have a well formed tree
	p._parseTree.build(prgm._program._theProgramJsonMap)

	//first parse root
	//recursively parse all other sub-arguments for n-node tree
	return "", errors.New("todo")
}

type StringParseable interface {
	GetParseArgs() []string
}

type KeyParseable interface {
	GetParseKeys() []Key
}

func check(query string, def string, m map[string]interface{}) string {
	_str := ""
	if str, ok := m[query].(string); ok {
		_str = str
	} else {
		_str = def
	}
	return _str
}
