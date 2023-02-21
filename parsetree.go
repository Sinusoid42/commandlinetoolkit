package commandlinetoolkit

import (
	"errors"
	"fmt"
	"strings"
)

type parsetree struct {
	_depth int

	_root *argnode

	_settings *settings
}

type argnode struct {
	_arg *Argument

	_sub []*argnode

	_parent *argnode
}

func newparsetree() *parsetree {

	t := &parsetree{_depth: 0,
		_settings: &settings{}}

	t._root = newNode(nil)

	return t
}

func newNode(arg *Argument) *argnode {
	return &argnode{
		_arg:    arg,
		_sub:    []*argnode{},
		_parent: nil,
	}
}

func (p *parsetree) build(m map[string]interface{}) {

	p._settings.build(m)

	if args, ok := m["arguments"].([]interface{}); ok {

		for _, pArg := range args {

			if argMap, argOk := pArg.(map[string]interface{}); argOk {

				p._root.addArgument(argMap)

			}
		}
	}
}

func (n *argnode) addArgument(m map[string]interface{}) {

	if !checkParseableArgFromProgramFile(m) {
		return
	}

	argN := &argnode{

		_parent: n,
		_sub:    []*argnode{},
	}

	argType, err := parseArgType(m)
	if err != nil {
		return
	}

	longFlag, _ := m[LONGFLAGSTRING].(string)

	shortFlag, _ := m[SHORTFLAGSTRING].(string)

	arg := &Argument{
		arg_type: argType,
		lflag:    longFlag,
		sflag:    shortFlag,
	}

	//append the argument in a new argument node in the tree
	argN._arg = arg
	n._sub = append(n._sub, argN)

	subArgs := m[ARGUMENTSSTRING]

	if subArgs != nil {
		var ok bool

		tt, ok := subArgs.([]map[string]interface{})
		fmt.Println(tt)
		if ok {

			for _, v := range tt {

				argN.addArgument(v)
			}

		}

	}

	//manual unmarshalling here required, to check for non existing variables in the tree

}

func (p *parsetree) String() string {

	s := ""

	s += p._root.String()

	return s

}

func (n *argnode) String() string {

	s := ""
	if n._arg != nil {

		s += argtypeString(n._arg.arg_type) + ": " + n._arg.lflag + "\n"

	}

	for _, k := range n._sub {
		s += "->" + k.String() + "\n"
	}

	return s
}

func parseArgType(m map[string]interface{}) (ArgumentType, error) {

	theType := ArgumentType(0)

	typeStr := ""
	var ok bool

	if typeStr, ok = m["type"].(string); ok {

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

	if (theType&COMMAND > 0 && (theType&(OPTION|WILDCARD) > 0 || theType&PARAMETER > 0)) || (theType&(OPTION|WILDCARD) > 0 && (theType&PARAMETER > 0 || theType&COMMAND > 0)) {
		return 0, errors.New("Could not parse the Type")
	}
	return theType, nil
}

func checkParseableArgFromProgramFile(m map[string]interface{}) bool {
	if str, ok := m[TYPESTRING].(string); ok {
		if strings.Compare(str, OPTIONSTRING) == 0 {
			_, okl := m[LONGFLAGSTRING]
			_, oks := m[SHORTFLAGSTRING]
			if okl || oks {
				return true
			}

		}
		if strings.Compare(str, PARAMETERSTRING) == 0 {
			return true
		}
	}
	return false
}
