package commandlinetoolkit

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//here we define a struct that is actually our live shell that we have available when running the commandline in --interactive mode
//the --interactive or -i mode allows running the commandline with a given for loop
//we require a waitgroup here to be able to wait in the main routine, if this shell shall be run interactively
//if it is run interactively, this struct maintains only history and allows usage for arrow keys when a given input was given to the commandline
//in interactive mode the shell and cli requires a hot reloading of commands, so can rebuild with --rebuild

type shell struct {
	_inputs *list.List
	
	//stores the prev input
	_input     string
	_currIndex int
	
	_logging bool
	
	_sysCall CLICODE
	
	_sysCallInterrupt CLICODE
	
	_wg sync.WaitGroup
}

func newShell(programName string, _logging bool) *shell {
	s := &shell{
		//use the default unix/linux keyboardInterrupt
		_sysCallInterrupt: syscall.SIGINT,
		//most recent syscall input
		_sysCall: 0,
		//logging
		_logging: true,
		
		_currIndex: 0,
		_inputs:    list.New(),
	}
	return s
}

func (s *shell) log(input string) {

}

func (s *shell) run(cmdline *CommandLine) {
	
	scanner := bufio.NewReader(os.Stdin)
	
	s._wg = sync.WaitGroup{}
	s._wg.Add(1)
	sysSignal := make(chan os.Signal, 1)
	signal.Notify(sysSignal, os.Interrupt)
	
	f := func() {
		for sig := range sysSignal {
			// sig is a ^C, handle it
			if sig == nil {
			}
			fmt.Println("Keyboard Interrupt")
			fmt.Println("Exit? y/n")
			fmt.Print(">")
			s._sysCall = syscall.SIGINT //sysExit
		}
		
	}
	
	go f()
	
	for true {
		fmt.Print(">")
		input, _ := scanner.ReadString('\n')
		
		if s._input == input {
			continue
		} else {
			
			if s._sysCall == syscall.SIGINT {
				
				if len(input) > 1 {
					if input[0] == 'y' {
						os.Exit(0)
					} else {
						s._sysCall = 0
					}
				}
			} else {
				if s._logging {
					s.log(input)
				}
				
			}
			
			// we have no signals that come from the syste,
			//we can run our own commands from this current commandline OR from a new binary that we could execute
			
			//here commandline.parse(input
			
			s._input = input
		}
		
		//fmt.Print(thetabprefix)
	}
	
}
