package commandlinetoolkit

type ArgumentType int32

type ArgumentDataType string

var FILETYPE ArgumentDataType = "file" //the software checks wether the file exists and can be read
var URLTYPE ArgumentDataType = "url"   //the software checks, wether the given string is a valid url (ssh|http|https or <name>@<ip>
var STRINGTYPE ArgumentDataType = "string"
var NUMBERTYPE ArgumentDataType = "number"
var CUSTOMTYPE ArgumentDataType = "custom"

/*
	Define all the different ArgumentType required to be parsed inside of the argument command line
*/

const OPTION ArgumentType = 0b0000000000000001
const PARAMETER ArgumentType = 0b0000000000000010
const WILDCARD ArgumentType = 0b0000000000000100
const FLAG ArgumentType = 0b0000000000001000
const COMMAND ArgumentType = 0b0000000000010000
const __NULL_ARG__ ArgumentType = 0b0000000000100000

/*
	This is how a command line argument is per default defined

		arg_type : The type of the argument
					Defines behaviour, Parsing differences between methods, options and parameters, flags/wildcars and options all are with leading dash/dashes


*/

type Argument struct {
	arg_type  ArgumentType
	data_type ArgumentDataType
	
	custom_type ArgumentDataType
	
	lflag string
	sflag string
	lhelp string
	shelp string
	
	choices []string //is generated dynamically when parsing trailing parameters for a given command
	
	callback func() CLICODE //is called when this argument is parsed
	
	command func(parameters []*Argument, options []*Argument)
	
	arguments []*Argument //All available subarguments or commands
	
	parent *Argument //The parent argument
}
