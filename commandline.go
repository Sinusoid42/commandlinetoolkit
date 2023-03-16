package commandlinetoolkit

import (
	`fmt`
	`os`
)

//base file for building a commandline
//here struct definition and runtime for a commandline is given
//for a command line to work we need the possibility of parsing a non cyclic directed, n-dimensional tree from arguments that follow conditionally
//either by setting the paramters required or not,

//const
var VERSION = "0.1.2"

type CommandLine struct {
	
	//here we have the commandline current parsetree
	
	//either provided with the json
	
	_program *program
	
	/*
		The parsers job is to traverse the parse tree (Parsing Explained - Computerphile 13.11.2019)
		Here we give the parser the given input from the commandline as array of strings,
		and the parsers job is to make sure either the input is error checked, well formattet,
		in case of callbacks that those are called and that options/parameters are jumped correctly
	*/
	_parser *parser
	
	_debugHandler *debugHandler
	
	/*
		the builder creates the argument tree from the provided json or programmed input
		either we can use a template from the commandlinetemplate.go lib implementation
			//as per default the root command --help execution cannot be overwritten, aswell as the --interactive or -i flags are not overwriteable
			//for running the -i --interactive there is also the --_logging OPTION, if _logging is provided either we can just provide it, resulting in TRUE, otherwise given --_logging True|TRUE|true, False|FALSE|false
			//they are reserved for either running the commandline within an own shell or printing the dynamically build short or full help menu
	
	*/
	
	_runtimeCodes CLICODE
	
	/*
		creates a .clihistory text file
	*/
	
	_logging bool
	
	_booted bool
	
	//available arguments in total
	_size int32
	
	//available methods
	_methods int32
	
	//amount of options
	_options int32
	
	_enabled    bool
	_printTitle bool
	
	//run in verbose mode
	_verbose      CLICODE
	_verboseColor Color
	
	//the interactive shell
	_enabledShell bool
	
	_shell *shell
}

func NewCommandLine() *CommandLine {
	
	cli := &CommandLine{
		_parser:       newparser(),
		_program:      newprogram("config.json"),
		_debugHandler: newDebugHandler(),
		_verbose:      0,
		_enabledShell: false,
		_verboseColor: COLOR_RED_I,
	}
	
	cli._shell = newShell(cli._program._programName, false, false, cli)
	cli.Rebuild()
	
	return cli
}

func (c *CommandLine) ReadJSON(path string) {
	
	c._program.readJsonProgram(path)
	
	c._program.checkInputProgram()
	
	c._program.genTitle()
	
	if title, err := c._parser.parseProgram(c._program); err == nil {
		
		c._enabled = true
		if c._printTitle {
			fmt.Print(title)
		}
		
	} else {
		
		fmt.Print(err)
	}
	
	//print the parsetree after fully reading a configuration file or
	//fmt.Println(c._parser._parseTree)
	
}

func (c *CommandLine) Set(attrib PROGRAM_ARGUMENT, clicode CLICODE) {
	
	c._shell.set(attrib, clicode)
	
	if c._shell._shellHandler.getAttributeCode(SHELL) == CLI_TRUE {
		c._enabledShell = true
	} else {
		c._enabledShell = false
	}
	
}

func (c *CommandLine) GetCode(attrib PROGRAM_ARGUMENT) CLICODE {
	if c._enabledShell {
		return CLI_SUCCESS
	}
	return c._shell.getCode(attrib)
}

func (c *CommandLine) Get() PROGRAM_ARGUMENT {
	if c._enabledShell {
		return c._shell.get() | SHELL
	}
	return c._shell.get()
}

func (c *CommandLine) newShell() {
	if c._shell == nil {
		c._shell = newShell(c._program._programName, false, false, c)
	}
}

func (c *CommandLine) Rebuild() CLICODE {
	
	//c.Clear()
	
	return CLI_SUCCESS
}

func (c CommandLine) JSON() string {
	return "TODO"
}

func (c *CommandLine) Clear() CLICODE {
	c._shell._shellHandler.clearTerminal()
	return CLI_SUCCESS
}

func (c CommandLine) String() string {
	return "commandlinetoolkit"
}

func (c *CommandLine) Exit() {
	if c._enabledShell {
		c._shell._shellHandler._osHandler._wg.Done()
		c._shell._shellHandler._osHandler._wg.Done()
		
	}
}

