package commandlinetoolkit

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

func (p *parsetree) Get(s string) *argnode {
	
	for i := 0; i < len(p._root._sub); i++ {
		if p._root._sub[i]._arg.lflag == s ||
			p._root._sub[i]._arg.data_type.data_flag == s {
			
			return p._root._sub[i]
			
		}
	}
	
	return &argnode{}
}

func (a *argnode) String() string {
	s := "{"
	
	if a._arg != nil {
		s += "\n \"type\" : " + argTypeToString(a._arg.arg_type) + "\n"
		s += " \"flag\" : " + a._arg.lflag + ",\n"
		s += " \"required\" :" + strconv.FormatBool(a._arg.required) + "\n"
	}
	
	return s + "}"
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
	
	//fmt.Println("BUILD TOKEN")
	
	//fmt.Println(a._sub[0].String("  "), args[0])
	
	for a._sub != nil && len(a._sub) > 0 {
		if index >= len(a._sub) {
			
			break
			
		}
		arg = a._sub[index]
		if arg._arg.arg_type&OPTION > 0 {
			
			//fmt.Println("TEST: " + arg._arg.lflag)
			newarg, newargs, ok := tokenizeOption(arg, args)
			
			//fmt.Println("CREATED OPTION")
			//fmt.Println(newarg)
			
			if a != nil && newarg != nil {
				
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
				
				newarg._parent = node
				node._sub = append(node._sub, newarg)
				
			}
			
			if ok && len(newargs) == 0 {
				
				return node, true
			}
			
			args = newargs
		}
		
		if arg._arg.arg_type&PARAMETER > 0 {
			
			newarg, newargs, ok := tokenizeParameter(arg, args)
			
			if newarg != nil && ok {
				args = newargs
				
				node._sub = append(node._sub, newarg)
				newarg._parent = node
			}
			
		}
		
		index++
	}
	return node, true
}

