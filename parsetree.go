package commandlinetoolkit

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

	_depth int
}

func newparsetree() *parsetree {

	t := &parsetree{_depth: 0,
		_settings: &settings{}}

	t._root = newNode(nil)

	return t
}

func newNode(arg *Argument) *argnode {
	node := &argnode{
		_arg:    arg.copy(),
		_sub:    []*argnode{},
		_parent: nil,
		_depth:  1,
	}

	if arg != nil && arg._sub != nil {
		for _, a := range arg._sub {
			_arg := newNode(a)
			node._sub = append(node._sub, _arg)
			if _arg._depth > node._depth {
				node._depth = _arg._depth + 1
			}
		}
		arg._sub = nil
	}

	return node
}

func (p *parsetree) Get(s string) *argnode {

	for i := 0; i < len(p._root._sub); i++ {
		if p._root._sub[i]._arg.lflag == s ||
			p._root._sub[i]._arg.data_type.data_flag == s {

			return p._root._sub[i]

		}
	}

	return &argnode{}
}

func (a *argnode) Value() any {
	if a._arg != nil {
		return a._arg.data_type.data
	}
	return nil
}

func (a *argnode) String() string {
	s := "{"

	if a._arg != nil {
		s += "\n \"type\" : " + argTypeToString(a._arg.arg_type) + "\n"
		s += " \"flag\" : " + a._arg.lflag + ",\n"
		s += " \"required\" :" + strconv.FormatBool(a._arg.required) + "\n"
		s += " \"depth\" : " + strconv.Itoa(a._depth) + "\n"
	}

	if len(a._sub) > 0 {
		for _, arg := range a._sub {
			s += "  " + arg.String()
		}
	}

	return s + "}"
}

func (a *argnode) HelpString() string {
	s := "\t"
	if a._arg == nil {
		s += "<nil>"
		return s
	}
	if a._arg.arg_type == OPTION {
		s += "--"
	}
	if len(a._arg.lflag) > 0 {

		s += string(COLOR_CYAN_I) + a._arg.lflag + string(COLOR_RESET) + " "

	}

	if len(s) < 22 {
		s += "\t"
	}
	if len(s) < 10 {
		s += " \t"
	}

	if len(s) > 24 {
		s = s[:24]
	}
	if len(a._sub) > 0 {
		for _, arg := range a._sub {
			s += arg._HelpString("\t"+string(COLOR_RESET)+string(COLOR_YELLOW_I), ": [ ", " ]"+string(COLOR_RESET))
		}
	}
	if len(s) < 36 {
		s += " \t"
	}
	if len(s) < 44 {
		s += " \t"
	}
	s += "\t" + string(COLOR_RESET) + ": "

	i := 0
	if a._arg.arg_type == OPTION {
		i += 1
		s += "O"
	}
	if a._arg.arg_type == COMMAND {
		i += 1
		s += "C"
	}
	if a._arg.arg_type == PARAMETER {
		i += 1
		s += "P"
	}
	if a._arg.arg_type == WILDCARD {
		i += 1
		s += "W"
	}
	if a._arg.arg_type == FLAG {
		i += 1
		s += "F"
	}
	for q := 0; q < 3-i; q++ {
		s += "-"
	}
	if a._arg.callback != nil {
		s += "C"
	} else {
		s += "-"
	}
	if a._arg.run != nil {
		s += "R"
	} else {
		s += "-"
	}

	s += "\t:"

	if len(a._arg.shelp) <= 0 {
		s += " " + a._arg.lhelp
	} else {
		s += " " + a._arg.shelp
	}
	/*
		if len(a._sub) > 0 {
			for _, arg := range a._sub {
				s += arg._HelpString("\n\t\t\t\t\t"+string(COLOR_CYAN_I), ": [ ", " ]"+string(COLOR_RESET))
			}
		}*/

	return s
}

func (a *argnode) _HelpString(s string, start string, end string) string {

	if a._arg.arg_type == PARAMETER {
		s += start + a._arg.data_type.data_flag
	}
	if a._sub != nil && len(a._sub) > 0 {
		for _, v := range a._sub {
			s += v._HelpString("", ", ", "")
		}
	}
	s += end
	return s
}

func (a *argnode) Get(s string) *argnode {

	for i := 0; i < len(a._sub); i++ {
		if a._sub[i]._arg.lflag == s ||
			a._sub[i]._arg.data_type.data_flag == s {

			return a._sub[i]

		}
	}
	return &argnode{}
}

func (a *argnode) Argument() *Argument {
	if a._arg == nil {
		return &Argument{}
	}
	ar := a._arg.copy()
	return ar
}

func (a *argnode) Next() *argnode {
	if len(a._sub) > 0 {
		return a._sub[0]
	}
	return &argnode{}
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

	p._root._depth = 0

	for _, arg := range args {

		_ok := p._root.addArgument(arg)
		if !_ok {
			ok = ok
		}

	}

	return ok
}

