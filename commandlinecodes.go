package commandlinetoolkit

type CLICODE int32

/*
	NO commandline error was detected, exiting with status 0, just like any other binary program
*/
const CLI_NO_ERROR CLICODE = 0
const CLI_SUCCESS CLICODE = 0

/*
	The input given to the commandline was unknown, either there was no definition given at all, or the input is not useable
	Is thrown when a wrong option or command is entered into the program, if a given parameter is requested this program
	exists with invalid input, required argument not found or wrong data code
*/
const CLI_UNKOWNINPUT_ERROR CLICODE = 1

/*
	If a given input is invalid, when a parameter is parsed and there was a positional argument, but a different one was given
*/
const CLI_INVALIDINPUT_ERROR CLICODE = 2

/*
	A trailing argument was required to complete the parsing of the tree, but the required argument could not be detected
*/
const CLI_REQARGUMENTNOTFOUND_ERROR CLICODE = 4

/*
	The provided data given to the commandline was incorrect
	Either the callback does not work, so you need to check your custom datatype callback for a given argument
	Or you have to recheck your defined json file and read it again
*/
const CLI_WRONGDATA_ERROR CLICODE = 8

/*
	The command line exists because of a help wildcard
*/
const CLI_HELPEXIT_ERROR CLICODE = 16

/*
	In the recursive parsing, when a argument is given by the userinput and the argument is not found in that
	tree level, throws a argument not found error
*/
const CLI_ARGUMENTNOTFOUND_ERROR CLICODE = 32

/*
	If the builder tries to overwrite a default or protected argument, this error is thrown
	Upon booting, the commandlinetoolkit will eject
	If interactive, the commandlinetool will just not read the config and throw an error but resume operation
*/
const CLI_ARGUMENTNOTOVERWRITEABLE_ERROR CLICODE = 64
const CLI_ARGUMENTPROTECTED_ERROR CLICODE = 64

/*
	This error is thrown when the commmandline is in interactive mode and experiences a fatal error of known or unkown origin
	Can be returned by a callback and will halt the stop execution of the entire program, similar to System.exit()
*/
const CLI_RUNTIME_ERROR = 128

/************************************************************************************************************************
Interacrive Shell Mode

providing a --interactive or -i when starting the commandline will lead to an executeable binary to keep running in a waitgroup shell struct
that handles constant and continuous input via the commandline input structure, (in future that can be overwritten by using a restcall api)
*/
const CLI_INTERACTIVE_MODE = 256