func tokenizeOption(a *argnode, args []string) (*argnode, []string, bool) {
	ok := false
	if len(args) <= 0 {
		return nil, args, false
	}
	
	//fmt.Println("THE ARGUMENTS", args)
	///fmt.Println("THE ARGUMENT", a)
	
	if strings.Count(args[0], "=") == 1 {
		s := strings.Split(args[0], "=")
		for i := 1; i < len(args); i++ {
			s = append(s, args[i])
		}
		args = s
	}
	
	if len(a._arg.lflag) > 0 && strings.Index(args[0], FULLOPTIONPREFIX+a._arg.lflag) == 0 {
		ok = true
	}
	if len(a._arg.sflag) > 0 && strings.Index(args[0], SHORTOPTIONPREFIX+a._arg.sflag) == 0 {
		ok = true
	}
	
	//fmt.Println(args[0], ok, a._arg.lflag)
	if !ok {
		return nil, args, false
	}
	
	index := 1
	
	newargnode := newNode(a._arg.copy())
	
	//fmt.Println(args)
	/*
		if len(args) == 0 {
			//split on =
			data, ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, "")
	
			if ok {
				newargnode._arg.data_type.data = data
			} else {
	
			}
		}*/
	
	if index < len(args) {
		if len(a._arg.data_type.data_flag) > 0 {
			//fmt.Println("TEadsgasdgasdgyST", a._arg.data_type.data_flag, len(args))
			//fmt.Println(len(args))
			
			data, _ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index])
			
			//fmt.Println(_ok)
			if _ok {
				ok = _ok
				newargnode._arg.data_type.data = data
				//check next node for parameter
				
				if len(a._sub) == 1 {
					na, ok := a.tokenizeArg(args[index:])
					
					//fmt.Println("WAT THE FUCK")
					//fmt.Println(na._sub[0].String(""))
					
					na._parent = newargnode
					newargnode._sub = []*argnode{}
					
					newsubnode := na._sub[0].clone()
					newsubnode._arg.data_type.data = na._sub[0]._arg.data_type.data
					//fmt.Println("WAT THE FUCK")
					//fmt.Println(newsubnode.String(""))
					//fmt.Println("WAT THE FUCK")
					newargnode._sub = append(newargnode._sub, newsubnode)
					
					if !ok {
						return newargnode, args[index:], false
					}
				}
				
				index = 2
			} else {
			
			}
			
		}
		if a._sub == nil || len(a._sub) == 0 {
		}
	} else {
		
		if len(a._sub) == 1 {
			// fmt.Println("kjabesgplajbsrgpojlbs")
			//fmt.Println(a._sub[0]._arg)
			if a._sub[0]._arg.required {
				return newargnode, args[index:], false
			}
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

func tokenizeParameter(a *argnode, args []string) (*argnode, []string, bool) {
	ok := false
	//fmt.Println("TOKENIZING PARANETER")
	//fmt.Println(a.String(" "))
	
	if len(args) <= 0 {
		
		if a._arg.required {
			
			newDebugHandler().printError("Required Argument not provided:" + a._arg.data_type.data_flag)
			if len(args) <= 0 {
				return nil, args, false
			}
		}
		
		return nil, args, false
	}
	
	//fmt.Println(a.String("   "))
	
	newargnode := newNode(a._arg.copy())
	if a._arg.required {
		
		if len(args) <= 0 {
			return nil, args, false
		}
	}
	
	//fmt.Println(newargnode.String("   "))
	index := 0
	
	if data, _ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index]); _ok {
		
		newargnode._arg.data_type.data = data
		ok = _ok
		
		if _ok {
			ok = _ok
			newargnode._arg.data_type.data = data
			//check next node for parameter
			//fmt.Println("WEITERE ARGUMENTE")
			if len(a._sub) == 1 {
				
				na, _ok := a.tokenizeArg(args[index+1:])
				if _ok {
					ok = ok
				}
				
				if len(na._sub) > 0 {
					
					na._parent = newargnode
					newargnode._sub = []*argnode{}
					
					newsubnode := na._sub[0].clone()
					newsubnode._arg.data_type.data = na._sub[0]._arg.data_type.data
					//fmt.Println("WAT THE FUCK")
					//fmt.Println(newsubnode.String(""))
					//fmt.Println("WAT THE FUCK")
					newargnode._sub = append(newargnode._sub, newsubnode)
					
				}
				
				if !ok {
					return newargnode, args[index:], false
					
				}
				
				index++
			} else {
			
			}
			
		}
		
		if len(a._sub) > 0 {
		
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

func (p *parsetree) String() string {
	
	s := "Parsetree:"
	
	s += p._root.String_("")
	
	return s
	
}

func (n *argnode) String_(str string) string {
	
	s := ""
	if n._arg != nil {
		
		data := fmt.Sprintln(n._arg.data_type.data)
		
		if len(data) > 25 {
			
			data = string(data[:25]) + "...\n"
		}
		
		s += "" + "> type: " + argTypeToString(n._arg.arg_type) + "\n"
		s += str + "> flag: " + n._arg.lflag + "\n"
		s += str + "> help: " + n._arg.lhelp + "\n"
		s += str + "> req.: " + strconv.FormatBool(n._arg.required) + "\n"
		s += str + "> dtype: " + n._arg.data_type.data_flag + "\n"
		s += str + "> data: " + data + ""
	}
	
	if n._sub != nil {
		if len(n._sub) > 0 {
			s += str + "> arguments: " + "\n"
		}
		for _, k := range n._sub {
			
			s += str + string(COLOR_CYAN) + "--> " + string(COLOR_RESET) + k.String_(str+"    ") + ""
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
	}
	
	if a._arg != nil {
		na._arg = a._arg.copy()
		
	}
	if len(a._sub) > 0 {
		
		for i := 0; i < len(a._sub); i++ {
			
			na._sub = append(na._sub, a._sub[i].clone())
			na._sub[i]._parent = na
		}
		
	}
	return na
}