/*
**********************************************************************
Add a new Argument as a result of traversing the input tree from a given program
*/
func (n *argnode) addArgument(m map[string]interface{}) bool {

	/*	if !checkParseableArgFromProgramFile(m) {
		return
	}*/

	//create a node
	argN := &argnode{
		//if same layer depth = 1 otherwise increasing with depth
		_depth:  1,
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
	if len(n._sub) == 0 {
		n._depth = 1
	} else {

		depth := 1

		for _, arg := range n._sub {
			if arg._arg.arg_type == PARAMETER {
				n._arg.choices = append(n._arg.choices, arg._arg.lflag)
			}
			if arg._depth >= depth {
				depth = arg._depth + 1
			}

		}
		n._depth = depth
	}
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
			if cmd._arg.callback != nil {
				cmd._arg.callback(cmdline)
			}
			if cmd._arg.run != nil {
				if code := cmd._arg.run(oparams, oargs, cmdline); code == CLI_SUCCESS {
					continue
				}
			}

		}
	}
	for _, cmd := range n._sub {
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
			if len(cmd._sub) == 0 && cmd._arg.data_type.data != nil {
				if untokenizedSubArgs, ok := cmd._arg.data_type.data.([]string); ok {

					for _, argument := range untokenizedSubArgs {

						m := map[string]interface{}{}

						m[TYPEKEY] = OPTIONSTRING
						m[LONGFLAGKEY] = argument

						newargument, _ := NewArgument(m)

						newargument.data_type.data = argument

						oargs = append(oargs, newargument)

					}
				}
			}
			if cmd._arg.callback != nil {
				cmd._arg.callback(cmdline)
			}
			if cmd._arg.run != nil {
				if code := cmd._arg.run(oparams, oargs, cmdline); code == CLI_SUCCESS {
					continue
				}
			}
			if len(cmd._arg.runCommand) > 0 {
				if _, ok := os.OpenFile("./"+cmd._arg.runCommand, os.O_RDONLY, 0666); ok == nil {
					args := []string{}
					if cmd._arg.data_type.data != nil {
						if data, ok := cmd._arg.data_type.data.([]string); ok {
							args = data
						}
					}
					_cmd := exec.Command("./"+cmd._arg.runCommand, args...)
					_cmd.Stdin = os.Stdin
					_cmd.Stdout = os.Stdout
					_cmd.Run()
				} else {
					_cmd := exec.Command(cmd._arg.runCommand)
					_cmd.Stdin = os.Stdin
					_cmd.Stdout = os.Stdout
					_cmd.Run()
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

func (p *parsetree) String() string {

	s := "Parsetree:"

	s += p._root.String_("")

	return s

}

func (n *argnode) String_(str string) string {

	s := ""
	if n._arg != nil {
		data := fmt.Sprintln(n._arg.data_type.data)
		if len(data) > 20 {

			data = string(data[:20]) + "..."
		}
		s += "\n"
		s += str + "> type  : " + argTypeToString(n._arg.arg_type) + "\n"
		s += str + "> flag  : " + n._arg.lflag + "\n"
		s += str + "> help  : " + n._arg.lhelp + "\n"
		s += str + "> req.  : " + strconv.FormatBool(n._arg.required) + "\n"
		s += str + "> dtype : " + n._arg.data_type.data_flag + "\n"
		s += str + "> data  : " + data + ""
		s += str + "> depth : " + strconv.Itoa(n._depth) + "\n"
	}
	if n._sub != nil {
		if len(n._sub) > 0 {
			s += str + string(COLOR_CYAN) + "> arguments: " + string(COLOR_CYAN) + "\n"
		}
		for _, k := range n._sub {

			s += str + string(COLOR_CYAN) + "     --> " + string(COLOR_RESET) + k.String_(str+"         ") + ""
		}
	}
	return "" + s
}

func (p *parsetree) clone() *parsetree {

	np := &parsetree{
		_depth:    p._depth,
		_settings: p._settings.clone(),
		_root:     p._root.clone(),
	}

	return np

}

func (a *argnode) clone() *argnode {
	na := &argnode{
		_arg:    nil,
		_parent: nil,
		_sub:    []*argnode{},
		_depth:  a._depth,
	}
	if a._arg != nil {
		na._arg = a._arg.copy()
		if a._arg._sub != nil {
			//for _, sarg := range a._arg._sub {

			//newargnode := newNode(sarg)

			//na._sub = append(na._sub, newargnode)

			//}

		}
	}

	for i := 0; i < len(a._sub); i++ {
		na._sub = append(na._sub, a._sub[i].clone())
		na._sub[i]._parent = na
	}

	return na
}

func (a *argnode) appendArgument(arg *Argument) {

	a._sub = append(a._sub, newNode(arg))

}
