package commandlinetoolkit

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
Helper file that represents a runnable program,

	_programName        The name of the executeable file, is parsed when getting os.Args from the first run
	_programFile        The name of the configuration json file, that holds all the available Arguments for the CLI Parser
	_program            The command line template struct, that holds the interface storage for the actual program
	_verbose            A given debugging verbosity that is passed to the debug handler
	_debugHandler       The debug handler that manages the printing and output
*/

type program struct {
	_programName  string
	_programFile  string
	_program      *commandlinetemplate
	_verbose      CLICODE
	_debugHandler *debugHandler
}

type programArgs struct {
	hasInteractiveOption       bool
	hasLoggingOption           bool
	hasHistoryOption           bool
	hasHelpOption              bool
	hasHistoryFileOption       bool
	hasConfigurationFileOption bool
	hasVerbosityOption         bool
}

/*
********************************************************************

Builds a new Program
*/
func newprogram(filename string) *program {
	
	p := &program{
		_programName:  "Command Line: " + VERSION,
		_programFile:  filename,
		_program:      DefaultCommandLineTemplate(),
		_debugHandler: newDebugHandler(),
	}
	
	p._debugHandler._verbose = 0
	
	return p
}

/*
********************************************************************

Read a json cli configuration file
*/
func (p *program) readJsonProgram(filename string) string {
	
	//check for file forst
	if _, err := os.OpenFile(filename, os.O_RDONLY, 0666); err != nil {
		filename = p._programFile
		p._debugHandler.printError("-->readJson: File not available, falling back to default\n")
	}
	
	p._programFile = filename
	_programFile, _error := os.OpenFile(p._programFile, os.O_CREATE, 0666)
	if _error != nil {
		p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "-->readJson: Error\n")
	} else {
		p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "-->readJson: Read Success\n")
	}
	
	fileInfo, _ := _programFile.Stat()
	fileSize := fileInfo.Size()
	
	data := make([]byte, fileSize)
	
	p._program._theProgramJsonMap = make(map[string]interface{})
	
	_, _ = _programFile.Read(data)
	
	text := string(data)
	
	err := json.Unmarshal([]byte(text), &p._program._theProgramJsonMap)
	
	if err != nil {
		if len(text) > 0 {
			
			saveErrorFile(p._programFile, text)
			p._program._theProgramJsonMap = DefaultTemplate()
			p._debugHandler.printError("-->readJson: Error while reading the JsonProgram: writing Default\n-->readJson: savefile created!\n")
			p.write()
		}
		
	} else {
		p._programFile = filename
		
		p._program._theProgramJson = text
		
		p._program._theProgramJsonMap = make(map[string]interface{})
		
		err := json.Unmarshal([]byte(text), &p._program._theProgramJsonMap)
		
		if err != nil {
			p._program._theProgramJsonMap = DefaultTemplate()
			
			p._debugHandler.printVerbose(CLI_VERBOSE_PROGRAM, "-->readJson: Error Reading: writing Default\n")
			
			p.write()
		}
	}
	_programFile.Close()
	
	return text
}

/*
********************************************************************

Write the json to disc
*/
func (p *program) write() {
	
	b, _ := json.MarshalIndent(p._program._theProgramJsonMap, " ", "   ")
	
	p._program._theProgramJson = string(b)
	
	_programFile, _ := os.OpenFile(p._programFile, os.O_RDWR|os.O_TRUNC, 0666)
	
	writer := bufio.NewWriter(_programFile)
	
	writer.WriteString(p._program._theProgramJson)
	
	writer.Flush()
	
	_programFile.Close()
}

/*
********************************************************************

Do a first tree sweep check via typecast checking
if coliding parameters are present, they are excluded,savefiles will be written
*/
func (p *program) check() {
	
	//check for required defaults, that are required to run the cli, if the contents are not present parsing will not hapen
	allArgs := p._program._theProgramJsonMap[ARGUMENTSKEY]
	
	//we expect here either a well formed json will default immuteable (replaced if removed) arguments
	//on simple read error we create a new json with the given commands
	//if successfull todo is check if the structure is well formed
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
			//check recursively
			checkedPrgm, ok = checkTopLevel(allArgs.([]map[string]interface{}))
		}
	case []interface{}:
		{
			//check recursively
			checkedPrgm, ok = parseInterfaceArray(allArgs.([]interface{}))
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
		
		p.write()
	}
	
	p._program._theProgramJsonMap["arguments"] = checkedPrgm
	
	p.write()
	
}

