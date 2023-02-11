package commandlinetoolkit

import (
	`fmt`
)

//base file for building a commandline
//here struct definition and runtime for a commandline is given
//for a command line to work we need the possibility of parsing a non cyclic directed, n-dimensional tree from arguments that follow conditionally
//either by setting the paramters required or not,

//const
var VERSION = "0.0.5"

type CommandLine struct {
	
	//here we have the commandline current parsetree
	_rootArgument *Argument
	
	//either provided with the json
	
	_program *program
	
	/*
		The parsers job is to traverse the parse tree (Parsing Explained - Computerphile 13.11.2019)
		Here we give the parser the given input from the commandline as array of strings,
		and the parsers job is to make sure either the input is error checked, well formattet,
		in case of callbacks that those are called and that options/parameters are jumped correctly
	*/
	_parser *parser
	
	/*
		the builder creates the argument tree from the provided json or programmed input
		either we can use a template from the commandlinetemplate.go lib implementation
			//as per default the root command --help execution cannot be overwritten, aswell as the --interactive or -i flags are not overwriteable
			//for running the -i --interactive there is also the --_logging OPTION, if _logging is provided either we can just provide it, resulting in TRUE, otherwise given --_logging True|TRUE|true, False|FALSE|false
			//they are reserved for either running the commandline within an own shell or printing the dynamically build short or full help menu
	
	*/
	_builder *builder
	
	_runtimeCodes CLICODE
	/*
		creates a .clihistory text file
	*/
	_logging    bool
	_usehistory bool
	
	//available arguments in total
	_size int32
	
	//available methods
	_methods int32
	
	//amount of options
	_options int32
	
	//run in verbose mode
	_verbose      int32
	_verboseColor Color
	
	//the interactive shell
	_shell *shell
}

func NewCommandLine() *CommandLine {
	
	cli := &CommandLine{
		_parser:       newparser(),
		_program:      newprogram(),
		_usehistory:   true,
		_verbose:      0,
		_verboseColor: COLOR_RED_I,
	}
	
	cli.Rebuild()
	return cli
}

func (c *CommandLine) ReadJSON(path string) {

}

func (c *CommandLine) Rebuild() CLICODE {
	c.Clear()
	
	c._builder.Rebuild(c)
	
	return CLI_SUCCESS
}

func (c CommandLine) JSON() string {
	return "TODO"
}

func (c *CommandLine) Clear() CLICODE {
	
	return CLI_SUCCESS
}

func (c CommandLine) String() string {
	return "commandlinetoolkit"
}

func (c *CommandLine) Parse(args []string) CLICODE {
	
	//program name is at args[0] always by definition
	//read out the name of the application in the first parameter
	//run the interactive mode of this commandlineparser
	
	//go c.runInteractive()
	
	//fmt.Println("DOING THE PARSING AND BLABLABLA")
	
	//either we have the programname from the executeable by running or from the json file provided
	
	c._shell = newShell(c._program._programName, c._logging, c._usehistory, c)
	
	//fmt.Println("DONE:..")
	
	c._shell.run(c)
	return CLI_SUCCESS
	
}

func (c *CommandLine) log(input string) {
	
	//to log previous commands into the commandline
}

func (c *CommandLine) runInteractive() {

}

func (c *CommandLine) Wait() {
	c._shell._osHandler._wg.Wait()
	
	//reset the original bash
	setSttyState(c._shell._originalSttyState)
	//dont need it anymore
	//os.Exit(0) //here we simulate the CTRL+C in case the syscall didnt get registered
}

func (c *CommandLine) checkPredictions(searchPrefix string, layer int32) (string, CLICODE) {
	
	if c._verbose&CLI_VERBOSE_PREDICT > 0 {
		
		c.printVerbose("\n-->CLI: Layer: ")
		c.printVerbose(layer)
		
	}
	
	// check available commands in the corresponding layer
	if len(searchPrefix) == 1 && searchPrefix[0] == 'a' {
		return "tabcompletion", CLI_SUCCESS
	}
	
	if len(searchPrefix) == 1 && searchPrefix[0] == 'c' {
		return "clear", CLI_SUCCESS
	}
	
	return "", CLI_NO_PREDICTION_ERROR
}

func (c *CommandLine) numberOfSuggestions(layer int32) int {
	
	return 69
	
}

/**
Prints with the verbose color overlay
*/
func (c *CommandLine) printVerbose(str interface{}) {
	fmt.Print(c._verboseColor)
	fmt.Print(str)
	fmt.Print(COLOR_RESET)
}
