package commandlinetoolkit

import `encoding/json`

//here we build the commandline tree and structure from the program json
type builder struct {
}

func (b *builder) Rebuild(c *CommandLine) {

}

//from here we also generate the JSON or read the JSON
//builds the commandline tree dynamicall from the json
//
func (p *parser) parse(commandlinejson string, c *CommandLine) {
	
	m := map[string]string{}
	
	json.Unmarshal([]byte(commandlinejson), &m)
	
}
