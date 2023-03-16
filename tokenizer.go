package commandlinetoolkit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const verbose = false

/*
**********************************************************************
tokenize the input from arguments, given by the os or the shell
build a tree recurively
*/
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

	/*
		New definition from parsing the tree,
		args: the arguments left, after parsing the input tree
		newArg being the root of the entire tree
		ok: the current subtree was successfully parsed
		skip: the current
	*/
	args, newArg, ok, skip := r.tokenizeArg(args)

	if !ok && skip {
		//fmt.Println(newArg)
	} else {
		code = CLI_SUCCESS
	}
	np._root = newArg

	//fmt.Println("NEW Parsetree", np)
	return np, code
}

/*
**********************************************************************
Tokenize a single argument

returns args, thesubtreenode, tokenizing was ok, skip to nextnode in the same layer, if ok, the command will execute either load now cli or will run other binary with set options
for trailing options or commands, tree can tokenize recurisvely again
*/
func (a *argnode) tokenizeArg(args []string) ([]string, *argnode, bool, bool) {

	newArgumentNode := newNode(nil)

	newArgumentNode._depth = 1

	if verbose {
		fmt.Println(string(COLOR_YELLOW_I) + ">>>>>>>>:::::::::: INPUT ::::::::::<<<<<<<<" + string(COLOR_RESET))
		fmt.Println(args)
		fmt.Println(a)
	}

	//debugging

	if a._sub == nil || len(a._sub) == 0 {
		return args, newArgumentNode, false, true
	}

	ok := false
	skip := false
	argindex := 0
	subArgument := &argnode{}
	scan := 0

	for a._sub != nil && len(a._sub) > 0 {

		//time.Sleep(1 * time.Second)
		if verbose {
			bufio.NewScanner(os.Stdin).Scan()
		}
		if argindex >= len(a._sub) {
			//cycle through the subarguments, always check option first, then command, then paramter
			if len(args) > 0 && scan < 1 {
				argindex = 0
				scan = 1
			} else {
				break
			}
		}
		subArgument = a._sub[argindex]
		argindex++

		if verbose {
			fmt.Println(string(COLOR_PINK_I) + ":::::::::: TOKENIZE SUB ARGUMENT::::::::::" + string(COLOR_RESET))
			fmt.Println(subArgument)
		}

		/*
			Parse now an option
		*/

		if subArgument._arg.arg_type&OPTION > 0 {

			if verbose {
				fmt.Println(string(COLOR_PINK_I) + ":::::::::: BEFORE Tokenize Option ::::::::::" + string(COLOR_RESET))
			}

			newarg, newargs, _ok, _skip := tokenizeOption(subArgument, args)

			if verbose {
				fmt.Println(string(COLOR_PINK_I) + ":::::::::: AFTER Tokenize Option ::::::::::" + string(COLOR_RESET))
				fmt.Println(newarg, newargs, _ok, _skip)

			}

			args = newargs
			if newarg != nil {

				newarg._parent = newArgumentNode
				newArgumentNode._sub = append(newArgumentNode._sub, newarg)
				newArgumentNode._depth = newarg._depth + 1
			}

			if _skip {
				args = newargs
				scan = 0
			}

			if _ok {
				args = newargs
				scan = 0
				ok = _ok
			}

			if ok && len(newargs) == 0 {
				return newargs, newArgumentNode, ok, skip
			}

		}

		if subArgument._arg.arg_type&COMMAND > 0 {
			if verbose {
				fmt.Println(string(COLOR_PINK_I) + ":::::::::: BEFORE Tokenize Command ::::::::::" + string(COLOR_RESET))
			}

			newarg, newargs, _ok, _skip := tokenizeCommand(subArgument, args)

			if verbose {
				fmt.Println(string(COLOR_PINK_I) + ":::::::::: AFTER Tokenize Command ::::::::::" + string(COLOR_RESET))
				fmt.Println(newarg, newargs, _ok, _skip)
			}

			if a != nil && newarg != nil && _ok {
				args = newargs
				scan = 0
				newarg._parent = newArgumentNode
				newArgumentNode._sub = append(newArgumentNode._sub, newarg)
				newArgumentNode._depth = newarg._depth + 1
				//here we have to jump into the corresponding layer for the command
				a = subArgument

			}

			if _skip {
				args = newargs
			}

			if _ok && len(newargs) == 0 {
				return newargs, newArgumentNode, ok, skip
			}

		}

		if subArgument._arg.arg_type&PARAMETER > 0 {
			if verbose {
				fmt.Println(string(COLOR_PINK_I) + ":::::::::: BEFORE Tokenize Parameter ::::::::::" + string(COLOR_RESET))
			}
			newarg, newargs, _ok, _skip := tokenizeParameter(subArgument, args)

			if verbose {
				fmt.Println(string(COLOR_PINK_I) + ":::::::::: AFTER Tokenize Parameter ::::::::::" + string(COLOR_RESET))
				fmt.Println(newarg, newargs, _ok, _skip)
			}
			args = newargs
			if newarg != nil {
				newArgumentNode._sub = append(newArgumentNode._sub, newarg)
				newarg._parent = newArgumentNode
				newArgumentNode._depth = newarg._depth + 1
			}

			if _ok {
				skip = _skip
				//scan = 0
			}
			if _skip {
				args = newargs
			}

			if _skip {
				return args, newArgumentNode, _ok, skip
			}
			if !_ok {
				return args, newArgumentNode, _ok, true
			}

		}

	}

	if verbose {
		fmt.Println(string(COLOR_YELLOW_I) + ">>>>>>>>>>:::::::::: END TOKENIZE ::::::::::<<<<<<<<<<" + string(COLOR_RESET))
		fmt.Println(args, newArgumentNode._sub, ok, skip)
	}

	return args, newArgumentNode, ok, skip
}

