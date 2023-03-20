package commandlinetoolkit

import (
	"bufio"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
)

type ArgumentType int32

type ArgumentDataType struct {
	data_flag             string
	dtype_custom_callback func(a *ArgumentDataType, arg string) (any, bool)
	data                  any
	attrib                string
}

const (
	TYPEKEY        string = "type"
	LONGFLAGKEY           = "flag"
	SHORTFLAGKEY          = "sflag"
	HELPKEY               = "help"
	SHORTHELPKEY          = "shelp"
	DATATYPEKEY           = "datatype"
	RUNKEY                = "run"
	REQUIREDKEY           = "required"
	MUTEABLEKEY           = "muteable"
	ARGUMENTSKEY          = "arguments"
	AUTHORKEY             = "author"
	DESCRIPTIONKEY        = "description"
	MANUALKEY             = "man"
	TITLEKEY              = "title"
	VERSIONKEY            = "version"
	EXECUTEABLEKEY        = "executeable"
	STYLETITLEKEY         = "styleTitle"
)

const FULLOPTIONPREFIX string = "--"
const SHORTOPTIONPREFIX string = "-"

var FILETYPE ArgumentDataType = ArgumentDataType{"file", checkForFile, nil, ""} //the software checks wether the file exists and can be read
var URLTYPE ArgumentDataType = ArgumentDataType{"url", checkForURL, nil, ""}    //the software checks, wether the given string is a valid url (ssh|http|https or <name>@<ip>
var STRINGTYPE ArgumentDataType = ArgumentDataType{"string", checkForString, nil, ""}
var NUMBERTYPE ArgumentDataType = ArgumentDataType{"number", checkForNumber, nil, ""}
var CUSTOMTYPE ArgumentDataType = ArgumentDataType{"custom", nil, nil, ""}
var BOOLTYPE ArgumentDataType = ArgumentDataType{"bool", checkForBool, nil, ""}

/*
	Define all the different ArgumentType required to be parsed inside of the argument command line
*/

const OPTION ArgumentType = 0b0000000000000001
const PARAMETER ArgumentType = 0b0000000000000010
const WILDCARD ArgumentType = 0b0000000000000100
const FLAG ArgumentType = 0b0000000000001000
const COMMAND ArgumentType = 0b0000000000010000
const __NULL_ARG__ ArgumentType = 0b0000000000100000

const OPTIONSTRING = "OPTION"
const PARAMETERSTRING = "PARAM"
const WILDCARDSTRING = "WILDCARD"
const FLAGSTRING = "FLAG"
const COMMANDSTRING = "COMMAND"

/*
	This is how a command line argument is per default defined

		arg_type : The type of the argument
					Defines behaviour, Parsing differences between methods, options and parameters, flags/wildcars and options all are with leading dash/dashes


*/

type Argument struct {

	//the type of this argument
	arg_type ArgumentType

	//a string representing the datatype of the argument, if custom, data_type is 'custom' and custom_type holds the to be checked string
	data_type *ArgumentDataType

	//a custom callback function, if a given ArgumentDataType is custom

	//bools for conditional parsing of the argument
	required bool //is requried
	muteable bool //is 'overwriteable' or default (see defaults.go & program.go)(is relevant for the parsetree and program file)

	//the long flag of the argument or command
	lflag string

	//the short flag, non parametrized should allow combinations of nonexclusive short flags (-h exits)
	sflag string

	//the long help for the --help menu
	lhelp string

	//the short help msg for the -h menu
	shelp string

	runCommand string

	choices []string //is generated dynamically when parsing trailing parameters for a given command

	//the callback function that is called when the argument is parsed
	callback func() CLICODE

	paramCallback func(string) bool

	_sub []*Argument

	//after parsing, the arguments array will be filled with trailing !new! parameters or options instances
	//will only sublayers from the directed subgraph, will not contain parsed arguments from other subtrees of the parsetree

	//@parameters stores the parameters for a given subcommand => allows multiple executions of the last sub cmd
	//arguments store all options and --verbose etc. wildcards that the

	//TODO add exec to program builder so we can run other binaries just by json alone

	//the run function for the given argument in case the argument has a 'run' keyword
	//in case a 'self' is provided, a custom run command can be provided for the argument that will be run after building the parsetree
	run func(parameters []*Argument, arguments []*Argument, cmdline *CommandLine) CLICODE //TODO
}

