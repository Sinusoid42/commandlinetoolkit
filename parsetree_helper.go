package commandlinetoolkit

import (
	"errors"
	"strings"
)

/*
**********************************************************************
Return the Tree visually
*/
func (p *parsetree) String() string {
	
	s := "Parsetree:\n"
	
	s += p._root.String("  ")
	
	return s
	
}

func (n *argnode) String(str string) string {
	
	s := ""
	if n._arg != nil {
		
		s += argtypeString(n._arg.arg_type) + ": " + n._arg.lflag + "\n"
		s += str + ">flag: " + n._arg.lflag + "\n"
		s += str + ">help: " + n._arg.lhelp + "\n"
		
	}
	
	for _, k := range n._sub {
		s += str + "->" + k.String(str+"   ") + "\n"
	}
	
	return " " + s
}

/*
**********************************************************************
parse a girven Argument Type
*/
func parseArgType(m map[string]interface{}) (ArgumentType, error) {
	
	theType := ArgumentType(0)
	
	typeStr := ""
	var ok bool
	if typeStr, ok = m[TYPEKEY].(string); ok {
		
		typeStrArrOR := strings.Split(typeStr, "|")
		typeStrArrComma := strings.Split(typeStr, ",")
		
		theTypes := []string{}
		
		if len(typeStrArrOR) > 0 {
			theTypes = typeStrArrOR
			
		} else if len(typeStrArrComma) > 0 {
			theTypes = typeStrArrComma
		}
		
		for _, value := range theTypes {
			
			if strings.Compare(value, OPTIONSTRING) == 0 {
				theType |= OPTION
			}
			if strings.Compare(value, PARAMETERSTRING) == 0 {
				theType |= PARAMETER
			}
			if strings.Compare(value, WILDCARDSTRING) == 0 {
				theType |= WILDCARD
			}
			if strings.Compare(value, FLAGSTRING) == 0 {
				theType |= FLAG
			}
			if strings.Compare(value, COMMANDSTRING) == 0 {
				theType |= COMMAND
			}
		}
	}
	
	if theType&COMMAND > 0 && (theType&(OPTION|WILDCARD|PARAMETER) > 0) {
		return theType, errors.New("Could not parse the Type")
	}
	return theType, nil
}