/*
Peeks into the single option, if it can parse the option, given a
single datapoint afterwards seperated by a '=' or ' ', then it will be able to
If multiple parameters are required, then it does not accept the = sign, only with
current depth of the subtree is <= 2
*/
func tokenizeOption(a *argnode, args []string) (*argnode, []string, bool, bool) {
	ok := false
	skip := false

	if len(args) <= 0 {
		return nil, args, ok, skip
	}

	//does it contain the option or the shortflag option
	if len(a._arg.lflag) > 0 && strings.Index(args[0], FULLOPTIONPREFIX+a._arg.lflag) == 0 {
		ok = true
	}
	if len(a._arg.sflag) > 0 && strings.Index(args[0], SHORTOPTIONPREFIX+a._arg.sflag) == 0 {
		ok = true
	}

	if !ok {
		if verbose {
			fmt.Println("Not parsed...")
		}
		return nil, args, ok, skip
	}

	newOptionArgument := newNode(a._arg.copy())

	//check wether a single = is present
	if strings.Count(args[0], "=") == 1 {
		if len(a._sub) >= 1 && a._depth == 2 {
			if verbose {
				fmt.Println(string(COLOR_YELLOW_I) + ":::::::::: OPTION SPLIT ::::::::::" + string(COLOR_RESET))
			}
			s := strings.Split(args[0], "=")
			for i := 1; i < len(args); i++ {
				s = append(s, args[i])
			}
			args = s

		} else {
			if verbose {
				fmt.Println(string(COLOR_YELLOW_I) + ":::::::::: OPTION Doesnt allow SPLIT ::::::::::" + string(COLOR_RESET))
			}
			ok = false
			skip = false

			return newOptionArgument, args, ok, skip
		}
	}

	//here we remove the first option, if there was a splittable part, it was already split before
	args = args[1:]
	if verbose {
		fmt.Println(string(COLOR_GREEN_I) + ":::::::::: OPTION WAS ACCEPTED::::::::::" + string(COLOR_RESET))
		fmt.Println(newOptionArgument, args)

		//if len(args) > 0 {

		//check next node for trailing

		fmt.Println(string(COLOR_RED_I) + ":::::::::: BEFORE Tokenize in OPTION ::::::::::" + string(COLOR_RESET))
		fmt.Println(a, args)
	}
	newargs, newarg, _ok, _skip := a.tokenizeArg(args)

	args = newargs

	if len(newarg._sub) == 1 && newarg._sub[0]._arg != nil {
		args = newargs
		newarg._sub[0]._parent = newOptionArgument

		newOptionArgument._sub = append(newOptionArgument._sub, newarg._sub[0])
		newOptionArgument._depth = newOptionArgument._depth + 1
		if newOptionArgument._arg.data_type.data_flag == newarg._sub[0]._arg.data_type.data_flag {
			newOptionArgument._arg.data_type.data = newarg._sub[0]._arg.data_type.data
		}
	}
	if verbose {
		fmt.Println(":::::::::: AFTER Tokenize in OPTION ::::::::::")
		fmt.Println(newarg, newargs, _ok, _skip)
	}
	//could not properly parse the subtree, here we exit with ok=false
	if !_ok {
		return newOptionArgument, args, _ok, _skip
	}

	//if we should skip, previous iterations in tokenizeArg should have removed the elements in args
	//we return with _ok and true
	if _skip {
		return newOptionArgument, newargs, _ok, _skip
	}

	return newOptionArgument, args, ok, skip
}