func NewArgument(m map[string]interface{}) (*Argument, error) {
	//checkInputProgram the argument type first
	argType, err := decodeArgType(m)

	if err != nil {
		d := newDebugHandler()
		d.printError("-->Error: ParseTree: decodeArgType")
		return nil, errors.New("-->Error: ParseTree: decodeArgType")
	}
	arg := &Argument{
		data_type: createArgDataType(""),
		_sub:      []*Argument{},
	}
	//here we have an option, every time the parser reads leading '--' '-'
	//if we also have option, wildcard => apply option globally, if option and flag, the argument is consumed immideately globally once
	if argType&OPTION > 0 {
		arg, err = createOptionArgument(argType, m)

	}
	//with this we set an option, that is given to every subroutine or sub-execution recursively, only in the current tree of commands
	if argType&WILDCARD > 0 {
		arg, err = createWildcardArgument(argType, m)
	}
	//something we use, that alters the entire pace of the program, globally acceptable parameter
	if argType&FLAG > 0 {
		arg, err = createFlagArgument(argType, m)
	}
	//actually run something, a subroutine or a post, tree building added callback method with given options and parameters
	if argType&COMMAND > 0 {
		arg, err = createCommandArgument(argType, m)
	}
	//a parameter, if we have an option and a single paramter trailing then we shall allow also
	//[option=value] we can also allow positional:[option VALUE] or nonpositional:[option paramname=VALUE]
	//can only be interpreted at the first paramter, if multiple paramters are required, positional or nonpositional cant be interchanged
	//if no name for the paramter is given eg. flag, than  the name of the param is the datatype

	//internal parsing rules
	//if the datatype is number, we allow interval [a:b], [:b] or [a:]
	//if the datatype is url, we check with http? maybe ?
	//if the datatype is string, we allow the paramcheck(p string) custom callback function to be evalutated and store it afterwards
	//if the datatype is bool, just store bool value
	//if the datatype is file, we check if file is present and read its contents into the variable as string
	if argType&PARAMETER > 0 {

		arg, err = createParameterArgument(argType, m)

	}

	return arg, err
}

func (arg *Argument) AddArgument(newarg *Argument) CLICODE {

	for _, basearg := range arg._sub {
		if strings.Compare(basearg.lflag, arg.lflag) == 0 {

			newDebugHandler().printError("Could not add Argument as it overlaps with another Argument")

			return CLI_ERROR
		}
	}

	if arg._sub == nil {
		arg._sub = []*Argument{}
	}
	arg._sub = append(arg._sub, newarg)

	return CLI_SUCCESS
}

func (a *Argument) GetValue() any {

	if a.data_type == nil {
		return nil
	}
	return a.data_type.data

}

func (a *Argument) copy() *Argument {

	if a == nil {
		return nil
	}

	arg := &Argument{
		arg_type: a.arg_type,
		data_type: &ArgumentDataType{
			data:                  a.data_type.data,
			data_flag:             a.data_type.data_flag,
			attrib:                a.data_type.attrib,
			dtype_custom_callback: a.data_type.dtype_custom_callback,
		},

		required:      a.required,
		muteable:      a.muteable,
		lflag:         a.lflag,
		sflag:         a.sflag,
		lhelp:         a.lhelp,
		shelp:         a.shelp,
		runCommand:    a.runCommand,
		callback:      a.callback,
		paramCallback: a.paramCallback,
		run:           a.run,
	}

	arg.choices = make([]string, len(a.choices))
	copy(arg.choices, a.choices)

	return arg

}

