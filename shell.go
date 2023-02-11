package commandlinetoolkit

import (
	`bytes`
	"fmt"
	`log`
	"os"
	`os/exec`
	"os/signal"
	`strconv`
	"sync"
	"syscall"
)

//here we define a struct that is actually our live shell that we have available when running the commandline in --interactive mode
//the --interactive or -i mode allows running the commandline with a given for loop
//we require a waitgroup here to be able to wait in the main routine, if this shell shall be run interactively
//if it is run interactively, this struct maintains only history and allows usage for arrow keys when a given input was given to the commandline
//in interactive mode the shell and cli requires a hot reloading of commands, so can rebuild with --rebuild

//not finished yet
//for the entire program

//the shell struct
type shell struct {
	
	//store the string file name for a local hidden .history file
	_historyFileHandler *historyHandler
	_enabledHistoryFile bool
	_previnputs         [][]Key
	
	_inputDisplayBuffer       []Key
	_inputDisplayBufferLength int
	_parseDepth               int32
	_consumed                 bool
	
	_currentInputBuffer []Key
	
	//stores the prev input
	
	_newestPrediction             string
	_predictionDisplayed          bool
	_currentPredictionAvailable   bool
	_prevPredictionFDisplayLength int
	_searchPredictions            bool
	_latestFullWord               string
	_requestSuggestions           bool
	
	_prefixColor  Color
	_verboseColor Color
	
	_preFix       string
	_arrowAction  int
	_currIndex    int
	_rtFlag       int
	_exit         int
	_preFixLength int
	
	_enabledHistory bool
	_alert          bool
	_playAlert      bool
	_showBytes      bool
	
	_logging bool
	_verbose int32
	
	_osHandler *osHandler
	
	_originalSttyState *bytes.Buffer
}

//operating system helper struct for sys signals and callbacks
type osHandler struct {
	_sysCall syscall.Signal
	
	_sysCallInterrupt syscall.Signal
	
	_wg sync.WaitGroup
	
	_sysSignal chan os.Signal
}

