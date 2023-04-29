package commandlinetoolkit

import (
	"fmt"
	"strings"
)

//here we define a struct that is actually our live shell that we have available when running the commandline in --interactive mode
//the --interactive or -i mode allows running the commandline with a given for loop
//we require a waitgroup here to be able to wait in the main routine, if this shell shall be run interactively
//if it is run interactively, this struct maintains only history and allows usage for arrow keys when a given input was given to the commandline
//in interactive mode the shell and cli requires a hot reloading of commands, so can rebuild with --rebuild

//not finished yet
//for the entire program

// the shell struct
type shell struct {
	_running bool
	//global buffer
	cmdline *CommandLine

	/********************
	  interface for the commandline parser and top level state machine of the program
	*/
	StringParseable
	KeyParseable

	_shellHandler *shellHandler

	_enabledHistory bool

	_logging bool
}

// create a new shell
func newShell(_programName string, _logging bool, _usehistory bool, cmdline *CommandLine) *shell {
	s := &shell{
		//use the default unix/linux keyboardInterrupt
		//logging

		_enabledHistory: true,
		cmdline:         cmdline,
		_logging:        _logging,
		_running:        false,

		_shellHandler: newShellHandler(_programName, _logging, _usehistory, cmdline),
	}

	return s
}

func (s *shell) GetParseArgs() []string {
	str := strings.Split(string(s._shellHandler._currentInputBuffer), " ")
	return str
}

func (s *shell) GetParseKeys() []Key {
	return s._shellHandler._currentInputBuffer
}

func (s *shell) log(input string) {

}

func (s *shell) set(attrib PROGRAM_ARGUMENT, clicode CLICODE) {
	s._shellHandler.set(attrib, clicode)
}

func (s *shell) get() PROGRAM_ARGUMENT {
	return s._shellHandler.getAttributes()
}

func (s *shell) getCode(attrib PROGRAM_ARGUMENT) CLICODE {

	return s._shellHandler.getAttributeCode(attrib)
}

/*
Run the shell in a secondary go routine, catch system calls within another go routine

Transfer routine data from routines via the
*/
func (s *shell) run(cmdline *CommandLine) {

	s._running = true

	//reader := bufio.NewReader(os.Stdin)

	//arrowreader := bufio.NewScanner(os.Stdin)

	//input := []byte{}
	//arrowCallBackInput := []byte{}

	//lastByte := byte(0)

	//set sys calls and make sys buffering

	s._shellHandler.boot()

	sh := func() {
		//	s._shellHandler._osHandler.registerSystemSignalCallbacks(s._shellHandler)

		//reprint line header
		//s._shellHandler.printPrefix()

		for {

			//read most recent pressed key if any
			s._shellHandler.read()

			//everything reading finished, request newline is processed and returnflag is 1
			//if rtflag is one, we can also get the previous line input and parse it in the commandline parser
			//shell._currInput is now the storage of  the most recent full line commandline Input that was parsed, WITHOUT the prefix
			if s._shellHandler.processState() {

				//we store the actual current line before any linebreaks in s._currInput

				s._shellHandler.debugInputBuffer()

				if s._shellHandler.handleState() == CLI_EXIT {
					break
				}

				s._shellHandler.handleSuggestions()

				s._shellHandler.handleHistory()

				if s.stringInputBuffer() == "test" {
					//fmt.Print("\nHAHA")

				}

				if s.stringInputBuffer() == "--verbose" {

					fmt.Print(COLOR_PINK_IBG)

					fmt.Print("\n-->shell: >>> START SHELL VERBOSE MODE <<<")

					fmt.Print(COLOR_RESET)
					cmdline._verbose |= CLI_VERBOSE_SHELL_PARSE | CLI_VERBOSE_SHELL | CLI_VERBOSE_OS_SIG | CLI_VERBOSE_PREDICT | CLI_VERBOSE_SHELL_BUFFER
					s._shellHandler._debugHandler._verbose |= cmdline._verbose

				}

				if s.stringInputBuffer() == "--verbose=false" {

					fmt.Print(COLOR_PINK_IBG)

					fmt.Print("\n-->shell: >>> END SHELL VERBOSE MODE <<<")

					fmt.Print(COLOR_RESET)

					cmdline._verbose = 0
					s.setVerbose(0)

				}
				args := []string{"executionfile"}

				args = append(args, strings.Fields(string(s._shellHandler._currentInputBuffer))...)

				cmdline.Parse(args)

				/*if s._inputDisplayBufferLength > 0 && !s.consumed() {

					s._historyFileHandler.append(string(s._currentInputBuffer))

				}*/
				//we have no signals that come from the syste,
				//we can run our own commands from this current commandline OR from a new binary that we could execute

				//here commandline.parse(input

			}

			s._shellHandler.newLine()
		}

		s.Exit()

	}
	go sh()

}
func (s *shell) setVerbose(code CLICODE) {
	s._shellHandler._debugHandler._verbose = code
}

func (s *shell) stringInputBuffer() string {
	return string(s._shellHandler._currentInputBuffer)

}

func (s *shell) Exit() {
	s._shellHandler.exit()
}

/*
*
moves the cursor the the right
debug method
*/
func (s *shell) moveRight() {
	fmt.Print(string(ARROW_RIGHT))
}

func (s *shell) wait() {

	s._shellHandler._osHandler._wg.Wait()

}
