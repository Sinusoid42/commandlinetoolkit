package commandlinetoolkit

import (
	`bytes`
	"fmt"
	`log`
	"os"
	`os/exec`
	"os/signal"
	"sync"
	"syscall"
)

//here we define a struct that is actually our live shell that we have available when running the commandline in --interactive mode
//the --interactive or -i mode allows running the commandline with a given for loop
//we require a waitgroup here to be able to wait in the main routine, if this shell shall be run interactively
//if it is run interactively, this struct maintains only history and allows usage for arrow keys when a given input was given to the commandline
//in interactive mode the shell and cli requires a hot reloading of commands, so can rebuild with --rebuild

//not finished yet

var originalSttyState bytes.Buffer

//the shell struct
type shell struct {
	_previnputs [][]byte
	
	//stores the prev input
	_input           []byte
	_lastInput       []byte
	_lastInputLength int
	_preFix          string
	
	_action       int
	_currIndex    int
	_rtFlag       int
	_exit         int
	_preFixLength int
	
	_enabledHistory bool
	_alert          bool
	
	_playAlert bool
	
	_logging bool
	
	_osHandler osHandler
}

//operating system helper struct for sys signals and callbacks
type osHandler struct {
	_sysCall syscall.Signal
	
	_sysCallInterrupt syscall.Signal
	
	_wg sync.WaitGroup
	
	_sysSignal chan os.Signal
}

//create a new shell
func newShell(programName string, _logging bool, cmdline *CommandLine) *shell {
	s := &shell{
		//use the default unix/linux keyboardInterrupt
		
		//logging
		_logging:        true,
		_exit:           0,
		_currIndex:      0,
		_previnputs:     [][]byte{},
		_preFix:         ">>>",
		_preFixLength:   3,
		_enabledHistory: true,
		_playAlert:      true,
		_osHandler: osHandler{
			
			_sysCallInterrupt: syscall.SIGINT,
			//most recent syscall input
			_sysCall: 0,
		},
	}
	
	s.registerSystemSignalCallbacks(cmdline)
	
	return s
}

