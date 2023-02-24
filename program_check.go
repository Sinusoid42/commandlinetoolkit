package commandlinetoolkit

import "strings"

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
			checkedPrgm, ok = checkProgramTopLevel(allArgs.([]map[string]interface{}))
		}
	case []interface{}:
		{
			//checkInputProgram recursively
			checkedPrgm, ok = checkProgramInterfaceArray(allArgs.([]interface{}))

		}
	default:
		{
			ok = false
		}

	}

	if !ok {

		p._debugHandler.printError("-->Program: Error: Could not read Arguments, writing defaults")

		saveErrorFile(p._programFile, p._program._theProgramJson)

		p._program._theProgramJsonMap = DefaultTemplate()

		p.writeJsonProgram()
	} else {

		p._program._theProgramJsonMap["arguments"] = checkedPrgm

		p.writeJsonProgram()
	}
}

func checkProgramArgs(args []map[string]interface{}) ([]map[string]interface{}, bool) {

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

				//hasLoggingOption = true
				if !theLoggingOption[MUTEABLEKEY].(bool) {
					m = append(m, theLoggingOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultHistoryOption) == 0 {

				//hasHistoryOption = true
				if !theHistoryOption[MUTEABLEKEY].(bool) {
					m = append(m, theHistoryOption)
				}

			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultHelpOption) == 0 {

				//hasHelpOption = true
				if !theHelpOption[MUTEABLEKEY].(bool) {
					m = append(m, theHelpOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _ddefaultHistoryFileOption) == 0 {

				//hasHistoryFileOption = true
				if !theHistoryFileOption[MUTEABLEKEY].(bool) {
					m = append(m, theHistoryFileOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultConfigurationFileOption) == 0 {

				//just replace the option anyways, so we dont have to do recursive or deep layer checks
				//hasConfigurationFileOption = true
				if !theConfigFileOption[MUTEABLEKEY].(bool) {
					m = append(m, theConfigFileOption)
				}
			}
			if strings.Compare(arg[LONGFLAGKEY].(string), _defaultVerbosityOption) == 0 {

				//hasVerbosityOption = true
				if !theVerbosityOption[MUTEABLEKEY].(bool) {
					m = append(m, theVerbosityOption)
				}
			}
		}
	}

	return m, pargs.hasInteractiveOption
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

func checkProgramTopLevel(m []map[string]interface{}) ([]map[string]interface{}, bool) {

	args, ok := checkProgramArgs(m)

	if !ok {
		return nil, false
	}

	for _, value := range m {

		if checkForDefaultProgramArgument(value[LONGFLAGKEY].(string)) {
			continue
		}

		newmap, _ok := checkProgramMap(value)
		if !_ok {
			ok = _ok
			continue
		}

		args = append(args, newmap)

	}

	return args, ok
}

func checkProgramMap(inputmap map[string]interface{}) (map[string]interface{}, bool) {

	resultmap := map[string]interface{}{}

	theType, ok := inputmap[TYPEKEY].(string)

	if !ok {
		return resultmap, false
	}

	theTypes := []string{theType}

	theTypes = splitBySeperator(theTypes)

	if len(theTypes) == 1 {

		_, _ok := checkParameter(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		_, _ok = checkOption(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		_, _ok = checkWildcard(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		_, _ok = checkFlag(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		_, _ok = checkCommand(theTypes[0], &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
		return resultmap, ok
	}
	if len(theTypes) > 1 {
		_, _ok := checkMultiArgument(theTypes, &resultmap, &inputmap)
		if !_ok {
			ok = _ok
		}
	}

	return resultmap, ok
}

func checkMultiArgument(types []string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {

	ok := true

	argtype := ArgumentType(0)

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
		return (*resultmap), false
	}

	if (argtype & PARAMETER) > 0 {
		checkParameter(PARAMETERSTRING, resultmap, inputmap)
	}
	if (argtype & OPTION) > 0 {
		checkOption(OPTIONSTRING, resultmap, inputmap)
	}
	if (argtype & WILDCARD) > 0 {
		checkWildcard(WILDCARDSTRING, resultmap, inputmap)
	}
	if (argtype & FLAG) > 0 {
		checkFlag(FLAGSTRING, resultmap, inputmap)
	}
	if (argtype & COMMAND) > 0 {
		checkCommand(COMMANDSTRING, resultmap, inputmap)
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

	return (*resultmap), ok
}

func checkParameter(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	if strings.Compare(theType, PARAMETERSTRING) == 0 {
		//this is a parameter

		if (*inputmap)[ARGUMENTSKEY] != nil {
			return (*resultmap), false
		}
		copyProgramArgument(resultmap, inputmap)
	}
	return (*resultmap), true
}

func checkOption(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {

	ok := true
	if strings.Compare(theType, OPTIONSTRING) != 0 {
		return (*resultmap), true
	}
	//this is an option
	if str, _ok := (*inputmap)[LONGFLAGKEY].(string); !_ok || len(str) <= 0 {
		return (*resultmap), false
	}

	copyProgramArgument(resultmap, inputmap)

	if (*inputmap)[ARGUMENTSKEY] == nil {
		return (*resultmap), true
	}
	(*resultmap)[ARGUMENTSKEY] = []map[string]interface{}{}

	switch (*inputmap)[ARGUMENTSKEY].(type) {
	case []interface{}:
		{
			for _, v := range (*inputmap)[ARGUMENTSKEY].([]interface{}) {
				arg, _ok := v.(map[string]interface{})
				if !_ok {
					ok = _ok
					continue
				}
				newarg, _ok := checkProgramMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
			}
			break
		}
	case []map[string]interface{}:
		{
			for _, arg := range (*inputmap)[ARGUMENTSKEY].([]map[string]interface{}) {
				newarg, _ok := checkProgramMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
				if !_ok {
					ok = _ok
				}
			}
			break
		}
	}
	return *resultmap, ok
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

func checkCommand(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {

	ok := true
	if strings.Compare(theType, COMMANDSTRING) != 0 {
		return (*resultmap), true
	}
	//this is an option
	if str, _ok := (*inputmap)[LONGFLAGKEY].(string); !_ok || len(str) <= 0 {
		return (*resultmap), false
	}

	if str, _ok := (*inputmap)[RUNKEY].(string); !_ok || len(str) <= 0 {
		return (*resultmap), false
	}

	copyProgramArgument(resultmap, inputmap)

	if (*inputmap)[ARGUMENTSKEY] == nil {
		return (*resultmap), true
	}
	(*resultmap)[ARGUMENTSKEY] = []map[string]interface{}{}

	switch (*inputmap)[ARGUMENTSKEY].(type) {
	case []interface{}:
		{
			for _, v := range (*inputmap)[ARGUMENTSKEY].([]interface{}) {
				arg, _ok := v.(map[string]interface{})
				if !_ok {
					ok = _ok
					continue
				}
				newarg, _ok := checkProgramMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
			}
			break
		}
	case []map[string]interface{}:
		{
			for _, arg := range (*inputmap)[ARGUMENTSKEY].([]map[string]interface{}) {

				newarg, _ok := checkProgramMap(arg)

				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)

				if !_ok {
					ok = _ok
				}
			}
			break
		}
	}
	return *resultmap, ok
}

func checkProgramInterfaceArray(i []interface{}) ([]map[string]interface{}, bool) {

	m := []map[string]interface{}{}
	var ok bool = true
	for _, v := range i {

		if arg, __ok := (v.(map[string]interface{})); __ok {

			_arg, _ok := checkProgramMap(arg)
			if !_ok {
				ok = _ok
				continue
			}
			m = append(m, _arg)
		} else {
			ok = __ok
		}
	}
	return m, ok
}