func createOptionArgument(argType ArgumentType, m map[string]interface{}) (*Argument, error) {

	//if we have an option, at least a long flag is required, that is already checked before
	longFlag, ok := m[LONGFLAGKEY].(string)

	if !ok {
		return nil, errors.New("Flag not available")
	}

	arg := &Argument{
		arg_type:  argType,
		lflag:     longFlag,
		data_type: createArgDataType(""),
	}
	if shortFlag, _ok := m[SHORTFLAGKEY].(string); _ok {
		arg.sflag = shortFlag
	}
	if help, _ok := m[HELPKEY].(string); _ok {
		arg.lhelp = help
	}
	if shelp, _ok := m[SHORTHELPKEY].(string); _ok {
		arg.shelp = shelp
	}
	if required, _ok := m[REQUIREDKEY].(bool); _ok {
		arg.required = required
	}
	if muteable, _ok := m[MUTEABLEKEY].(bool); _ok {
		arg.muteable = muteable
	}
	if data_type, _ok := m[DATATYPEKEY].(string); _ok {
		arg.data_type = createArgDataType(data_type)
	} else {
		arg.data_type = &ArgumentDataType{}
	}
	if runCmd, _ok := m[RUNKEY].(string); _ok {
		arg.runCommand = runCmd

		if strings.ContainsAny(runCmd, ",") {
			cmds := strings.Split(runCmd, ", ")
			for i, c := range cmds {
				cmds[i] = c
			}
			arg.run = getRunCommands(cmds)
		} else {
			if isLibCommand(runCmd) {
				arg.run = getRunCommand(runCmd)
			} else {
				arg.run = nil
			}
		}

	}
	return arg, nil
}

func createWildcardArgument(argType ArgumentType, m map[string]interface{}) (*Argument, error) {

	//if we have an option, at least a long flag is required, that is already checked before
	longFlag, ok := m[LONGFLAGKEY].(string)

	if !ok {
		return nil, errors.New("Flag not available")
	}

	arg := &Argument{
		arg_type:  argType,
		lflag:     longFlag,
		data_type: createArgDataType(""),
	}

	if shortFlag, _ok := m[SHORTFLAGKEY].(string); _ok {
		arg.sflag = shortFlag
	}
	if help, _ok := m[HELPKEY].(string); _ok {
		arg.lhelp = help
	}
	if shelp, _ok := m[SHORTHELPKEY].(string); _ok {
		arg.shelp = shelp
	}
	if required, _ok := m[REQUIREDKEY].(bool); _ok {
		arg.required = required
	}
	if muteable, _ok := m[MUTEABLEKEY].(bool); _ok {
		arg.muteable = muteable
	}
	if data_type, _ok := m[DATATYPEKEY].(string); _ok {

		arg.data_type = createArgDataType(data_type)
	} else {
		arg.data_type = &ArgumentDataType{}
	}

	return arg, nil

}

