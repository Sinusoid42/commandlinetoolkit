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

var originalSttyState bytes.Buffer

//the shell struct
type shell struct {
	_previnputs [][]byte
	
	//stores the prev input
	_input     []byte
	_lastInput []byte
	
	_action    int
	_currIndex int
	_rtFlag    int
	_exit      int
	
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
		_logging:    true,
		_exit:       0,
		_currIndex:  0,
		_previnputs: [][]byte{},
		
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
				fmt.Print(">")
				
				s._osHandler._sysCall = syscall.SIGINT //sysExit
				//reset buffers
				s._input = []byte{}
				s._lastInput = []byte{}
				
				for {
					
					if s._exit == 1 {
						if cmdline._verbose&CLI_VERBOSE_OS_SIG > 0 {
							fmt.Println("\n-->osHandler: Exiting out of os handling subroutine")
						}
						
						//run at the end once
						defer s._osHandler._wg.Add(-1)
						
						return
					}
					if s._osHandler._sysCall == 0 {
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
	
	input := []byte{}
	s._lastInput = []byte{}
	s._action = -1
	s._rtFlag = 0
	//lastByte := byte(0)
	
	s._osHandler._wg = sync.WaitGroup{}
	s._osHandler._wg.Add(1)
	
	removeTerminalBuffering()
	
	sh := func() {
		var bt = make([]byte, 1)
		
		fmt.Print(">")
		
		if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
			fmt.Println("\n-->shell: Booted shell subroutine")
		}
		
		for {
			
			os.Stdin.Read(bt)
			
			b := bt[0]
			
			//if b != lastByte || b == byte('\n') {
			
			//lastByte = b
			
			if b == byte('\n') {
				
				if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
					fmt.Println("\n-->shell: Registered CR")
				}
				
				s._previnputs = append(s._previnputs, input)
				s._currIndex = 0
				
				s._rtFlag = 1
				s._lastInput = input
				//s._input = s._lastInput
				input = []byte{}
				if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
					fmt.Print("\n-->shell: Previous input: ")
					fmt.Print(s._lastInput)
					fmt.Print("\n")
				}
			} else {
				
				fmt.Print(string(b))
				s._rtFlag = 0
				input = append(input, b)
				s._lastInput = input
			}
			//}
			
			//handle a delete in the same line
			if b == 127 && len(input) > 1 {
				//remove last char
				newInput := input[0 : len(input)-2]
				
				//replace char sequence in the current terminal line with empty string
				fmt.Print("\r")
				for i := 0; i < len(input); i++ {
					fmt.Print(" ")
				}
				//fill it back up from the beginning with full chars up to n-1
				fmt.Print("\r")
				fmt.Print(">" + string(newInput))
				
				input = newInput
				s._lastInput = newInput
				
				newInput = nil
				
			}
			
			//check for arrow input
			//handle arrow UP
			arrowInput := s._lastInput
			l := len(arrowInput)
			if l > 3 && arrowInput[l-3] == 27 && arrowInput[l-2] == 91 && arrowInput[l-1] == 65 {
				
				fmt.Print("\r") //keep cursor in the current line
				
				//clear the current line
				fmt.Print("\r")
				for i := 0; i < len(input); i++ {
					fmt.Print(" ")
				}
				
				fmt.Print("\r") //start the line at the beginning
				if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
					fmt.Println("\n-->shell: Arrow up")
				}
				
				fmt.Print(">ARROW UP")
				/*
					linputs := len(s._previnputs)
					if s._currIndex >= 0 && linputs > 0 && linputs > s._currIndex {
				
						s._lastInput = s._previnputs[linputs-1-s._currIndex]
						input = s._lastInput
						s._rtFlag = 0
						s._currIndex++
				
						fmt.Print("\n>" + string(s._lastInput))
				
					} else {
						input = []byte{}
						s._lastInput = []byte{}
				
						fmt.Print(">")
					}
				*/
				//go up
				//s._action = 1
			}
			
			//handle arrow down
			/*
				if l > 3 && arrowInput[l-3] == 27 && arrowInput[l-2] == 91 && arrowInput[l-1] == 66 {
					fmt.Print("\033[F>")
					//fmt.Printf("\033[F")
					if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
						fmt.Println("\n-->shell: Arrow down")
					}
			
					//fmt.Print("ARROW DOWN")
			
					linputs := len(s._previnputs)
					if s._currIndex > 0 && linputs > 0 && linputs > s._currIndex {
						s._lastInput = s._previnputs[linputs-1-s._currIndex]
						input = s._lastInput
						s._rtFlag = 0
						s._currIndex--
			
						fmt.Print(">" + string(s._lastInput))
			
					} else {
						input = []byte{}
						s._lastInput = []byte{}
					}
					//go down
					//s._action = 2
				}*/
			
			//everything reading finished, request newline adn returnflag is 1
			//if rtflag is one, we can also get the previous line input and parse it in the commandline parser
			
			if s._rtFlag == 1 {
				
				s._input = s._lastInput
				
				if cmdline._verbose&CLI_VERBOSE_SHELL_PARSE > 0 {
					fmt.Print("\n-->shell: Previous parseable input: ")
					fmt.Print(s._input)
					fmt.Print("\n")
				}
				
				if s._osHandler._sysCall == syscall.SIGINT {
					
					if len(s._input) > 0 {
						
						//maybe the user entered y|Y as first char
						if s._input[0] == byte('y') ||
							s._input[0] == byte('Y') ||
							//last chat can be \n, so we check last - 1
							(len(s._input) > 3 &&
								s._input[len(s._input)-2] == byte('y') ||
								s._input[len(s._input)-2] == byte('Y')) {
							
							fmt.Println("Exit 0")
							
							s._exit = 1
							
							setSttyState(&(originalSttyState))
							
							break
						} else {
							s._osHandler._sysCall = 0
							s._input = []byte{}
							s._lastInput = []byte{}
							fmt.Println("aborting...")
						}
					}
				}
				if s._logging {
					s.log(string(s._input))
				}
				
				// we have no signals that come from the syste,
				//we can run our own commands from this current commandline OR from a new binary that we could execute
				
				//here commandline.parse(input
				
			}
			//fmt.Print(thetabprefix)
			if s._rtFlag == 1 {
				s._rtFlag = 0
				fmt.Print("\n>")
			}
		}
		
		//code here is run, but sometimes the printing to the console takes longer
		setSttyState(&(originalSttyState))
		
		//run at exiting the scope
		defer s._osHandler._wg.Add(-1)
		
		//os.Exit(0)
	}
	go sh()
	
}
