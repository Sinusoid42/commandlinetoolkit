package commandlinetoolkit

import (
	`bufio`
	`fmt`
	`os`
	`os/signal`
)

//base file for building a commandline
//here struct definition and runtime for a commandline is given
//for a command line to work we need the possibility of parsing a non cyclic directed, n-dimensional tree from arguments that follow conditionally
//either by setting the paramters required or not,

type CommandLine struct {
	_rootArgument *Argument
	
	//either provided with the json
	_programName string
	
	_program string //a string defining the json file from this command line, can be stored using framework commands using --storeCLI <filename>
	
	/*
		The parsers job is to traverse the parse tree (Parsing Explained - Computerphile 13.11.2019)
		Here we give the parser the given input from the commandline as array of strings,
		and the parsers job is to make sure either the input is error checked, well formattet,
		in case of callbacks that those are called and that options/parameters are jumped correctly
	*/
	_parser parser
	
	/*
		the builder creates the argument tree from the provided json or programmed input
		either we can use a template from the commandlinetemplate.go lib implementation
			//as per default the root command --help execution cannot be overwritten, aswell as the --interactive or -i flags are not overwriteable
			//for running the -i --interactive there is also the --_logging OPTION, if _logging is provided either we can just provide it, resulting in TRUE, otherwise given --_logging True|TRUE|true, False|FALSE|false
			//they are reserved for either running the commandline within an own shell or printing the dynamically build short or full help menu
	
	*/
	_builder builder
	
	_runtimeCodes CLICODE
	/*
		creates a .clihistory text file
	*/
	_logging bool
	
	//available arguments in total
	_size int32
	
	//available methods
	_methods int32
	
	//amount of options
	_options int32
	
	//run in verbose mode
	_verbose bool
	
	//the interactive shell
	_shell shell
}

func NewCommandLine() *CommandLine {
	
	cli := &CommandLine{
		_parser:  newparser(),
		_program: DefaultCommandLineTemplate(),
	}
	cli.Rebuild()
	return cli
}

func (c *CommandLine) ReadJSON(path string) {

}

func (c *CommandLine) Rebuild() CLICODE {
	
	c._builder.build(c._program, c)
	
	return CLI_SUCCESS
}

func (c CommandLine) JSON() string {
	return "TODO"
}

func (c CommandLine) String() string {
	return "commandlinetoolkit"
}

func (c *CommandLine) Parse(args []string) CLICODE {
	
	//program name is at args[0] always by definition
	//read out the name of the application in the first parameter
	//run the interactive mode of this commandlineparser
	
	//go c.runInteractive()
	
	fmt.Println("DOING THE PARSING AND BLABLABLA")
	
	//either we have the programname from the executeable by running or from the json file provided
	
	c._shell = newShell(c._programName, c._logging)
	
	fmt.Println("DONE:..")
	
	return CLI_SUCCESS
	
}

func (c *CommandLine) log(input string) {

}

func (c *CommandLine) runInteractive() {
	
	//writer := bufio.NewWriter(os.Stdout)
	
	//mysyscall := os.Signal(syscall.SIGINT)
	
	scanner := bufio.NewReader(os.Stdin)
	
	sysSignal := make(chan os.Signal, 1)
	signal.Notify(sysSignal, os.Interrupt)
	
	sysCallCode := 0
	
	f := func() {
		for sig := range sysSignal {
			// sig is a ^C, handle it
			if sig == nil {
			}
			fmt.Println("Keyboard Interrupt")
			fmt.Println("Exit? y/n")
			fmt.Print(">")
			sysCallCode = 1 //sysExit
		}
		
	}
	
	go f()
	
	fmt.Println("TESTING")
	
	_input := ""
	
	for true {
		fmt.Print(">")
		input, _ := scanner.ReadString('\n')
		
		if _input == input {
			continue
		} else {
			
			if sysCallCode != 0 {
				
				if len(input) > 1 {
					if input[0] == 'y' {
						os.Exit(0)
					} else {
						sysCallCode = 0
					}
				}
			} else {
				if c._logging {
					c.log(input)
				}
				
			}
			
			// we have no signals that come from the syste,
			//we can run our own commands from this current commandline OR from a new binary that we could execute
			
			//here commandline.parse(input
			
			_input = input
		}
		
		//fmt.Print(thetabprefix)
	}
}
