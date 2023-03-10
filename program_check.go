package commandlinetoolkit

import (
	"strings"
)

/*
********************************************************************

Do a first tree sweep check via typecast checking, the parsetree has to apply the given map then to the different arguments int he tree seperately
if coliding parameters are present, they are excluded,savefiles will be written
*/
func (p *program) checkInputProgram() {

	//checkInputProgram for required defaults, that are required to run the cli, if the contents are not present parsing will not hapen
	allArgs := p._program._theProgramJsonMap[ARGUMENTSKEY]

	//we expect here either a well formed json will default immuteable (replaced if removed) arguments
	//on simple read error we create a new json with the given commands
	//if successfull todo is checkInputProgram if the structure is well formed
	//recursive type binding for []interface{} and []map[string]interface{} trees nd tree
	pargs := programArgs{}

	programDepth := 0

	if allArgs == nil {
		allArgs = []map[string]interface{}{}
		//can already exit the subroutine and create a new default json file
	}

	ok := false
	checkedPrgm := []map[string]interface{}{}
	// inspiration from https://stackoverflow.com/questions/29366038/looping-iterate-over-the-second-level-nested-json-in-go-lang
	switch allArgs.(type) {
	case []map[string]interface{}:
		{
			//checkInputProgram recursively
			checkedPrgm, ok, programDepth = checkProgramTopLevel(allArgs.([]map[string]interface{}))
		}
	case []interface{}:
		{

			//checkInputProgram recursively
			checkedPrgm, ok, programDepth = checkProgramInterfaceArray(allArgs.([]interface{}))

			_, pargs = checkProgramArgs(checkedPrgm)

			ok = pargs.success()

		}
	default:
		{
			ok = false
		}

	}
	p._programDepth = programDepth

	if !ok {

		p._debugHandler.printError("-->Program: Error: Could not read Arguments, writing defaults")

		saveErrorFile(p._programFile, p._program._theProgramJson)

		p._program._theProgramJsonMap = DefaultTemplate()

		p.writeJsonProgram()
	} else {

		if !pargs.hasHelpOption {
			checkedPrgm = append(checkedPrgm, defaultHelpOption())
		}
		if !pargs.hasVerbosityOption {
			checkedPrgm = append(checkedPrgm, defaultVerbosityOption())
		}

		p._program._theProgramJsonMap["arguments"] = checkedPrgm

		p.writeJsonProgram()
	}
}