//create a new shell
func newShell(_programName string, _logging bool, _usehistory bool, cmdline *CommandLine) *shell {
	s := &shell{
		//use the default unix/linux keyboardInterrupt
		//logging
		_exit:               0,
		_currIndex:          0,
		_previnputs:         [][]Key{},
		_preFix:             ">>>",
		_preFixLength:       3,
		_enabledHistory:     true,
		_enabledHistoryFile: _usehistory,
		_logging:            _logging,
		_playAlert:          true,
		_consumed:           true,
		_originalSttyState:  &bytes.Buffer{},
		_prefixColor:        GenColor(ITALIC_COLORFONT, INTENSITY_COLORTYPE, CYAN_COLOR),
		_verboseColor:       GenColor(DARK_COLORFONT, INTENSITY_COLORTYPE, GRAY_COLOR),
		_searchPredictions:  true,
		_osHandler: &osHandler{
			
			_sysCallInterrupt: syscall.SIGINT,
			//most recent syscall input
			_sysCall: 0,
		},
	}
	
	if s._enabledHistoryFile {
		
		//logging is debug information in the shell log
		//is being stored with prefix flag, so the debug information is not rendered back into the shell ubon loading
		
		s._historyFileHandler = newHistoryFileHandler(_programName)
		
		keys_history := s._historyFileHandler._bufferedLines
		
		if len(keys_history) > 0 {
			s._previnputs = s._historyFileHandler._keyLines
			s._currIndex = -1
		}
		
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

func (s *shell) removeTerminalBuffering() {
	//https://gist.github.com/mrnugget/9582788
	err := getSttyState(s._originalSttyState)
	if err != nil {
		log.Fatal(err)
	}
	//run in this function!
	//"/dev/tty", "raw", "-echo", "cbreak", "-g"
	//setSttyState(bytes.NewBufferString("icanon raw cbreak -g"))
	//setSttyState(bytes.NewBufferString(" cbreak -echo"))
	//setSttyState(bytes.NewBufferString("min 1"))
	//setSttyState(bytes.NewBufferString("-raw"))
	//setSttyState(bytes.NewBufferString("cbreak"))
	//setSttyState(bytes.NewBufferString("-echo"))
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
		
		if s._verbose&CLI_VERBOSE_OS_SIG > 0 {
			fmt.Println("osHandler: Booted os signal handling subroutine")
		}
		
		for sig := range s._osHandler._sysSignal {
			// sig is a ^C, handle it
			
			if sig == nil {
				continue
			}
			
			if sig == syscall.SIGINT {
				
				if s._verbose&CLI_VERBOSE_OS_SIG > 0 {
					fmt.Println("\n-->osHandler: syscall.SIGINT")
				}
				
				fmt.Println("Keyboard Interrupt")
				fmt.Println("Exit? y/n")
				s.printPrefix()
				
				s._osHandler._sysCall = syscall.SIGINT //sysExit
				//reset buffers
				//s._currInput = s._lastInput
				s._inputDisplayBuffer = []Key{}
				
				for {
					
					if s._exit == 1 {
						if s._verbose&CLI_VERBOSE_OS_SIG > 0 {
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
	s._inputDisplayBuffer = []Key{}
	s._currentInputBuffer = []Key{}
	
	s._arrowAction = -1
	s._rtFlag = 0
	//lastByte := byte(0)
	
	s._osHandler._wg = sync.WaitGroup{}
	s._osHandler._wg.Add(1)
	
	s.removeTerminalBuffering()
	
	sh := func() {
		var bt = make([]byte, 1)
		
		if s._verbose&CLI_VERBOSE_SHELL > 0 {
			s.printVerbose("\n-->shell: Booted shell subroutine")
			
		}
		
		s.printPrefix()
		
		for {
			
			//read every new char, when it is entered into the console
			os.Stdin.Read(bt)
			
			//from the byte buffer, get the first char alwayss
			byteInput := Key(bt[0])
			
			/*
				Handle the input of a linebreak
				s._rtFlag (shell.returnFlag)
				Store Inputs, reset current new input
			*/
			s._arrowAction = 0
			if byteInput == KEY_RETURN {
				
				s._rtFlag = 1
				
			} else {
				s.handleKeyInput(byteInput, cmdline)
				
				s._inputDisplayBufferLength = len(s._inputDisplayBuffer)
			}
			
			if s._verbose&CLI_VERBOSE_SHELL_BUFFER > 0 {
				s.printVerbose("\n-->shell: numBytes: ")
				s.printVerbose(numBytesAvailable())
				s.printVerbose("\n")
				s.printVerbose("-->shell: inputbyte: ")
				s.printVerbose(bt)
			}
			
			s.handleLineBreakInput(cmdline)
			
			//if prev input is an arrow up or down, remove the
			if s.checkForArrow() {
			
			}
			
			s.handleDelete(byteInput)
			
			//check for arrow input
			//handle arrow UP
			
			s.handleArrowUp(cmdline)
			
			//handle arrow down
			
			s.handleArrowDown(cmdline)
			
			//if history is enabled, we scan through the previous inputs of the commandline
			
			s.iterateHistory()
			
			s.checkForCurrentPrediction(cmdline, byteInput)
			
			s.requestSuggestionsOnTab(cmdline, byteInput)
			
			//everything reading finished, request newline is processed and returnflag is 1
			//if rtflag is one, we can also get the previous line input and parse it in the commandline parser
			//shell._currInput is now the storage of  the most recent full line commandline Input that was parsed, WITHOUT the prefix
			if s._rtFlag == 1 {
				
				//we store the actual current line before any linebreaks in s._currInput
				
				if s._verbose&CLI_VERBOSE_SHELL_PARSE > 0 {
					s.printVerbose("\n-->shell: Previous parseable input: ")
					s.printVerbose(s._currentInputBuffer)
					s.printVerboseBuffer()
				}
				
				if s.handleSIGINTExit(cmdline) {
					break
				}
				
				s.handleSuggestions(cmdline)
				
				if s._logging {
					s.log(string(s._currentInputBuffer))
				}
				
				if string(s._currentInputBuffer) == "test" {
					fmt.Print("\nHAHA")
					
				}
				
				if string(s._currentInputBuffer) == "verbose" {
					
					fmt.Print(COLOR_PINK_IBG)
					
					fmt.Print("\n-->shell: >>> START SHELL VERBOSE MODE <<<")
					
					fmt.Print(COLOR_RESET)
					cmdline._verbose |= CLI_VERBOSE_SHELL_PARSE | CLI_VERBOSE_SHELL | CLI_VERBOSE_OS_SIG | CLI_VERBOSE_PREDICT | CLI_VERBOSE_SHELL_BUFFER
					s._verbose |= cmdline._verbose
					
				}
				
				if string(s._currentInputBuffer) == "!verbose" {
					
					fmt.Print(COLOR_PINK_IBG)
					
					fmt.Print("\n-->shell: >>> END SHELL VERBOSE MODE <<<")
					
					fmt.Print(COLOR_RESET)
					
					cmdline._verbose = 0
					s._verbose = cmdline._verbose
					
				}
				
				if string(s._currentInputBuffer) == "clear" {
					s.clearTerminal()
					
				}
				
				if s._inputDisplayBufferLength > 0 && !s._consumed {
					
					s._historyFileHandler.append(string(s._currentInputBuffer))
					
				}
				
				//we have no signals that come from the syste,
				//we can run our own commands from this current commandline OR from a new binary that we could execute
				
				//here commandline.parse(input
				
			}
			//fmt.Print(thetabprefix)
			if s._rtFlag == 1 {
				s._rtFlag = 0
				fmt.Print("\n")
				s.printPrefix()
			}
		}
		
		s.exit()
		
	}
	go sh()
	
}

func (s *shell) exit() {
	//code here is run, but sometimes the printing to the console takes longer
	setSttyState(s._originalSttyState)
	//reset raw
	setSttyState(bytes.NewBufferString("-raw"))
	setSttyState(bytes.NewBufferString("-icanon"))
	//run at exiting the scope
	s._osHandler._wg.Add(-1)
	
	if s._enabledHistoryFile {
		
		s._historyFileHandler.close()
		
	}
	
	os.Exit(0)
}

/**
Iterate the previous history in the present shell

TODO
Add serializing and dezerialising inputs from the previous history
*/
func (s *shell) iterateHistory() {
	if s._enabledHistory && (s._arrowAction == 2 || s._arrowAction == 3) {
		linputs := len(s._previnputs)
		
		if s._currIndex >= 0 && linputs > s._currIndex {
			s._inputDisplayBuffer = s._previnputs[linputs-1-s._currIndex]
			s._rtFlag = 0
		} else {
			if -1 >= s._currIndex {
				s._inputDisplayBuffer = []Key{}
				s.clearCurrentLine()
				s.printPrefix()
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
		s.reprintCurrentLine()
	}
}

/**
Process the input, when a arrow up action is present
*/
func (s *shell) handleArrowDown(cmdline *CommandLine) {
	//l := len(s._lastInput)
	//if l > 2 && s._lastInput[l-3] == 27 && s._lastInput[l-2] == 91 && s._lastInput[l-1] == 66 {
	
	if s._arrowAction == 2 {
		
		// "\033[F"
		
		//fmt.Print("\033[F") //keep the cursor in the line#
		//remove the arrow bytes from the buffer
		
		//s._lastInputLength = len(s._lastInput)
		s.removePrediction()
		
		//s.removeArrowKeyStrokeFromBuffer()
		
		s.clearCurrentLine()
		
		s.debug(cmdline._verbose, "\n-->shell: Arrow down")
		
		s.printPrefix()
		
		//s.moveRight()
		
		s._currIndex--
		
		if s._currIndex <= 0 {
			s.clearCurrentLine()
			s.printPrefix()
		}
		
	}
}

/**
Process the input, when a arrow up action is present
*/
func (s *shell) handleArrowUp(cmdline *CommandLine) {
	//l := len(s._lastInput)
	//if l > 2 && s._lastInput[l-3] == 27 && s._lastInput[l-2] == 91 && s._lastInput[l-1] == 65 {
	if s._arrowAction == 3 {
		//fmt.Print("\n") //keep the cursor in the line
		//remove the arrow bytes from the buffer
		
		//s._lastInputLength = len(s._lastInput)
		s.removePrediction()
		
		//s.removeArrowKeyStrokeFromBuffer()
		
		//clear the current line
		
		s.clearCurrentLine()
		
		//debug?
		s.debug(cmdline._verbose, "\n-->shell: Arrow up")
		
		s.printPrefix()
		
		//s.moveRight()
		
		s._currIndex++
	}
}

/**
Handles the exit when CTRL+C for unix/linux Keyboard Interrupt
*/
func (s *shell) handleSIGINTExit(cmdline *CommandLine) bool {
	if s._osHandler._sysCall == syscall.SIGINT {
		
		//maybe the user entered y|Y as first char
		if s.yesNoConfirm() {
			
			s._consumed = true
			
			fmt.Println("\nExit 0")
			
			s._exit = 1
			
			setSttyState(s._originalSttyState)
			
			return true
		} else {
			if s._verbose&CLI_VERBOSE_SHELL > 0 {
				fmt.Print("\nThe prev input: ")
				fmt.Print(s._currentInputBuffer)
			}
			s._exit = 0
			s._osHandler._sysCall = 0
			s._currentInputBuffer = []Key{}
			s._inputDisplayBuffer = []Key{}
			
			fmt.Print("\naborting...")
		}
		
	}
	return false
}

/**
	Transfer the parseable inputstring from the current commandline into a second application buffer
that can be read from outside or callbacks can be fired for
(need to register callback functions)
*/
func (s *shell) handleLineBreakInput(cmdline *CommandLine) {
	if s._rtFlag != 1 {
		return
	}
	
	s.removePrediction()
	
	if s._verbose&CLI_VERBOSE_SHELL > 0 {
		s.printVerbose("\n-->shell: Registered RETURN KEY")
	}
	
	if len(s._inputDisplayBuffer) > 0 {
		s._previnputs = append(s._previnputs, s._inputDisplayBuffer)
	}
	//s._lastInputLength = len(s._lastInput)
	
	s._currIndex = -1
	
	s._currentInputBuffer = s._inputDisplayBuffer
	
	//s._currInput = s._lastInput
	s._inputDisplayBuffer = []Key{}
	if cmdline._verbose&CLI_VERBOSE_SHELL > 0 {
		s.printVerbose("\n-->shell: Previous input: ")
		s.printVerbose(s._inputDisplayBuffer)
	}
}

/**
Handles the deletion from a given previous character if valid and byteInput was 127
*/
func (s *shell) handleDelete(byteInput Key) {
	l := len(s._inputDisplayBuffer)
	//handle a delete in the same line
	if byteInput == KEY_DELETE && l > 0 {
		
		s._inputDisplayBuffer = s._inputDisplayBuffer[:l-1]
		s._inputDisplayBufferLength = l - 1
		
		s.reprintCurrentLine()
		
		//remove last char
		//replace char sequence in the current terminal line with empty string
		
		//fill it back up from the beginning with full chars up to n-1
		
	}
}

/**
Bundle Function to reprint the current line from a given inpput and add the previously recieved prediction from the parse tree
*/
func (s *shell) handleTabCompletion() {
	s.reprintCurrentLine()
	s.addPreviousPrediction()
}

/**
clear the current line within the active shell
*/
func (s *shell) clearCurrentLine() {
	
	inputlength := s._inputDisplayBufferLength
	
	if s._verbose&CLI_VERBOSE_SHELL > 0 {
		s.printVerbose("\n-->shell: currentInputLength_ClearCurrLine: ")
		s.printVerbose(inputlength + s._preFixLength + s._prevPredictionFDisplayLength)
		
	}
	
	fmt.Print("\033[2K")
	
	s._prevPredictionFDisplayLength = 0
	
	fmt.Print("\033[0G")
	//fmt.Print("\u001b[{n}")
	
}

/**
Handle the Input of a given byte from the raw stdin console callback by the os
resets state specific buffer variables such as autocompletion
*/
func (s *shell) handleKeyInput(byteInput Key, cmdline *CommandLine) {
	
	s.removePrediction()
	
	//buffer und input darstellung trennen, bzw einen weiteren buffer hinzufügen, der dann den input festhält, der den arrow key etc. beinhaltet und den ESC nicht blockiert, ist der ESC drin, dann soll der einfach nur die line clearen
	if byteInput == KEY_DELETE ||
		byteInput == KEY_TAB ||
		s.checkForArrowInput(byteInput) {
		return
	}
	
	l := len(s._inputDisplayBuffer)
	
	multiKeyDebug := l > 0
	if multiKeyDebug {
		multiESC := s._inputDisplayBuffer[l-1] == KEY_ESC && byteInput == KEY_ESC
		multiSPACE := byteInput == KEY_SPACE && s._inputDisplayBuffer[l-1] == KEY_SPACE
		
		if multiSPACE || multiESC {
			return
		}
	}
	
	fmt.Print(string(byteInput))
	
	if s._showBytes {
		s.printVerbose([]Key{byteInput})
	}
	
	s._rtFlag = 0
	s._inputDisplayBuffer = append(s._inputDisplayBuffer, byteInput)
	
	s._predictionDisplayed = false
	s._prevPredictionFDisplayLength = 0
	
	//s._lastInputLength = len(s._lastInput)
	
}

//TODO or @deprecated
func (s *shell) checkForArrow() bool {
	
	return false
}

func (s *shell) checkForArrowInput(keyInput Key) bool {
	
	//remove the arrow input
	
	l := len(s._inputDisplayBuffer)
	if l < 2 {
		return false
	}
	
	if (keyInput == ARROW_UP[2] ||
		keyInput == ARROW_DOWN[2] ||
		keyInput == ARROW_LEFT[2] ||
		keyInput == ARROW_RIGHT[2]) && (
		s._inputDisplayBuffer[l-1] == ARROW_UP[1] && s._inputDisplayBuffer[l-2] == ARROW_UP[0]) {
		
		//action 0 is button left
		//action 1 is button right
		//action 2 is button down
		//action 3 is button up
		
		s._arrowAction = 68 - int(keyInput)
		
		return true
	}
	
	return false
}

/**
Print the prefix for the custom Shell environment
*/
func (s *shell) printPrefix() {
	
	fmt.Print("\r")
	fmt.Print(s._prefixColor)
	fmt.Print("\r")
	fmt.Print(s._preFix)
	fmt.Print(COLOR_RESET)
}

/**
moves the cursor the the right
debug method
*/
func (s *shell) moveRight() {
	fmt.Print(string(ARROW_RIGHT))
}

/**
General debug function to print debug messages in the printVerbose handle function
*/
func (s *shell) debug(verbose int32, msg string) {
	if verbose&CLI_VERBOSE_SHELL > 0 {
		s.printVerbose(msg)
	}
}

/**
Removes a triplet of bytes from the input buffer when possible
*/
func (s *shell) removeArrowKeyStrokeFromDisplayBuffer() {
	l := len(s._inputDisplayBuffer)
	if l > 2 {
		s._inputDisplayBuffer = s._inputDisplayBuffer[0 : l-3]
		
	}
}

/**
Reprints the current line
*/
func (s *shell) reprintCurrentLine() {
	s.clearCurrentLine()
	
	s.printPrefix()
	
	fmt.Print(string(s._inputDisplayBuffer))
	s._inputDisplayBufferLength = len(s._inputDisplayBuffer)
}

/**
Displays a possible completion prediction
*/
func (s *shell) displayPrediction() {
	
	l := len(s._newestPrediction)
	q := len(s._latestFullWord)
	k := l - q - 1
	
	if k <= 0 {
		s._prevPredictionFDisplayLength = 0
		return
	}
	
	s._prevPredictionFDisplayLength = k
	
	s._predictionDisplayed = true
	
	fmt.Print(COLOR_GRAY_D)
	
	fmt.Print(s._newestPrediction[q:])
	
	fmt.Print(COLOR_RESET)
}

/**
Adds the previously displayewd Tab completion prediction to the current input string buffer
*/
func (s *shell) addPreviousPrediction() {
	
	l := len(s._newestPrediction)
	q := len(s._latestFullWord)
	k := l - q - 1
	
	if k <= 0 {
		s._prevPredictionFDisplayLength = 0
		return
	}
	
	s._inputDisplayBuffer = append(s._inputDisplayBuffer, []Key(s._newestPrediction[q:])...)
	
	s._predictionDisplayed = false
	s._prevPredictionFDisplayLength = 0
	s._newestPrediction = ""
	
	s.reprintCurrentLine()
}

/**
Removes a current Prediction and reprints the current displayed line with the most recent entered buffer,
if arrows were used, the buffer is emnpty
*/
func (s *shell) removePrediction() {
	if !s._searchPredictions {
		return
	}
	
	s.clearKeys(s._prevPredictionFDisplayLength + s._inputDisplayBufferLength)
	
	s.reprintCurrentLine()
	s._predictionDisplayed = false
}

/**
Clears n-Keys
*/
func (s *shell) clearKeys(n int) {
	
	for i := 0; i < n; i++ {
		fmt.Print("\b")
	}
}

/**
Returns the last valid subword before a given Spacebar from the current inputbuffer that was not entered yet
*/
func (s *shell) latestFullInput() (string, int32) {
	l := len(s._inputDisplayBuffer)
	s._latestFullWord = ""
	if l <= 0 {
		s._parseDepth = -1
		return "", -1
	} else {
		str := ""
		a := true
		count := int32(0)
		for i := 0; i < l; i++ {
			c := s._inputDisplayBuffer[l-1-i]
			{
				if c == KEY_SPACE && i > 0 {
					a = false
					// if the last key in the current line is not a space, we at least have min a one char praseable argument if i > 0 {
					count++
					
				}
				if a {
					s._latestFullWord += string(c)
					str += string(c)
				}
			}
		}
		
		s._parseDepth = count
		
		return str, count
	}
	return "", -1
}

/**
Checks on every Key enter a possible current prediction based on the current output
If more match with a given substring displays only the first hit
WIll also be the same on TAB complete
*/
func (s *shell) checkForCurrentPrediction(cmdline *CommandLine, byteInput Key) {
	if s._searchPredictions {
		
		code := CLICODE(-1)
		
		s._newestPrediction, code = cmdline.checkPredictions(s.latestFullInput())
		
		s._currentPredictionAvailable = true
		
		if code&CLI_SUCCESS > 0 {
			s.displayPrediction()
			
		} else {
			s._currentPredictionAvailable = false
		}
	}
	
	if s._searchPredictions && byteInput == KEY_TAB && s._currentPredictionAvailable {
		
		s.handleTabCompletion()
	}
}

/**
Creates a UserInterfaceInteraction Request when pressing TAB
*/
func (s *shell) requestSuggestionsOnTab(cmdline *CommandLine, byteInput Key) {
	
	if s._currentPredictionAvailable || byteInput != KEY_TAB {
		return
	}
	
	if !s._requestSuggestions {
		
		s._consumed = true
		fmt.Print("There are " + strconv.Itoa(cmdline.numberOfSuggestions(s._parseDepth)) + " available Options. \nDisplay? y/n?\n")
		s.printPrefix()
	}
	
	s._requestSuggestions = true
	
}

/**
Helps displaying all possible suggestions for a given VALID parse tree layer
*/
func (s *shell) handleSuggestions(cmdline *CommandLine) {
	
	if s._verbose&CLI_VERBOSE_SHELL_PARSE > 0 {
		s.printVerbose("\n-->shell: Requesting current-layer suggestions: ")
		s.printVerbose("Layer")
		s.printVerbose(s._parseDepth)
		s.printVerbose("; Request: ")
		s.printVerbose(s._requestSuggestions)
		
	}
	if s._requestSuggestions {
		
		s._requestSuggestions = false
		if s.yesNoConfirm() {
			
			s._consumed = true
			
			fmt.Print("\nPrinting" + strconv.Itoa(cmdline.numberOfSuggestions(s._parseDepth)) + "Options")
			
		} else {
			fmt.Print("\naborting...")
		}
		
	}
}

/**
Registers a 'y' or 'Y' confirm for a given Inputquestion
Returns true or false
*/
func (s *shell) yesNoConfirm() bool {
	
	if len(s._currentInputBuffer) < 1 || len(s._currentInputBuffer) > 1 {
		return false
	}
	
	if s._currentInputBuffer[0] == Key('y') ||
		s._currentInputBuffer[0] == Key('Y') ||
		//last chat can be \n, so we check last - 1
		(len(s._currentInputBuffer) > 3 &&
			(s._currentInputBuffer[len(s._currentInputBuffer)-2] == Key('y') ||
				s._currentInputBuffer[len(s._currentInputBuffer)-2] == Key('Y'))) {
		return true
	}
	return false
	
}

/**
Clears the current Terminal Window
*/
func (s *shell) clearTerminal() {
	
	fmt.Print("\033[2J\033[H")
	
}

/**
Prints with the verbose color overlay
*/
func (s *shell) printVerbose(str interface{}) {
	fmt.Print(s._verboseColor)
	fmt.Print(str)
	fmt.Print(COLOR_RESET)
}

/**
Prints the current input buffer with spaces, in debug mode
*/
func (s *shell) printVerboseBuffer() {
	sr := "["
	for _, i := range s._currentInputBuffer {
		if i == KEY_ESC {
			sr += " ESC,"
		} else if i == KEY_DELETE {
			sr += " DEL,"
		} else {
			sr += "" + string(i) + ", "
		}
	}
	l := len(sr)
	if l > 1 {
		sr = sr[:l-2]
	}
	sr += "]"
	s.printVerbose(sr)
}