func createCommandArgument(argType ArgumentType, m map[string]interface{}) (*Argument, error) {

	//if we have an option, at least a long flag is required, that is already checked before
	longFlag, ok := m[LONGFLAGKEY].(string)

	if !ok {
		return nil, errors.New("Flag not available")
	}

	runCmd, _ok := m[RUNKEY].(string)
	if !_ok {

		//need to make sure, there is a callback or run that at least

		//return nil, errors.New("Flag not available")
	}

	arg := &Argument{
		arg_type:   argType,
		lflag:      longFlag,
		runCommand: runCmd,
		data_type:  createArgDataType(""),
	}

	if shortFlag, _ok := m[SHORTFLAGKEY].(string); _ok {
		arg.sflag = shortFlag
	}
	if help, _ok := m[HELPKEY].(string); _ok {
		arg.lhelp = help
	}
	if shelp, _ok := m[SHORTHELPKEY].(string); _ok {
		arg.shelp = shelp
	}
	if required, _ok := m[REQUIREDKEY].(bool); _ok {
		arg.required = required
	}
	if muteable, _ok := m[MUTEABLEKEY].(bool); _ok {
		arg.muteable = muteable
	}
	if data_type, _ok := m[DATATYPEKEY].(string); _ok {

		arg.data_type = createArgDataType(data_type)
	} else {
		arg.data_type = &ArgumentDataType{}
	}

	return arg, nil
}
func createFlagArgument(argType ArgumentType, m map[string]interface{}) (*Argument, error) {
	//if we have an option, at least a long flag is required, that is already checked before
	longFlag, ok := m[LONGFLAGKEY].(string)

	if !ok {
		return nil, errors.New("Flag not available")
	}

	arg := &Argument{
		arg_type:  argType,
		lflag:     longFlag,
		data_type: createArgDataType(""),
	}
	if shortFlag, _ok := m[SHORTFLAGKEY].(string); _ok {
		arg.sflag = shortFlag
	}
	if help, _ok := m[HELPKEY].(string); _ok {
		arg.lhelp = help
	}
	if help, _ok := m[HELPKEY].(string); _ok {
		arg.lhelp = help
	}
	if shelp, _ok := m[SHORTHELPKEY].(string); _ok {
		arg.shelp = shelp
	}
	if required, _ok := m[REQUIREDKEY].(bool); _ok {
		arg.required = required
	}
	if muteable, _ok := m[MUTEABLEKEY].(bool); _ok {
		arg.muteable = muteable
	}
	if data_type, _ok := m[DATATYPEKEY].(string); _ok {
		arg.data_type = createArgDataType(data_type)
	} else {
		arg.data_type = &ArgumentDataType{}
	}

	return arg, nil
}

func createParameterArgument(argType ArgumentType, m map[string]interface{}) (*Argument, error) {

	arg := &Argument{
		arg_type:  argType,
		data_type: createArgDataType(""),
	}
	if longFlag, _ok := m[SHORTFLAGKEY].(string); _ok {
		arg.lflag = longFlag
	}
	if shortFlag, _ok := m[SHORTFLAGKEY].(string); _ok {
		arg.sflag = shortFlag
	}
	if help, _ok := m[HELPKEY].(string); _ok {
		arg.lhelp = help
	}
	if help, _ok := m[HELPKEY].(string); _ok {
		arg.lhelp = help
	}
	if shelp, _ok := m[SHORTHELPKEY].(string); _ok {
		arg.shelp = shelp
	}
	if required, _ok := m[REQUIREDKEY].(bool); _ok {
		arg.required = required
	}
	if muteable, _ok := m[MUTEABLEKEY].(bool); _ok {
		arg.muteable = muteable
	}
	if data_type, _ok := m[DATATYPEKEY].(string); _ok {
		arg.data_type = createArgDataType(data_type)
	} else {
		arg.data_type = &ArgumentDataType{}
	}
	return arg, nil
}

func createArgDataType(dtype string) *ArgumentDataType {
	if strings.Compare(dtype, FILETYPE.data_flag) == 0 {
		return &ArgumentDataType{FILETYPE.data_flag, checkForFile, nil, FILETYPE.attrib}
	}
	if strings.Compare(dtype, URLTYPE.data_flag) == 0 {
		return &ArgumentDataType{URLTYPE.data_flag, checkForURL, nil, FILETYPE.attrib}
	}
	if strings.Compare(dtype, NUMBERTYPE.data_flag) == 0 {
		return &ArgumentDataType{NUMBERTYPE.data_flag, checkForNumber, nil, FILETYPE.attrib}
	}
	if strings.Compare(dtype, BOOLTYPE.data_flag) == 0 {
		return &ArgumentDataType{BOOLTYPE.data_flag, checkForBool, nil, FILETYPE.attrib}
	}
	if strings.Compare(dtype, STRINGTYPE.data_flag) == 0 {
		return &ArgumentDataType{STRINGTYPE.data_flag, checkForString, nil, FILETYPE.attrib}
	}
	if strings.Index(dtype, NUMBERTYPE.data_flag) == 0 && strings.Index(dtype, "[") > 0 && strings.Index(dtype, "]") > 0 {
		theoption := strings.TrimPrefix(dtype, NUMBERTYPE.data_flag)
		return &ArgumentDataType{NUMBERTYPE.data_flag, checkForNumber, nil, theoption}
	}
	return nil
}

