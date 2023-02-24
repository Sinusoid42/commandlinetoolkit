package commandlinetoolkit

type CLICODE int32

type ATTRIBUTE int32

const SHELL ATTRIBUTE = 0b00000000001
const HISTORY ATTRIBUTE = 0b00000000010
const HISTORYFILE ATTRIBUTE = 0b00000000100
const SUGGESTIONS ATTRIBUTE = 0b00000001000
const PREDICTIONS ATTRIBUTE = 0b00000010000
const LOGGING ATTRIBUTE = 0b00000100000

var CLI_NULL CLICODE = 0b00000000000000

const CLI_FALSE CLICODE = 0b000000000000
const CLI_ERROR CLICODE = 0b000000000000

/*
NO commandline error was detected, exiting with status 0, just like any other binary program
*/
const CLI_NO_ERROR CLICODE = 0b00000000000001
const CLI_SUCCESS CLICODE = 0b00000000000001
const CLI_TRUE CLICODE = 0b000000000001

/*
The input given to the commandline was unknown, either there was no definition given at all, or the input is not useable
Is thrown when a wrong option or command is entered into the program, if a given parameter is requested this program
exists with invalid input, required argument not found or wrong data code
*/
const CLI_UNKOWNINPUT_ERROR CLICODE = 0b00000000000010

/*
If a given input is invalid, when a parameter is parsed and there was a positional argument, but a different one was given
*/
const CLI_INVALIDINPUT_ERROR CLICODE = 0b00000000000100

/*
A trailing argument was required to complete the parsing of the tree, but the required argument could not be detected
*/
const CLI_REQARGUMENTNOTFOUND_ERROR CLICODE = 0b00000000001000

/*
The provided data given to the commandline was incorrect
Either the callback does not work, so you need to checkInputProgram your custom datatype callback for a given argument
Or you have to recheck your defined json file and read it again
*/
const CLI_WRONGDATA_ERROR CLICODE = 0b00000000010000

/*
The command line exists because of a help wildcard
*/
const CLI_HELPEXIT_ERROR CLICODE = 0b00000000100000

/*
In the recursive parsing, when a argument is given by the userinput and the argument is not found in that
tree level, throws a argument not found error
*/
const CLI_ARGUMENTNOTFOUND_ERROR CLICODE = 0b00000001000000

/*
If the builder tries to overwrite a default or protected argument, this error is thrown
Upon booting, the commandlinetoolkit will eject
If interactive, the commandlinetool will just not read the config and throw an error but resume operation
*/
const CLI_ARGUMENTNOTOVERWRITEABLE_ERROR CLICODE = 0b00000010000000
const CLI_ARGUMENTPROTECTED_ERROR CLICODE = 0b0000001000000

/*
This error is thrown when the commmandline is in interactive mode and experiences a fatal error of known or unkown origin
Can be returned by a callback and will halt the stop execution of the entire program, similar to System.exit()
*/
const CLI_RUNTIME_ERROR CLICODE = 0b00000010000000

const CLI_LEXING_ERROR CLICODE = 0b00000100000000

/*
***********************************************************************************************************************
Interacrive Shell Mode

providing a --interactive or -i when starting the commandline will lead to an executeable binary to keep running in a waitgroup shell struct
that handles constant and continuous input via the commandline input structure, (in future that can be overwritten by using a restcall api)
*
	//see @defaults

*/

// when using the TAB and there is no predction available, usually we want to display all possible inputs for the given parse tree layer
const CLI_NO_PREDICTION_ERROR CLICODE = 0b00001000000000

const CLI_EXIT CLICODE = 0b100000000000

const CLI_VERBOSE_SIMPLE CLICODE = 0b000000001

const CLI_VERBOSE_COMPLEX CLICODE = 0b000000010

const CLI_VERBOSE_PARSING CLICODE = 0b000000100

const CLI_VERBOSE_SHELL CLICODE = 0b000001000

const CLI_VERBOSE_OS_SIG CLICODE = 0b000010000

const CLI_VERBOSE_SHELL_PARSE CLICODE = 0b000100000

const CLI_VERBOSE_PREDICT CLICODE = 0b001000000

const CLI_VERBOSE_SHELL_BUFFER CLICODE = 0b010000000

const CLI_VERBOSE_FILE CLICODE = 0b100000000

const CLI_VERBOSE_DEBUG CLICODE = 0b1000000000

const CLI_VERBOSE_PROGRAM CLICODE = 0b10000000000