func checkProgramArgs(args []map[string]interface{}) ([]map[string]interface{}, programArgs) {

	m := []map[string]interface{}{}

	pargs := programArgs{}

	for _, arg := range args {

		if arg[TYPEKEY] == OPTIONSTRING {

			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultInteractiveOption) == 0 {
				//hasInteractiveOption = true
				//just in case replacement
				pargs.hasInteractiveOption = true
				if !theInteractiveOption[MUTEABLEKEY].(bool) {
					m = append(m, theInteractiveOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultLoggingOption) == 0 {

				pargs.hasLoggingOption = true
				if !theLoggingOption[MUTEABLEKEY].(bool) {
					m = append(m, theLoggingOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultHistoryOption) == 0 {

				pargs.hasHistoryOption = true
				if !theHistoryOption[MUTEABLEKEY].(bool) {
					m = append(m, theHistoryOption)
				}

			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultHelpOption) == 0 {

				pargs.hasHelpOption = true
				if !theHelpOption[MUTEABLEKEY].(bool) {
					m = append(m, theHelpOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _ddefaultHistoryFileOption) == 0 {

				pargs.hasHistoryFileOption = true
				if !theHistoryFileOption[MUTEABLEKEY].(bool) {
					m = append(m, theHistoryFileOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultConfigurationFileOption) == 0 {

				//just replace the option anyways, so we dont have to do recursive or deep layer checks
				pargs.hasConfigurationFileOption = true
				if !theConfigFileOption[MUTEABLEKEY].(bool) {
					m = append(m, theConfigFileOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultVerbosityOption) == 0 {

				pargs.hasVerbosityOption = true
				if !theVerbosityOption[MUTEABLEKEY].(bool) {
					m = append(m, theVerbosityOption)
				}
			}
			checkOption(arg[LONGFLAGKEY].(string), &arg, &map[string]interface{}{})
		}
	}

	return m, pargs
}

func checkForDefaultProgramArgument(str string) bool {
	switch str {
	case _defaultHelpOption:
		{
			return true
		}
	case _defaultHistoryOption:
		{
			return true
		}
	case _defaultInteractiveOption:
		{
			return true
		}
	case _defaultLoggingOption:
		{
			return true
		}
	case _defaultConfigurationFileOption:
		{
			return true
		}
	case _defaultVerbosityOption:
		{
			return true
		}
	case _ddefaultHistoryFileOption:
		{
			return true
		}
	}
	return false
}

func checkProgramTopLevel(m []map[string]interface{}) ([]map[string]interface{}, bool, int) {

	args, pargs := checkProgramArgs(m)

	if !pargs.success() {
		return nil, false, 0
	}
	ok := true
	depth := 0
	for _, value := range m {

		if checkForDefaultProgramArgument(value[LONGFLAGKEY].(string)) {
			continue
		}

		newmap, _ok, _depth := checkProgramMap(value)
		if !_ok {
			ok = _ok
			continue
		}

		args = append(args, newmap)
		if _depth > depth {
			depth = _depth
		}

	}

	return args, ok, depth
}

func checkProgramMap(inputmap map[string]interface{}) (map[string]interface{}, bool, int) {

	resultmap := map[string]interface{}{}

	theType, ok := inputmap[TYPEKEY].(string)

	depth := 0

	if !ok {
		return resultmap, false, 0
	}

	theTypes := []string{theType}

	theTypes = splitBySeperator(theTypes)

	if len(theTypes) == 1 {
		_depth := 0
		_, _ok, depth := checkParameter(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		_, _ok, _depth = checkOption(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		if _depth > depth {
			depth = _depth
		}
		_, _ok = checkWildcard(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		_, _ok = checkFlag(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		_, _ok, _depth = checkCommand(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		if _depth > depth {
			depth = _depth
		}

		return resultmap, ok, depth + 1
	}
	if len(theTypes) > 1 {
		_, _ok, _depth := checkMultiArgument(theTypes, &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		if _depth > depth {
			depth = _depth
		}
	}

	//if has datatype, append datatype as sub element, if has single agrument proceeding as datatype, add datatype

	return resultmap, ok, depth + 1
}

func checkMultiArgument(types []string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool, int) {

	ok := true

	argtype := ArgumentType(0)

	depth := 0

	for _, t := range types {

		if strings.Compare(t, OPTIONSTRING) == 0 {

			argtype |= OPTION
			continue
		}

		if strings.Compare(t, FLAGSTRING) == 0 {

			argtype |= FLAG
			continue
		}
		if strings.Compare(t, COMMANDSTRING) == 0 {

			argtype |= COMMAND
			continue
		}
		if strings.Compare(t, PARAMETERSTRING) == 0 {

			argtype |= PARAMETER
			continue
		}
		if strings.Compare(t, WILDCARDSTRING) == 0 {

			argtype |= WILDCARD
			continue
		}

	}

	if !validateArgType(argtype) {
		return (*resultmap), false, 0
	}

	if (argtype & PARAMETER) > 0 {
		checkParameter(PARAMETERSTRING, resultmap, inputmap)
	}
	if (argtype & OPTION) > 0 {
		_, _, _depth := checkOption(OPTIONSTRING, resultmap, inputmap)
		if _depth > depth {
			depth = _depth
		}
	}
	if (argtype & WILDCARD) > 0 {
		checkWildcard(WILDCARDSTRING, resultmap, inputmap)
	}
	if (argtype & FLAG) > 0 {
		checkFlag(FLAGSTRING, resultmap, inputmap)
	}
	if (argtype & COMMAND) > 0 {
		_, _, _depth := checkCommand(COMMANDSTRING, resultmap, inputmap)
		if _depth > depth {
			depth = _depth
		}
	}

	(*resultmap)[TYPEKEY] = ""

	str := []string{}

	if (argtype & PARAMETER) > 0 {
		str = append(str, PARAMETERSTRING)

	}
	if (argtype & OPTION) > 0 {
		str = append(str, OPTIONSTRING)

	}
	if (argtype & WILDCARD) > 0 {
		str = append(str, WILDCARDSTRING)

	}
	if (argtype & FLAG) > 0 {
		str = append(str, FLAGSTRING)

	}
	if (argtype & COMMAND) > 0 {
		str = append(str, COMMANDSTRING)

	}

	for i, s := range str {
		if i == 0 {
			(*resultmap)[TYPEKEY] = s
		} else {
			(*resultmap)[TYPEKEY] = (*resultmap)[TYPEKEY].(string) + " | " + s
		}
	}

	return (*resultmap), ok, depth
}

func checkParameter(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool, int) {

	depth := 0
	if strings.Compare(theType, PARAMETERSTRING) == 0 {
		//this is a parameter

		copyProgramArgument(resultmap, inputmap)

		if (*inputmap)[ARGUMENTSKEY] != nil {

			if i, ok := ((*inputmap)[ARGUMENTSKEY].([]interface{})); ok {
				m, ok, q := checkProgramInterfaceArray(i)
				depth = q
				if ok {
					(*resultmap)[ARGUMENTSKEY] = m
				}
			}

			return (*resultmap), false, depth + 1
		}
	}
	return (*resultmap), true, depth + 1
}

func checkOption(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool, int) {

	ok := true
	if strings.Compare(theType, OPTIONSTRING) != 0 {
		return (*resultmap), true, 0
	}
	//this is an option
	if str, _ok := (*inputmap)[LONGFLAGKEY].(string); !_ok || len(str) <= 0 {
		return (*resultmap), false, 0
	}

	copyProgramArgument(resultmap, inputmap)

	//we have no trailing arguments, we can eject here and have a depth of 1
	if (*inputmap)[ARGUMENTSKEY] == nil {

		if dtype, hasDtype := (*resultmap)[DATATYPEKEY].(string); hasDtype {

			//create Parameter argument

			newArg := map[string]interface{}{}
			newArg[TYPEKEY] = PARAMETERSTRING
			newArg[DATATYPEKEY] = dtype

			(*resultmap)[ARGUMENTSKEY] = []map[string]interface{}{}

			(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newArg)

		}

		return (*resultmap), true, 1
	}
	(*resultmap)[ARGUMENTSKEY] = []map[string]interface{}{}

	depth := 0
	trailingArguments := 0

	switch (*inputmap)[ARGUMENTSKEY].(type) {
	case []interface{}:
		{

			for _, v := range (*inputmap)[ARGUMENTSKEY].([]interface{}) {
				arg, _ok := v.(map[string]interface{})
				if !_ok {
					ok = _ok
					continue
				}
				newarg, _ok, _depth := checkProgramMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
				trailingArguments++
				if _depth > depth {
					_depth = depth
				}

			}
			break
		}
	case []map[string]interface{}:
		{
			for _, arg := range (*inputmap)[ARGUMENTSKEY].([]map[string]interface{}) {
				newarg, _ok, _depth := checkProgramMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
				if !_ok {
					ok = _ok
				}
				trailingArguments++
				if _depth > depth {
					_depth = depth
				}
			}
			break
		}
	}

	if dtype, hasDtype := (*resultmap)[DATATYPEKEY].(string); hasDtype {

		if trailingArguments == 0 {

			//create Parameter argument

			newArg := map[string]interface{}{}
			newArg[TYPEKEY] = PARAMETERSTRING
			newArg[DATATYPEKEY] = dtype

			(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newArg)
		}
	}

	if _, hasDtype := (*resultmap)[DATATYPEKEY].(string); !hasDtype {

		if trailingArguments == 1 && (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][TYPEKEY] == PARAMETERSTRING {

			//create Parameter argument

			(*resultmap)[DATATYPEKEY] = (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][DATATYPEKEY]
		}
	}
	if dtype, hasDtype := (*resultmap)[DATATYPEKEY].(string); hasDtype {
		if trailingArguments == 1 && (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][TYPEKEY] == PARAMETERSTRING {
			if (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][DATATYPEKEY] != dtype {
				dtype = (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][DATATYPEKEY].(string)
				(*resultmap)[DATATYPEKEY] = dtype
			}
		}
	}
	if dtype, hasDtype := (*resultmap)[DATATYPEKEY].(string); hasDtype {
		if trailingArguments == 1 && (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][TYPEKEY] == PARAMETERSTRING {
			if (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][DATATYPEKEY] != dtype {
				dtype = (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][DATATYPEKEY].(string)
				(*resultmap)[DATATYPEKEY] = dtype
			}
		}
	}

	return *resultmap, ok, depth + 1
}

func checkWildcard(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	if strings.Compare(theType, WILDCARDSTRING) == 0 {
		if str, _ok := (*inputmap)[LONGFLAGKEY].(string); _ok && len(str) > 0 {
			copyProgramArgument(resultmap, inputmap)
		}
		return (*resultmap), true
	}
	return (*resultmap), true
}

func checkFlag(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	if strings.Compare(theType, FLAGSTRING) == 0 {
		if str, _ok := (*inputmap)[LONGFLAGKEY].(string); _ok && len(str) > 0 {
			copyProgramArgument(resultmap, inputmap)
		}
		return (*resultmap), true
	}
	return (*resultmap), true
}

func checkCommand(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool, int) {

	ok := true
	if strings.Compare(theType, COMMANDSTRING) != 0 {
		return (*resultmap), true, 0
	}
	//this is an option
	if str, _ok := (*inputmap)[LONGFLAGKEY].(string); !_ok || len(str) <= 0 {
		return (*resultmap), false, 0
	}

	if str, _ok := (*inputmap)[RUNKEY].(string); !_ok || len(str) <= 0 {
		return (*resultmap), false, 0
	}

	copyProgramArgument(resultmap, inputmap)

	if (*inputmap)[ARGUMENTSKEY] == nil {
		if dtype, hasDtype := (*resultmap)[DATATYPEKEY].(string); hasDtype {

			//create Parameter argument

			newArg := map[string]interface{}{}
			newArg[TYPEKEY] = PARAMETERSTRING
			newArg[DATATYPEKEY] = dtype

			(*resultmap)[ARGUMENTSKEY] = []map[string]interface{}{}

			(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newArg)

		}

		return (*resultmap), true, 1
	}
	(*resultmap)[ARGUMENTSKEY] = []map[string]interface{}{}

	depth := 0
	trailingArguments := 0

	switch (*inputmap)[ARGUMENTSKEY].(type) {
	case []interface{}:
		{

			for _, v := range (*inputmap)[ARGUMENTSKEY].([]interface{}) {
				arg, _ok := v.(map[string]interface{})
				if !_ok {
					ok = _ok
					continue
				}
				newarg, _ok, _depth := checkProgramMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
				trailingArguments++
				if _depth > depth {
					_depth = depth
				}

			}
			break
		}
	case []map[string]interface{}:
		{
			for _, arg := range (*inputmap)[ARGUMENTSKEY].([]map[string]interface{}) {
				newarg, _ok, _depth := checkProgramMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
				if !_ok {
					ok = _ok
				}
				trailingArguments++
				if _depth > depth {
					_depth = depth
				}
			}
			break
		}
	}

	if dtype, hasDtype := (*resultmap)[DATATYPEKEY].(string); hasDtype {

		if trailingArguments == 0 {

			//create Parameter argument

			newArg := map[string]interface{}{}
			newArg[TYPEKEY] = PARAMETERSTRING
			newArg[DATATYPEKEY] = dtype

			(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newArg)
		}
	}

	if _, hasDtype := (*resultmap)[DATATYPEKEY].(string); !hasDtype {

		if trailingArguments == 1 && (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][TYPEKEY] == PARAMETERSTRING {

			//create Parameter argument

			(*resultmap)[DATATYPEKEY] = (*resultmap)[ARGUMENTSKEY].([]map[string]interface{})[0][DATATYPEKEY]
		}
	}

	return *resultmap, ok, depth
}

func checkProgramInterfaceArray(i []interface{}) ([]map[string]interface{}, bool, int) {

	m := []map[string]interface{}{}
	var ok bool = true

	programDepth := 0

	for _, v := range i {

		if arg, __ok := (v.(map[string]interface{})); __ok {

			_arg, _ok, _programDepth := checkProgramMap(arg)
			if !_ok {
				ok = _ok
				continue
			}
			if _programDepth > programDepth {
				programDepth = _programDepth
			}
			m = append(m, _arg)
		} else {
			ok = __ok
		}
	}
	return m, ok, programDepth
}