func saveErrorFile(fileName string, content string) {
	
	yr, month, day := time.Now().Date()
	hr := time.Now().Hour()
	min := time.Now().Minute()
	
	theTime := strconv.Itoa(yr) + "_" + strconv.Itoa(int(month)) + "_" + strconv.Itoa(day) + "_" + strconv.Itoa(hr) + "_" + strconv.Itoa(min) + "_"
	
	file, _ := os.OpenFile(theTime+"errsave_"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	
	writer := bufio.NewWriter(file)
	
	writer.WriteString(content)
	
	writer.Flush()
	
	file.Close()
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

func isDefault(str string) bool {
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

func checkTopLevel(m []map[string]interface{}) ([]map[string]interface{}, bool) {
	
	args, ok := checkProgramArgs(m)
	
	if !ok {
		return nil, false
	}
	
	for _, value := range m {
		
		if isDefault(value[LONGFLAGKEY].(string)) {
			continue
		}
		
		newmap, _ok := parseMap(value)
		if !_ok {
			ok = _ok
			continue
		}
		
		args = append(args, newmap)
		
	}
	
	return args, ok
}

func parseMap(inputmap map[string]interface{}) (map[string]interface{}, bool) {
	
	resultmap := map[string]interface{}{}
	
	theType, ok := inputmap[TYPEKEY].(string)
	
	if !ok {
		return resultmap, false
	}
	
	theTypes := []string{theType}
	
	theTypes = splitBySeperator(theTypes)
	
	if len(theTypes) == 1 {
		
		_, _ok := parseParameter(theTypes[0], &resultmap, &inputmap)
		if !ok {
			ok = _ok
		}
		_, _ok = parseOption(theTypes[0], &resultmap, &inputmap)
		if !ok {
			ok = _ok
		}
		_, _ok = parseWildcard(theTypes[0], &resultmap, &inputmap)
		if !ok {
			ok = _ok
		}
		_, _ok = parseFlag(theTypes[0], &resultmap, &inputmap)
		if !ok {
			ok = _ok
		}
		_, _ok = parseCommand(theTypes[0], &resultmap, &inputmap)
		if !ok {
			ok = _ok
		}
		
	}
	
	return resultmap, ok
}

func parseParameter(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	if strings.Compare(theType, PARAMETERSTRING) == 0 {
		//this is a parameter
		
		if (*inputmap)[ARGUMENTSKEY] != nil {
			return (*resultmap), false
		}
		copyArg(resultmap, inputmap)
	}
	return (*resultmap), true
}

func parseOption(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	
	ok := true
	if strings.Compare(theType, OPTIONSTRING) != 0 {
		return (*resultmap), true
	}
	//this is an option
	if str, _ok := (*inputmap)[LONGFLAGKEY].(string); !_ok || len(str) <= 0 {
		return (*resultmap), false
	}
	
	copyArg(resultmap, inputmap)
	
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
				newarg, _ok := parseMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
			}
			break
		}
	case []map[string]interface{}:
		{
			for _, arg := range (*inputmap)[ARGUMENTSKEY].([]map[string]interface{}) {
				newarg, _ok := parseMap(arg)
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

func parseWildcard(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	if strings.Compare(theType, WILDCARDSTRING) == 0 {
		if str, _ok := (*inputmap)[LONGFLAGKEY].(string); _ok && len(str) > 0 {
			copyArg(resultmap, inputmap)
		}
		return (*resultmap), true
	}
	return (*resultmap), true
}

func parseFlag(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	if strings.Compare(theType, FLAGSTRING) == 0 {
		if str, _ok := (*inputmap)[LONGFLAGKEY].(string); _ok && len(str) > 0 {
			copyArg(resultmap, inputmap)
		}
		return (*resultmap), true
	}
	return (*resultmap), true
}

func parseCommand(theType string, resultmap *map[string]interface{}, inputmap *map[string]interface{}) (map[string]interface{}, bool) {
	
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
	
	copyArg(resultmap, inputmap)
	
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
				newarg, _ok := parseMap(arg)
				(*resultmap)[ARGUMENTSKEY] = append((*resultmap)[ARGUMENTSKEY].([]map[string]interface{}), newarg)
			}
			break
		}
	case []map[string]interface{}:
		{
			for _, arg := range (*inputmap)[ARGUMENTSKEY].([]map[string]interface{}) {
				
				newarg, _ok := parseMap(arg)
				
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

func copyArg(dst *map[string]interface{}, src *map[string]interface{}) {
	_copyArg(dst, src, TYPEKEY)
	
	_copyArg(dst, src, LONGFLAGKEY)
	_copyArg(dst, src, SHORTFLAGKEY)
	_copyArg(dst, src, HELPKEY)
	_copyArg(dst, src, SHORTHELPKEY)
	_copyArg(dst, src, DATATYPEKEY)
	_copyArg(dst, src, RUNKEY)
	_copyArg(dst, src, REQUIREDKEY)
	_copyArg(dst, src, MUTEABLEKEY)
	
	//do this with check casting
	//_copyArg(dst, src, ARGUMENTSKEY] = (*src)[ARGUMENTSKEY] since here we have no string
	_copyArg(dst, src, AUTHORKEY)
	_copyArg(dst, src, DESCRIPTIONKEY)
	_copyArg(dst, src, MANUALKEY)
	_copyArg(dst, src, VERSIONKEY)
	_copyArg(dst, src, EXECUTEABLEKEY)
}

func _copyArg(dst *map[string]interface{}, src *map[string]interface{}, typeString string) {
	if str, ok := (*src)[typeString].(string); ok {
		(*dst)[typeString] = str
	}
	if str, ok := (*src)[typeString].(int); ok {
		(*dst)[typeString] = str
	}
	if str, ok := (*src)[typeString].(bool); ok {
		(*dst)[typeString] = str
	}
}

func splitBySeperator(str []string) []string {
	
	//split
	if strings.Contains(str[0], "|") {
		str = strings.Split(str[0], "|")
	}
	
	if strings.Contains(str[0], ",") {
		str = strings.Split(str[0], ",")
	}
	
	for i, s := range str {
		
		if strings.Contains(s, " ") {
			strings.ReplaceAll(s, " ", "")
		}
		str[i] = s
	}
	return str
}

func parseInterfaceArray(i []interface{}) ([]map[string]interface{}, bool) {
	
	m := []map[string]interface{}{}
	var ok bool = true
	for _, v := range i {
		
		if arg, ok := (v.(map[string]interface{})); ok {
			
			_arg, _ok := parseMap(arg)
			
			if !_ok {
				ok = _ok
				continue
			}
			m = append(m, _arg)
		} else {
		
		}
	}
	return m, ok
}
