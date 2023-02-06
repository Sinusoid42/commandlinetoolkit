package commandlinetoolkit

import (
	`bytes`
	"container/list"
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
	_inputs *list.List
	
	//stores the prev input
	_input     []byte
	_lastInput []byte
	
	_action    int
	_currIndex int
	_rtFlag    int
	
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
func newShell(programName string, _logging bool) *shell {
	s := &shell{
		//use the default unix/linux keyboardInterrupt
		
		//logging
		_logging: true,
		
		_currIndex: 0,
		_inputs:    list.New().Init(),
		
		_osHandler: osHandler{
			
			_sysCallInterrupt: syscall.SIGINT,
			//most recent syscall input
			_sysCall: 0,
		},
	}
	
	s.registerSystemSignalCallbacks()
	
	s.run(nil)
	
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

func (s *shell) registerSystemSignalCallbacks() {
	
	s._osHandler._sysSignal = make(chan os.Signal, 1)
	signal.Notify(s._osHandler._sysSignal, os.Interrupt)
	
	//the shell callback
	f := func() {
		for sig := range s._osHandler._sysSignal {
			// sig is a ^C, handle it
			if sig == nil {
			}
			
			if s._osHandler._sysCall != syscall.SIGINT {
				
				fmt.Println("Keyboard Interrupt")
				fmt.Println("Exit? y/n")
				fmt.Print(">")
				s._osHandler._sysCall = syscall.SIGINT //sysExit
				s._input = []byte{}
				s._lastInput = []byte{}
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
		for {
			
			os.Stdin.Read(bt)
			
			b := bt[0]
			
			//if b != lastByte || b == byte('\n') {
			
			//lastByte = b
			
			if b == byte('\n') {
				
				s._inputs.PushBack(input)
				s._currIndex = 0
				
				s._rtFlag = 1
				s._lastInput = input
				input = []byte{}
				fmt.Print("\n>")
			} else {
				
				fmt.Print(string(b))
				s._rtFlag = 0
				input = append(input, b)
				s._lastInput = input
			}
			//}
			
			//check for arrow input
			
			s._action = -1
			
			//delete
			if b == 127 && len(input) > 1 {
				
				newInput := input[0 : len(input)-2]
				//replace char sequence in the current terminal line
				
				fmt.Print("\r")
				for i := 0; i < len(input); i++ {
					fmt.Print(" ")
				}
				fmt.Print("\r")
				fmt.Print(">" + string(newInput))
				
				input = newInput
				s._lastInput = newInput
				
			}
			
			//handle arrow UP
			if len(input) >= 3 {
				l := len(input)
				if input[l-3] == 27 && input[l-2] == 91 && input[l-1] == 65 && len(input) > 1 {
					fmt.Print("\r")
					if s._inputs.Front() != nil {
						e := (s._inputs.Front().Next())
						if e != nil {
							input = (e.Value.([]byte))
							fmt.Print(">" + string(input))
						}
					} else {
						input = []byte{}
					}
					s._lastInput = input
					//go up
					s._action = 1
				} else if input[l-3] == 27 && input[l-2] == 91 && input[l-1] == 66 && len(input) > 1 {
					fmt.Print("\r")
					if s._inputs.Front() != nil {
						e := (s._inputs.Front()).Prev()
						if e != nil {
							input = (e.Value.([]byte))
							fmt.Print(">" + string(input))
						}
					} else {
						input = []byte{}
					}
					s._lastInput = input
					//go down
					s._action = 2
				}
			}
			
			if bytes.Compare(s._lastInput, s._input) == 0 {
				continue
			} else {
				if s._action == -1 && s._rtFlag == 1 {
					
					s._input = s._lastInput
					
					if s._osHandler._sysCall == syscall.SIGINT {
						
						if len(s._input) > 0 {
							
							if s._input[0] == byte('y') || s._input[0] == byte('Y') {
								s._osHandler._wg.Done()
								defer setSttyState(&(originalSttyState))
								
							} else {
								s._osHandler._sysCall = 0
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
			}
			//fmt.Print(thetabprefix)
		}
		defer setSttyState(&(originalSttyState))
	}
	go sh()
	
}
