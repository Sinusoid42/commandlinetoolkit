package commandlinetoolkit

import `encoding/json`

//here we build and store the pseudo tree for the parser
//from here we also generate the JSON or read the JSON
type parser struct {
	Executeable string
}

//builds the commandline tree dynamicall from the json
//
func (p *parser) parse(commandlinejson string, c *CommandLine) {
	
	m := map[string]string{}
	
	json.Unmarshal([]byte(commandlinejson), &m)
	
}

func newparser() parser {
	
	return parser{}
}