func getSttyState(state *bytes.Buffer) (err error) {
	//https://gist.github.com/mrnugget/9582788
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func setSttyState(state *bytes.Buffer) (err error) {
	//https://gist.github.com/mrnugget/9582788
	cmd := exec.Command("stty", state.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func removeTerminalBuffering() {
	//https://gist.github.com/mrnugget/9582788
	err := getSttyState(&(originalSttyState))
	if err != nil {
		log.Fatal(err)
	}
	//run in this function!
	//"/dev/tty", "raw", "-echo", "cbreak", "-g"
	//setSttyState(bytes.NewBufferString("-icanon"))
	//setSttyState(bytes.NewBufferString("min"))
	setSttyState(bytes.NewBufferString("-raw"))
	setSttyState(bytes.NewBufferString("cbreak"))
	setSttyState(bytes.NewBufferString("-echo"))
}

func (s *shell) log(input string) {

}

func (s *shell) registerSystemSignalCallbacks(cmdline *CommandLine) {
	
	s._osHandler._sysSignal = make(chan os.Signal, 1)
	signal.Notify(s._osHandler._sysSignal, os.Interrupt)
	
	//the shell callback
	f := func() {
		s._osHandler._wg.Add(1)
		//need double loop as syscall can be happening in different scope at certain times
		
		if cmdline._verbose&CLI_VERBOSE_OS_SIG > 0 {
			fmt.Println("osHandler: Booted os signal handling subroutine")
		}
		
		for sig := range s._osHandler._sysSignal {
			// sig is a ^C, handle it
			
			if sig == nil {
				continue
			}
			
			if sig == syscall.SIGINT {
				
				if cmdline._verbose&CLI_VERBOSE_OS_SIG > 0 {
					fmt.Println("\n-->osHandler: syscall.SIGINT")
				}
				
				fmt.Println("Keyboard Interrupt")
				fmt.Println("Exit? y/n")
				fmt.Print(s._preFix)
				
				s._osHandler._sysCall = syscall.SIGINT //sysExit
				//reset buffers
				//s._input = s._lastInput
				s._lastInput = []byte{}
				
				for {
					
					if s._exit == 1 {
						if cmdline._verbose&CLI_VERBOSE_OS_SIG > 0 {
							fmt.Println("\n-->osHandler: Exiting out of os handling subroutine")
						}
						
						//run at the end once
						s._osHandler._wg.Add(-1)
						
						return
					}
					
					if s._osHandler._sysCall == 0 || s._exit == 0 {
						break
					}
				}
			}
		}
		
	}
	
	go f()
}

/*
	Run the shell in a secondary go routine, catch system calls within another go routine

	Transfer routine data from routines via the
*/
func (s *shell) run(cmdline *CommandLine) {
	
	//reader := bufio.NewReader(os.Stdin)
	
	//arrowreader := bufio.NewScanner(os.Stdin)
	
	//input := []byte{}
	//arrowCallBackInput := []byte{}
	s._lastInput = []byte{}
	s._action = -1
	s._rtFlag = 0
	//lastByte := byte(0)
	
	s._osHandler._wg = sync.WaitGroup{}
	s._osHandler._wg.Add(1)
	
	removeTerminalBuffering()
	
	sh := func() {
		var bt = make([]byte, 1)
		
		fmt.Print(s._preFix)
		
		if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
			fmt.Print("\n-->shell: Booted shell subroutine")
		}
		
		for {
			
			//read every new char, when it is entered into the console
			os.Stdin.Read(bt)
			
			//from the byte buffer, get the first char alwayss
			byteInput := bt[0]
			
			/*
				Handle the input of a linebreak
				s._rtFlag (shell.returnFlag)
				Store Inputs, reset current new input
			*/
			if byteInput == byte('\n') {
				
				s.handleLineBreakInput(cmdline)
				
			} else {
				
				s.handleKeyInput(byteInput, cmdline)
				
			}
			
			//if prev input is an arrow up or down, remove the
			if s.checkForArrow() {
			}
			
			s.handleDelete(byteInput)
			
			//check for arrow input
			//handle arrow UP
			
			s._action = 0
			
			s.handleArrowUp(cmdline)
			
			//handle arrow down
			
			s.handleArrowDown(cmdline)
			
			//if history is enabled, we scan through the previous inputs of the commandline
			
			s.iterateHistory()
			
			//everything reading finished, request newline is processed and returnflag is 1
			//if rtflag is one, we can also get the previous line input and parse it in the commandline parser
			//shell._input is now the storage of  the most recent full line commandline Input that was parsed, WITHOUT the prefix
			if s._rtFlag == 1 {
				
				//s._input = s._lastInput
				
				if cmdline._verbose&CLI_VERBOSE_SHELL_PARSE > 0 {
					fmt.Print("\n-->shell: Previous parseable input: ")
					fmt.Print(s._input)
					fmt.Print("\n")
				}
				
				if s.handleSIGINTExit(cmdline) {
					break
				}
				
				if s._logging {
					s.log(string(s._input))
				}
				
				if string(s._input) == "test" {
					fmt.Print("\nHAHA")
				}
				
				if string(s._input) == "verbose" {
					fmt.Print("\n-->shell: Enabling verbose mode")
					cmdline._verbose |= CLI_VERBOSE_SHELL_PARSE
				}
				
				if string(s._input) == "!verbose" {
					fmt.Print("\n-->shell: Disabling verbose mode")
					cmdline._verbose = 0
				}
				
				//we have no signals that come from the syste,
				//we can run our own commands from this current commandline OR from a new binary that we could execute
				
				//here commandline.parse(input
				
			}
			//fmt.Print(thetabprefix)
			if s._rtFlag == 1 {
				s._rtFlag = 0
				fmt.Print("\n" + s._preFix)
			}
		}
		
		//code here is run, but sometimes the printing to the console takes longer
		setSttyState(&(originalSttyState))
		
		//run at exiting the scope
		s._osHandler._wg.Add(-1)
		
		os.Exit(0)
	}
	go sh()
	
}

func (s *shell) iterateHistory() {
	if s._enabledHistory && (s._action == 1 || s._action == 2) {
		linputs := len(s._previnputs)
		
		if s._currIndex >= 0 && linputs > s._currIndex {
			s._lastInput = s._previnputs[linputs-1-s._currIndex]
			s._rtFlag = 0
		} else {
			if -1 >= s._currIndex {
				s._lastInput = []byte{}
			}
			if s._currIndex < -1 {
				s._currIndex = -1
				s._alert = true
			}
			if s._currIndex >= linputs-1 {
				s._currIndex = linputs - 1
				s._alert = true
			}
		}
		if s._alert {
			s._alert = false
			if s._playAlert {
				fmt.Print("\a")
			}
		}
		fmt.Print(string(s._lastInput))
	}
}

func (s *shell) handleArrowDown(cmdline *CommandLine) {
	l := len(s._lastInput)
	if l > 2 && s._lastInput[l-3] == 27 && s._lastInput[l-2] == 91 && s._lastInput[l-1] == 66 {
		
		// "\033[F"
		
		//fmt.Print("\033[F") //keep the cursor in the line#
		//remove the arrow bytes from the buffer
		if l > 2 {
			s._lastInput = s._lastInput[0 : l-3]
		} else {
			s._lastInput = []byte{}
		}
		
		fmt.Print("\r")
		for i := 0; i < len(s._lastInput)+s._preFixLength+4; i++ {
			fmt.Print(" ") //clear the entire line
		}
		fmt.Print("\r") //start the line at the beginning again
		if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
			fmt.Println("\n-->shell: Arrow down")
		}
		fmt.Print(s._preFix)
		s._currIndex--
		s._action = 2
	}
}

func (s *shell) handleArrowUp(cmdline *CommandLine) {
	l := len(s._lastInput)
	if l > 2 && s._lastInput[l-3] == 27 && s._lastInput[l-2] == 91 && s._lastInput[l-1] == 65 {
		
		fmt.Print("\n") //keep the cursor in the line
		//remove the arrow bytes from the buffer
		if l > 2 {
			s._lastInput = s._lastInput[0 : l-3]
		} else {
			s._lastInput = []byte{}
		}
		
		//clear the current line
		fmt.Print("\r")
		for i := 0; i < len(s._lastInput)+s._preFixLength+4; i++ {
			fmt.Print(" ") //clear the entire line
		}
		fmt.Print("\r") //start the line at the beginning again
		//debug?
		if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
			fmt.Println("\n-->shell: Arrow up")
		}
		fmt.Print(s._preFix)
		s._currIndex++
		s._action = 1
	}
}