func tokenizeCommand(a *argnode, args []string) (*argnode, []string, bool, bool) {
	ok := false
	skip := false

	if len(args) <= 0 {
		return nil, args, false, skip
	}

	//the datapoint has to be exactly the same, as required by the lflag of the command
	//comamnds dont have shortflags, and if defined by the program, they are not used
	if len(a._arg.lflag) > 0 && strings.Compare(args[0], a._arg.lflag) == 0 {
		ok = true
	}
	if !ok {
		return nil, args, ok, skip
	}

	index := 1

	newargnode := newNode(a._arg.copy())
	//fmt.Println(args)
	if verbose {
		fmt.Println(string(COLOR_GREEN_I) + ":::::::::: COMMAND ACCEPTED ::::::::::" + string(COLOR_RESET))
	}
	args = args[index:]

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
	return newargnode, newargs, ok, skip
}

func tokenizeParameter(a *argnode, args []string) (*argnode, []string, bool, bool) {
	ok := false
	skip := false
	index := 0
	if verbose {
		fmt.Println("::::::::: TOKENIZE PARAMETER ::::::::::")
		fmt.Println(a)
		fmt.Println(args)
	}
	if len(args) <= 0 {

		if a._arg.required {

			newDebugHandler().printError("Required Argument not provided:" + a._arg.data_type.data_flag + "\n")

			if len(args) <= 0 {
				return nil, args, false, skip
			}
		}
		skip = true
		return nil, args, false, skip
	}
	//have to take care of paramters that potentially LOOK like options
	/*
		if strings.Index(args[0], "-") == 0 || strings.Index(args[0], "--") == 0 {
			return nil, args, false
		}*/

	//fmt.Println(a.String("   "))

	//a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index])

	newParameterArgument := &argnode{}
	/*
		if a._arg.required {

			if len(args) <= 0 {
				return nil, args, false, skip
			}
		}*/

	//fmt.Println(newargnode.String("   "))

	if data, _ok := a._arg.data_type.dtype_custom_callback(a._arg.data_type, args[index]); _ok {

		newParameterArgument = newNode(a._arg.copy())
		/*
			The Parameter was accepted by the specific callback function provided by the framework or the user

		*/
		if verbose {
			fmt.Println(string(COLOR_GREEN_I) + ":::::::::: PARAMETER ACCEPTED ::::::::::" + string(COLOR_RESET))
			fmt.Println(a)
		}
		args = args[1:]
		newParameterArgument._arg.data_type.data = data
		ok = _ok

		//check next node for parameter
		//fmt.Println("WEITERE ARGUMENTE")
		if len(a._sub) >= 0 {
			if verbose {
				fmt.Println(string(COLOR_RED_I) + ":::::::::: BEFORE Tokenize in PARAMETER ::::::::::" + string(COLOR_RESET))
			}
			newArgs, newSubArgumentNode, _ok, _skip := a.tokenizeArg(args)
			if verbose {
				fmt.Println(string(COLOR_RED_I) + ":::::::::: AFTER Tokenize in PARAMETER ::::::::::" + string(COLOR_RESET))
			}

			//fmt.Println(newArgs, newSubArgumentNode._sub, _ok, _skip)

			if _ok {
				ok = _ok

			}
			if _skip {

				args = newArgs
				skip = _skip
			}
			if newSubArgumentNode != nil && len(newSubArgumentNode._sub) == 1 && newSubArgumentNode._sub[0]._arg != nil {

				newParameterArgument._sub = []*argnode{}

				newsubnode := newSubArgumentNode._sub[0].clone()

				newsubnode._parent = newParameterArgument

				newsubnode._arg.data_type.data = newSubArgumentNode._sub[0]._arg.data_type.data

				if newsubnode._parent._arg.arg_type == OPTION {

					newsubnode._parent._arg.data_type.data = newsubnode._arg.data_type.data

				}
				newParameterArgument._sub = append(newParameterArgument._sub, newsubnode)

				newParameterArgument._depth = newsubnode._depth + 1
				args = newArgs

			}

			if !_ok || newArgs == nil {
				return newParameterArgument, args, ok, skip
			}

		} else {

			return newParameterArgument, args, ok, skip

		}

	} else {
		if verbose {
			fmt.Println(string(COLOR_BLUE) + ":::::::::: PARAMETER NOT ACCEPTED ::::::::::" + string(COLOR_RESET))
			fmt.Println(a)
		}
		if a._arg.required {
			ok = false
		} else {
			ok = true
		}
		skip = true

		return nil, args, ok, skip
	}
	return newParameterArgument, args, ok, skip
}