func (c *CommandLine) Parse(args []string) CLICODE {
	//program name is at args[0] always by definition
	//read out the name of the application in the first parameter
	//run the interactive mode of this commandlineparser
	
	//go c.runInteractive()
	
	//fmt.Println("DOING THE PARSING AND BLABLABLA")
	
	//either we have the programname from the executeable by running or from the json file provided
	
	//fmt.Println("DONE:..")
	
	if !c._enabled {
		c._debugHandler.printError("-->Commandline: Not enabled, required is a cli program!\n")
		v := c._debugHandler._verbose
		c._debugHandler._verbose = CLI_VERBOSE_DEBUG
		c._verboseColor = COLOR_PINK_L
		c.printVerbose("\nTo set a program, define a *.json file and run \n  --config=<myfile> or \n  --config <myfile>\n")
		c._debugHandler._verbose = v
		
		if !c._enabledShell {
			
			v := c._debugHandler._verbose
			c._debugHandler._verbose = CLI_VERBOSE_DEBUG
			c._verboseColor = COLOR_GRAY_debug
			c.printVerbose("\nshell not activated\naborting...\n")
			c._debugHandler._verbose = v
			
			os.Exit(1)
			
		}
	}
	
	//fmt.Println(c._parser._parseTree)
	
	//we can only run this function as a user
	
	//i gereally expect this first to be either a binary or the mainfile or a caller
	binaryFileName := args[0]
	
	c._program._programFile = binaryFileName
	
	arguments := args[1:]
	
	execTree, ok := c._parser.parse(arguments)
	
	c._parser._executeableTree = execTree
	
	if ok&CLI_SUCCESS == 0 && len(arguments) > 0 {
		fmt.Println("\n")
		c._debugHandler.printError("Unsuccessful\n")
	}
	
	ok = execTree.execute(c)
	
	//fmt.Println("Run:", ok)
	
	/*
		if c._enabledShell {
			c._shell.run(c)
	
		}*/
	return CLI_SUCCESS
	
}

func (c *CommandLine) log(input string) {
	
	//to log previous commands into the commandline
}

func (c *CommandLine) runInteractive() {

}

func (c *CommandLine) Wait() {
	
	if c._shell._shellHandler._attribs&SHELL > 0 {
		
		//programmer has to implement a WAIT in the end
		
		if c._enabledShell {
			c._shell.run(c)
		}
		
		c._enabledShell = true
	}
	
	c._shell.wait()
	
	c._shell.Exit()
	//reset the original bash
	
	//dont need it anymore
	//os.Exit(0) //here we simulate the CTRL+C in case the syscall didnt get registered
}

func (c *CommandLine) checkPredictions(args []string, searchPrefix string, layer int32) (string, CLICODE) {
	
	if c._verbose&CLI_VERBOSE_PREDICT > 0 {
		
		c.printVerbose("\n-->CLI: Layer: ")
		c.printVerbose(layer)
		
	}
	
	// checkInputProgram available commands in the corresponding layer
	if len(searchPrefix) == 1 && searchPrefix[0] == 'a' {
		return "tabcompletion", CLI_SUCCESS
	}
	
	if len(searchPrefix) == 1 && searchPrefix[0] == 'c' {
		return "clear", CLI_SUCCESS
	}
	
	return "", CLI_NO_PREDICTION_ERROR
}

func (c *CommandLine) numberOfSuggestions(args []string, layer int32) int {
	
	return 69
	
}

func (c *CommandLine) getSuggestions(args []string, layer int32) []string {
	return nil
}

func (c *CommandLine) Program() *parsetree {
	
	p := c._parser._parseTree.clone()
	
	return p
}

func (c *CommandLine) ParseTree() *parsetree {
	
	p := c._parser._executeableTree.clone()
	
	return p
}

/**
Prints with the verbose color overlay
*/
func (c *CommandLine) printVerbose(str interface{}) {
	c._debugHandler._verboseColor = c._verboseColor
	c._debugHandler.printVerbose(CLI_VERBOSE_DEBUG, str)
}

func (c *CommandLine) SetVerbosity(code CLICODE) {

}

func (c *CommandLine) PrintTitle(print bool) {
	c._printTitle = print
}
