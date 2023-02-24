package commandlinetoolkit

import (
	"fmt"
	"os"
	"os/exec"
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

/*
**********************************************************************
Build the tree
*/
func (p *parsetree) build(m map[string]interface{}) bool {

	p._settings.build(m)

	args, ok := m[ARGUMENTSKEY].([]map[string]interface{})

	if !ok {
		return false
	}

	for _, arg := range args {

		_ok := p._root.addArgument(arg)
		if !_ok {
			ok = ok
		}

	}

	return ok
}

func tokenize(p *parsetree, args []string) (*parsetree, CLICODE) {
	np := newparsetree()
	r := p._root
	if r._sub == nil {
		return np, CLI_LEXING_ERROR
	}
	code := CLI_ERROR

	if len(args) == 0 {
		return np, CLI_SUCCESS
	}

	newArg, ok := r.tokenizeArg(args)
	if !ok {
		//fmt.Println(newArg)
	} else {
		code = CLI_SUCCESS
	}
	np._root = newArg

	//fmt.Println("NEW Parsetree", np)
	return np, code
}

func (a *argnode) tokenizeArg(args []string) (*argnode, bool) {

	node := newNode(nil)

	if a._sub == nil || len(a._sub) == 0 {
		return nil, false
	}

	index := 0

	arg := &argnode{}
	none := false
	for a._sub != nil && len(a._sub) > 0 {
		if index >= len(a._sub) {

			if len(args) > 0 && !none {
				none = true
				index = 0
			} else {

				break
			}

		}
		arg = a._sub[index]
		if arg._arg.arg_type&OPTION > 0 {

			//fmt.Println("TEST: " + arg._arg.lflag)
			newarg, newargs, ok := tokenizeOption(arg, args)
			if a != nil && newarg != nil && ok {
				none = false
				newarg._parent = a
				node._sub = append(node._sub, newarg)

			}

			if ok && len(newargs) == 0 {

				return node, true
			}

			args = newargs

		}
		if arg._arg.arg_type&COMMAND > 0 {
			//fmt.Println("DEBUG:")
			//fmt.Println("TEST: " + arg._arg.lflag)
			newarg, newargs, ok := tokenizeCommand(arg, args)
			if a != nil && newarg != nil && ok {
				none = false
				newarg._parent = a
				node._sub = append(node._sub, newarg)

			}

			if ok && len(newargs) == 0 {

				return node, true
			}

			args = newargs

		}

		index++

	}
	return node, false
}

func tokenizeOption(a *argnode, args []string) (*argnode, []string, bool) {
	ok := false
	if len(args) <= 0 {
		return nil, args, false
	}

	if len(a._arg.lflag) > 0 && strings.Index(args[0], FULLOPTIONPREFIX+a._arg.lflag) == 0 {
		ok = true
	}
	if len(a._arg.sflag) > 0 && strings.Index(args[0], SHORTOPTIONPREFIX+a._arg.sflag) == 0 {
		ok = true
	}
	if !ok {
		return nil, args, false
	}

	index := 1

	newargnode := newNode(a._arg.copy())

	//fmt.Println(args)

	if len(args) == 0 {
		//split on =
		data, ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index])

		if ok {
			newargnode._arg.data_type.data = data
		} else {
		}
	}
	for index = 1; index < len(args); index++ {

		if len(a._arg.data_type.data_flag) > 0 {

			data, ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index])

			if ok {
				newargnode._arg.data_type.data = data
			} else {
				break
			}

		}
		if a._sub == nil || len(a._sub) == 0 {
			break
		}
	}
	newargs := args[index:]
	return newargnode, newargs, ok
}

func tokenizeCommand(a *argnode, args []string) (*argnode, []string, bool) {
	ok := false
	if len(args) <= 0 {
		return nil, args, false
	}

	if len(a._arg.lflag) > 0 && strings.Index(args[0], a._arg.lflag) == 0 {
		ok = true
	}
	if !ok {
		return nil, args, false
	}

	index := 1

	newargnode := newNode(a._arg.copy())
	//fmt.Println(args)

	if len(args) == 0 {
		data, ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index])

		if ok {
			newargnode._arg.data_type.data = data
		} else {
		}
	}
	for index = 1; index < len(args); index++ {

		if len(a._arg.data_type.data_flag) > 0 {

			data, ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index])

			if ok {
				newargnode._arg.data_type.data = data
			} else {
				break
			}

		}
		if a._sub == nil || len(a._sub) == 0 {
			break
		}
	}
	newargs := args[index:]
	return newargnode, newargs, ok
}

/*
**********************************************************************
Add a new Argument as a result of traversion the input tree
*/
func (n *argnode) addArgument(m map[string]interface{}) bool {

	/*	if !checkParseableArgFromProgramFile(m) {
		return
	}*/

	//create a node
	argN := &argnode{

		_parent: n,
		_sub:    []*argnode{},
	}

	arg, err := NewArgument(m)

	if err != nil {
		fmt.Println(err)
		return false
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

	n.update()

	//manual unmarshalling here required, to checkInputProgram for non existing variables in the tree
	return true
}

func (n *argnode) update() {
	//update all layer arguments in the tree recursively

}

func (p *parsetree) execute(cmdline *CommandLine) CLICODE {
	ok := CLICODE(CLI_FALSE)

	//fmt.Println(p)

	p._root.execute(cmdline)

	return ok
}

func (n *argnode) execute(cmdline *CommandLine) CLICODE {
	//fmt.Println(n)

	args := []*Argument{}

	for _, cmd := range n._sub {
		if cmd._arg != nil && cmd._arg.arg_type&OPTION > 0 {

			args = append(args, cmd._arg.copy())

			oparams := []*Argument{}
			oargs := []*Argument{}
			for _, p := range cmd._sub {
				if p._arg.arg_type&PARAMETER > 0 {
					oparams = append(oparams, p._arg.copy())
				}
			}
			if cmd._arg.run != nil {
				if code := cmd._arg.run(oparams, oargs, cmdline); code == CLI_SUCCESS {
					continue
				}
			}

		}

		if cmd._arg != nil && cmd._arg.arg_type&COMMAND > 0 {

			args = append(args, cmd._arg.copy())

			oparams := []*Argument{}
			oargs := []*Argument{}
			for _, p := range cmd._sub {
				if p._arg.arg_type&PARAMETER > 0 {
					oparams = append(oparams, p._arg.copy())
				}
				if p._arg.arg_type&OPTION > 0 {
					oargs = append(oargs, p._arg.copy())
				}
			}
			if cmd._arg.run != nil {
				if code := cmd._arg.run(oparams, oargs, cmdline); code == CLI_SUCCESS {
					continue
				}
			}
			if len(cmd._arg.runCommand) > 0 {

				if _, ok := os.OpenFile("./"+cmd._arg.runCommand, os.O_RDONLY, 0666); ok == nil {
					cmd := exec.Command("./" + cmd._arg.runCommand)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Run()
				} else {
					cmd := exec.Command(cmd._arg.runCommand)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Run()
				}
				return CLI_SUCCESS
			}

		}

	}

	if n._sub != nil && len(n._sub) > 0 {
		for _, cmd := range n._sub {
			if cmd._arg != nil && cmd._arg.arg_type&COMMAND > 0 {
				n._sub[0].execute(cmdline)
			}
		}

	}
	return CLI_SUCCESS
}
