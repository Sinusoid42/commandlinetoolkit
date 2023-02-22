package commandlinetoolkit

import "fmt"

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

/*
**********************************************************************
Build the tree
*/
func (p *parsetree) build(m map[string]interface{}) {

	p._settings.build(m)

	args, ok := m[ARGUMENTSKEY].([]map[string]interface{})

	if !ok {
		return
	}

	for _, arg := range args {

		p._root.addArgument(arg)

	}

	fmt.Println(p)
	//now we expect that every part from the input is well formattedd
	//map[string]interface{} etc .. 'm["arguments"].(type) = []map[string]interface{}' !!

	/*if args, ok := m["arguments"].([]interface{}); ok {

		for _, pArg := range args {

			if argMap, argOk := pArg.(map[string]interface{}); argOk {

				p._root.addArgument(argMap)

			}
		}
	}*/
}

/*
**********************************************************************
Add a new Argument as a result of traversion the input tree
*/
func (n *argnode) addArgument(m map[string]interface{}) {

	/*	if !checkParseableArgFromProgramFile(m) {
		return
	}*/

	argN := &argnode{

		_parent: n,
		_sub:    []*argnode{},
	}

	argType, err := parseArgType(m)
	if err != nil {
		return
	}

	longFlag, _ := m[LONGFLAGKEY].(string)

	shortFlag, _ := m[SHORTFLAGKEY].(string)

	help, _ := m[HELPKEY].(string)

	shelp, _ := m[SHORTHELPKEY].(string)

	arg := &Argument{
		arg_type: argType,
		lflag:    longFlag,
		sflag:    shortFlag,
		lhelp:    help,
		shelp:    shelp,
	}

	//append the argument in a new argument node in the tree
	argN._arg = arg
	n._sub = append(n._sub, argN)

	subArgs := m[ARGUMENTSKEY]

	if subArgs != nil {
		var ok bool

		tt, ok := subArgs.([]map[string]interface{})

		if ok {

			for _, v := range tt {

				argN.addArgument(v)
			}

		}

	}

	//manual unmarshalling here required, to check for non existing variables in the tree

}