func (s *shell) handleSIGINTExit(cmdline *CommandLine) bool {
	if s._osHandler._sysCall == syscall.SIGINT {
		
		if len(s._input) > 0 && len(s._input) < 2 {
			
			//maybe the user entered y|Y as first char
			if s._input[0] == byte('y') ||
				s._input[0] == byte('Y') ||
				//last chat can be \n, so we check last - 1
				(len(s._input) > 3 &&
					(s._input[len(s._input)-2] == byte('y') ||
						s._input[len(s._input)-2] == byte('Y'))) {
				
				fmt.Println("\nExit 0")
				
				s._exit = 1
				
				setSttyState(&(originalSttyState))
				
				return true
			} else {
				if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
					fmt.Print("\nThe prev input: ")
					fmt.Print(s._input)
				}
				s._exit = 0
				s._osHandler._sysCall = 0
				s._input = []byte{}
				s._lastInput = []byte{}
				
				fmt.Print("\naborting...")
			}
		} else {
			if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
				fmt.Print("\nThe prev input: ")
				fmt.Print(s._input)
			}
			s._exit = 0
			s._osHandler._sysCall = 0
			s._input = []byte{}
			s._lastInput = []byte{}
			
			fmt.Print("\naborting...")
		}
	}
	return false
}

func (s *shell) handleLineBreakInput(cmdline *CommandLine) {
	if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
		fmt.Print("\n-->shell: Registered CR")
	}
	
	if len(s._lastInput) > 0 {
		s._previnputs = append(s._previnputs, s._lastInput)
	}
	//s._lastInputLength = len(s._lastInput)
	
	s._currIndex = -1
	
	s._rtFlag = 1
	
	s._input = s._lastInput
	
	//s._input = s._lastInput
	s._lastInput = []byte{}
	if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
		fmt.Print("\n-->shell: Previous input: ")
		fmt.Print(s._lastInput)
	}
}

func (s *shell) handleDelete(byteInput byte) {
	//handle a delete in the same line
	if byteInput == 127 && len(s._lastInput) > 0 {
		//remove last char
		s._lastInput = s._lastInput[0 : len(s._lastInput)-1]
		
		//replace char sequence in the current terminal line with empty string
		fmt.Print("\r")
		
		inputlength := len(s._lastInput)
		
		for i := 0; i < inputlength+1+s._preFixLength; i++ {
			fmt.Print(" ")
		}
		//fill it back up from the beginning with full chars up to n-1
		fmt.Print("\r")
		fmt.Print(s._preFix + string(s._lastInput))
		
	}
}

func (s *shell) handleKeyInput(byteInput byte, cmdline *CommandLine) {
	if byteInput == 127 || byteInput == byte('\n') {
		return
	}
	
	fmt.Print(string(byteInput))
	
	s._rtFlag = 0
	s._lastInput = append(s._lastInput, byteInput)
}

func (s *shell) checkForArrow() bool {
	
	return false
}
