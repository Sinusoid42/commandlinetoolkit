package commandlinetoolkit

type ArgumentType int32

type ArgumentDataType string

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
	VERSIONKEY            = "version"
	EXECUTEABLEKEY        = "executeable"
)

var FILETYPE ArgumentDataType = "file" //the software checks wether the file exists and can be read
var URLTYPE ArgumentDataType = "url"   //the software checks, wether the given string is a valid url (ssh|http|https or <name>@<ip>
var STRINGTYPE ArgumentDataType = "string"
var NUMBERTYPE ArgumentDataType = "number"
var CUSTOMTYPE ArgumentDataType = "custom"
var BOOLTYPE ArgumentDataType = "bool"

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
	data_type ArgumentDataType

	//the custom type
	custom_type ArgumentDataType

	//a custom callback function, if a given ArgumentDataType is custom
	dtype_custom_callback func(arg string) bool

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

	choices []string //is generated dynamically when parsing trailing parameters for a given command

	//the callback function that is called when the argument is parsed
	callback func() CLICODE

	//after parsing, the arguments array will be filled with trailing !new! parameters or options instances
	//will only sublayers from the directed subgraph, will not contain parsed arguments from other subtrees of the parsetree
	arguments []*Argument

	//the run function for the given argument in case the argument has a 'run' keyword
	//in case a 'self' is provided, a custom run command can be provided for the argument that will be run after building the parsetree
	run func(parameters []*Argument, options []*Argument, cmdline *CommandLine) CLICODE //TODO
}

func argtypeString(argumentType ArgumentType) string {
	switch argumentType {
	case OPTION:
		{
			return OPTIONSTRING
		}
	case COMMAND:
		{
			return COMMANDSTRING
		}
	case WILDCARD:
		{
			return WILDCARDSTRING
		}
	case FLAG:
		{
			return FLAGSTRING
		}
	case PARAMETER:
		{
			return PARAMETERSTRING
		}
	}
	return ""
}