func argTypeToString(argumentType ArgumentType) string {
	s := ""

	if argumentType&OPTION > 0 {
		s += OPTIONSTRING
	}
	if argumentType&COMMAND > 0 {
		if len(s) > 0 {
			s += " | "
		}
		s += COMMANDSTRING
	}
	if argumentType&WILDCARD > 0 {
		if len(s) > 0 {
			s += " | "
		}
		s += WILDCARDSTRING
	}
	if argumentType&FLAG > 0 {
		if len(s) > 0 {
			s += " | "
		}
		s += FLAGSTRING
	}
	if argumentType&PARAMETER > 0 {
		if len(s) > 0 {
			s += " | "
		}
		s += PARAMETERSTRING
	}

	return s
}

func (a *Argument) String() string {

	s := ""
	if a.data_type != nil {
		s = "   dtype: " + a.data_type.data_flag + "\n"
	}

	return "Argument:{\n" +
		"   type: " + argTypeToString(a.arg_type) + "\n" +
		"   flag: " + a.lflag + "\n" + s +
		"   required: " + strconv.FormatBool(a.required) + "\n" +
		"   muteable: " + strconv.FormatBool(a.required) + "\n" +
		"}"
}

func checkForFile(a *ArgumentDataType, str string) (any, bool) {

	if file, err := os.OpenFile(str, os.O_RDWR, 0666); err == nil {
		// Create a scanner to read the file
		scanner := bufio.NewScanner(file)

		// Use a string builder to concatenate the lines into a single string
		var builder strings.Builder
		for scanner.Scan() {
			builder.WriteString(scanner.Text())
			builder.WriteString("\n")
		}

		// Print the resulting string
		textcontent := builder.String()
		a.data = textcontent
		file.Close()

		s := []string{}
		s = append(s, str)
		s = append(s, textcontent)
		return s, true
	}

	return "", false
}

func checkForURL(a *ArgumentDataType, str string) (any, bool) {

	if strings.Index(str, "http") == 0 {
		return true, true
	}

	if strings.Index(str, "https") == 0 {
		return true, true
	}

	//url ip has to have 4 dot + 1 : optinoally if we have an url or ipaddress
	if strings.Count(str, ".") == 4 && strings.Count(str, ":") <= 1 {
		return true, true
	}
	return false, false
}

func checkForNumber(a *ArgumentDataType, str string) (any, bool) {
	if len(a.attrib) == 0 {
		if intNbr, ok := strconv.Atoi(str); ok == nil {
			return intNbr, true
		}
		if floatNbr, ok := strconv.ParseFloat(str, 32); ok == nil {
			return floatNbr, true
		}
	}
	if len(str) < 1 {
		return 0, false
	}
	bounds := strings.Split(a.attrib, ":")
	if len(bounds) > 1 {
		bounds[0] = strings.TrimPrefix(bounds[0], "[")
		bounds[1] = strings.TrimSuffix(bounds[1], "]")
		min, ok := strconv.ParseFloat(bounds[0], 32)

		if ok != nil {
			min = math.Inf(-1)
		}
		max, ok := strconv.ParseFloat(bounds[1], 32)
		if ok != nil {
			max = math.Inf(1)
		}
		if intNbr, ok := strconv.Atoi(str); ok == nil {
			return intNbr, intNbr >= int(min) && intNbr <= int(max)
		}
		if floatNbr, ok := strconv.ParseFloat(str, 32); ok == nil {
			return floatNbr, floatNbr >= min && floatNbr <= max
		}
	}

	return 0.0, false
}

func checkForString(a *ArgumentDataType, str string) (any, bool) {
	return str, true
}

func checkForBool(a *ArgumentDataType, str string) (any, bool) {
	switch strings.ToLower(str) {
	case "true":
		{
			return true, true
		}
	case "false":
		{
			return false, true
		}
	}
	return false, false
}
