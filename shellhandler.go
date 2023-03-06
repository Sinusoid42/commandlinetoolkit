package commandlinetoolkit

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// handler for all printing operations from the current shell file, better software architecture
type shellHandler struct {

	/********************

	 */
	//global buffer
	cmdline *CommandLine

	_debugHandler *debugHandler

	//handler for all the syscalls
	_osHandler *osHandler

	_exit int

	_attribs ATTRIBUTE

	/********************
	  state machine
	*/

	_consumed bool

	_prevByteFromStdinBuffer []byte

	_prevKeyFromStdin Key

	/********************

	 */
	//the most recent buffer after a user pressed enter
	_currentInputBuffer []Key
	_inputDisplayBuffer []Key

	//the amount of words
	_parseDepth int32

	//buffered for prediction deletion
	_inputDisplayBufferLength int

	/********************

	 */

	//if an active history of the shell is enabled (eg. using a startup command or by using the config file provided)
	//_enabledHistoryFile bool
	//_enabledHistory     bool
	_history    *history
	_previnputs [][]Key

	/********************

	 */
	//Predictions
	_predictionDisplayed bool
	_newestPrediction    string

	_requestSuggestionsState bool
	_requestSuggestions      bool

	_searchPredictionsState bool
	_latestFullWord         string

	_currentPredictionAvailable   bool
	_prevPredictionFDisplayLength int

	/********************

	 */
	//prefix and line start
	_prefixColor  Color
	_preFix       string
	_preFixLength int

	/********************
	  state machine actions
	*/
	//actions
	_rtAction    int
	_arrowAction int

	/********************
	  outputs
	*/
	//io
	_alert     bool
	_playAlert bool
}

func newShellHandler(_programName string, _logging bool, _usehistory bool, cmdline *CommandLine) *shellHandler {

	s := &shellHandler{
		_inputDisplayBuffer:      []Key{},
		_currentInputBuffer:      []Key{},
		_prevByteFromStdinBuffer: make([]byte, 1),
		_preFix:                  ">>>",
		cmdline:                  cmdline,
		_exit:                    0,
		_requestSuggestionsState: false,
		_prefixColor:             GenColor(ITALIC_COLORFONT, INTENSITY_COLORTYPE, CYAN_COLOR),
		_searchPredictionsState:  true,

		_previnputs: [][]Key{},
		//_enabledHistoryFile: _usehistory,
		_playAlert:    true,
		_consumed:     true,
		_debugHandler: newDebugHandler(),
		_osHandler:    newOSHandler(),
	}

	//if s._enabledHistoryFile {

	//logging is debug information in the shell log
	//is being stored with prefix flag, so the debug information is not rendered back into the shell ubon loading

	s._history = newHistoryFileHandler(_programName)

	s._history._enabledHistory = _usehistory

	if len(s._history._keyLines) > 0 {
		s._previnputs = s._history._keyLines

	}

	//}

	s._attribs = s.getAttributes()

	return s
}

func (s *shellHandler) GetParseArgs() []string {
	return strings.Split(string(s._currentInputBuffer), " ")
}

func (s *shellHandler) GetParseKeys() []Key {
	return s._currentInputBuffer
}

func (s *shellHandler) set(attib ATTRIBUTE, clicode CLICODE) {
	if clicode == CLI_SUCCESS {

		//add attribute to the binary arg represenation
		s._attribs |= attib
	} else {

		//remove only the specified
		attribs := ^attib & 0xFFFFFFF

		//update
		s._attribs = s._attribs & attribs

	}

	if s._attribs&HISTORY > 0 {
		s._history._enabledHistory = true
	} else {
		s._history._enabledHistory = false
	}

	if s._attribs&HISTORYFILE > 0 {
		if s._attribs&HISTORY == 0 {
			s._debugHandler.printError("   Historyfile not possible:\n   Need to enable History first\n")
		} else {
			s._history._enabledHistoryFile = true
		}
	} else {
		s._history._enabledHistoryFile = false
	}

	if s._attribs&SUGGESTIONS > 0 {
		s._requestSuggestionsState = true
	} else {
		s._requestSuggestionsState = false
	}

	if s._attribs&PREDICTIONS > 0 {
		s._searchPredictionsState = true
	} else {
		s._searchPredictionsState = false
	}

}

func (s *shellHandler) getAttributes() ATTRIBUTE {

	attrib := ATTRIBUTE(0)

	if s._history._enabledHistoryFile {
		attrib |= HISTORYFILE
	}
	if s._history._enabledHistory {
		attrib |= HISTORY
	}
	if s._searchPredictionsState {
		attrib |= PREDICTIONS
	}
	if s._requestSuggestionsState {
		attrib |= SUGGESTIONS
	}

	return attrib

}

func (s *shellHandler) getAttributeCode(attrib ATTRIBUTE) CLICODE {
	if s._attribs&attrib > 0 {
		return CLI_TRUE
	}
	return CLI_FALSE
}

/*******************************************************************************************************************

State machine Logic and Input Handling




*/

func (s *shellHandler) boot() {
	//boot the debug mode, enable debug logging etc
	s._debugHandler.boot()
	s._osHandler.removeTerminalBuffering()
	s._osHandler._wg = sync.WaitGroup{}
	s._osHandler._wg.Add(1)
	s._osHandler.registerSystemSignalCallbacks(s)
}

func (s *shellHandler) exit() {
	//code here is run, but sometimes the printing to the console takes longer

	s._osHandler.exit()

	//run at exiting the scope
	//s._osHandler._wg.Add(-1)

	s._history.close()

	os.Exit(0)
}

/*
*
Registers a 'y' or 'Y' confirm for a given Inputquestion
Returns true or false
*/
func (s *shellHandler) yesNoConfirm() bool {

	if len(s._currentInputBuffer) < 1 || len(s._currentInputBuffer) > 1 {
		return false
	}

	if s._currentInputBuffer[0] == Key('y') ||
		s._currentInputBuffer[0] == Key('Y') ||
		//last chat can be \n, so we checkInputProgram last - 1
		(len(s._currentInputBuffer) > 3 &&
			(s._currentInputBuffer[len(s._currentInputBuffer)-2] == Key('y') ||
				s._currentInputBuffer[len(s._currentInputBuffer)-2] == Key('Y'))) {
		return true
	}
	return false

}

func (s *shellHandler) handleState() CLICODE {

	if s.handleSIGINTExit() {
		return CLI_EXIT
	}

	s.handleClear()

	return CLI_SUCCESS
}

func (s *shellHandler) handleClear() {
	if string(s._currentInputBuffer) == "clear" {
		s.clearTerminal()
		s._consumed = true
	}
}

func (s *shellHandler) handleHistory() {

	s._history.append(string(s._currentInputBuffer))
	s._history._currentHistoryBufferLength++

}

/*******************************************************************************************************************

PROCESS STATE




*/

func (s *shellHandler) processState() bool {

	s.handleLineBreakInput()

	//if prev input is an arrow up or down, remove the
	if s.checkForArrow() {

		//remove the arrow key input
		s._inputDisplayBuffer = s._inputDisplayBuffer[:len(s._inputDisplayBuffer)-2]
		s._inputDisplayBufferLength = len(s._inputDisplayBuffer)

		if !s._history._enabledHistory {
			s._arrowAction = 0
		}
	}

	s.handleDelete()

	//checkInputProgram for arrow input
	//handle arrow UP

	s.handleArrowUp()

	//handle arrow down

	s.handleArrowDown()

	//if history is enabled, we scan through the previous inputs of the commandline

	s._history.iterateHistory(s)

	s.checkForCurrentPrediction()

	s.requestSuggestionsOnTab()

	return s._rtAction == 1
}

/**
State Functions
*/

/*
*

	Transfer the parseable inputstring from the current commandline into a second application buffer

that can be read from outside or callbacks can be fired for
(need to register callback functions)
*/
func (s *shellHandler) handleLineBreakInput() {
	if s._rtAction != 1 {
		return
	}

	s.removePrediction()

	s._debugHandler.debugReturn()

	if len(s._inputDisplayBuffer) > 0 {

		//fmt.Println(s._previnputs)

		s._previnputs = append(s._previnputs, s._inputDisplayBuffer)
	}
	//s._lastInputLength = len(s._lastInput)

	s._history.reset()

	s._consumed = false

	s._currentInputBuffer = s._inputDisplayBuffer

	//s._currInput = s._lastInput
	s._inputDisplayBuffer = []Key{}

	s._debugHandler.printVerboseBuffer(s._inputDisplayBuffer)

}

/*
*
Handles the deletion from a given previous character if valid and byteInput was 127
*/
func (s *shellHandler) handleDelete() {
	l := len(s._inputDisplayBuffer)
	//handle a delete in the same line
	if s._prevKeyFromStdin == KEY_DELETE && l > 0 {

		s._inputDisplayBuffer = s._inputDisplayBuffer[:l-1]
		s._inputDisplayBufferLength = l - 1

		s.reprintCurrentLine()

		//remove last char
		//replace char sequence in the current terminal line with empty string
		//fill it back up from the beginning with full chars up to n-1
	}
}

/*
*
Process the input, when a arrow up action is present
*/
func (s *shellHandler) handleArrowDown() {
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

		s._debugHandler.printVerbose(CLI_VERBOSE_SHELL, "\n-->shell: Arrow down")
		s.printPrefix()

		//s.moveRight()

		if !s._history.up() {
			s.clearCurrentLine()
			s.printPrefix()
		}

	}
}

/*
*
Process the input, when a arrow up action is present
*/
func (s *shellHandler) handleArrowUp() {
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

		s._debugHandler.printVerbose(CLI_VERBOSE_SHELL, "\n-->shell: Arrow up")

		s.printPrefix()

		//s.moveRight()

		s._history.down()

	}
}

/*
*
Creates a UserInterfaceInteraction Request when pressing TAB
*/
func (s *shellHandler) requestSuggestionsOnTab() {

	//|| s._inputDisplayBufferLength > 1

	if s._attribs&SUGGESTIONS == 0 || !s._requestSuggestionsState || s._currentPredictionAvailable || s._prevKeyFromStdin != KEY_TAB {
		return
	}

	if !s._requestSuggestions {

		s._consumed = true
		fmt.Print("There are " + strconv.Itoa(s.cmdline.numberOfSuggestions(strings.Split(string(s._currentInputBuffer), " "), s._parseDepth)) + " available Options. \nDisplay? y/n?\n")
		s.printPrefix()
	}

	s._requestSuggestions = true

}

/*******************************************************************************************************************

IO




*/

func (s *shellHandler) read() {

	//read every new char, when it is entered into the console
	os.Stdin.Read(s._prevByteFromStdinBuffer)

	//from the byte buffer, get the first char alwayss
	s._prevKeyFromStdin = Key(s._prevByteFromStdinBuffer[0])

	/*
		Handle the input of a linebreak
		s._rtFlag (shell.returnFlag)
		Store Inputs, reset current new input
	*/
	s._arrowAction = 0

	if s._prevKeyFromStdin == KEY_RETURN {

		s._rtAction = 1

	} else {
		s.handleKeyInput()

		s._inputDisplayBufferLength = len(s._inputDisplayBuffer)
	}

	//print entire buffer
	s._debugHandler.printVerbose(CLI_VERBOSE_SHELL, "\n")
	s._debugHandler.printVerbose(CLI_VERBOSE_SHELL, s._prevByteFromStdinBuffer)

	//debug with input code
	s._debugHandler.debugBufferSingle(s._prevKeyFromStdin)

	s._debugHandler.debugBuffer(s._inputDisplayBuffer)
	if s._debugHandler._verbose&CLI_VERBOSE_SHELL_BUFFER > 0 || s._debugHandler._verbose&CLI_VERBOSE_DEBUG > 0 {
		fmt.Print("\n")
		s.printPrefix()
	}
}

/*
*
Handles the exit when CTRL+C for unix/linux Keyboard Interrupt
*/
func (s *shellHandler) handleSIGINTExit() bool {
	if s._osHandler._sysCall == syscall.SIGINT {

		//maybe the user entered y|Y as first char
		if s.yesNoConfirm() {

			s._consumed = true

			fmt.Println("\nExit 0")

			s._exit = 1

			s._osHandler.reset()

			return true
		} else {

			s._debugHandler.printVerbose(CLI_VERBOSE_SHELL, s._currentInputBuffer)

			s._exit = 0
			s._osHandler._sysCall = 0
			s._currentInputBuffer = []Key{}
			s._inputDisplayBuffer = []Key{}

			fmt.Print("\naborting...")
		}

	}
	return false
}

/*******************************************************************************************************************

KEY INPUT




*/

/*
*
Handle the Input of a given byte from the raw stdin console callback by the os
resets state specific buffer variables such as autocompletion
*/
func (s *shellHandler) handleKeyInput() {

	s.removePrediction()

	//buffer und input darstellung trennen, bzw einen weiteren buffer hinzufügen, der dann den input festhält, der den arrow key etc. beinhaltet und den ESC nicht blockiert, ist der ESC drin, dann soll der einfach nur die line clearen
	if s._prevKeyFromStdin == KEY_DELETE ||
		s._prevKeyFromStdin == KEY_TAB ||
		s.checkForArrowInput(s._prevKeyFromStdin) {
		return
	}

	l := len(s._inputDisplayBuffer)

	multiKeyDebug := l > 0
	if multiKeyDebug {
		multiESC := s._inputDisplayBuffer[l-1] == KEY_ESC && s._prevKeyFromStdin == KEY_ESC
		multiSPACE := s._prevKeyFromStdin == KEY_SPACE && s._inputDisplayBuffer[l-1] == KEY_SPACE

		if multiSPACE || multiESC {
			return
		}
	}

	//print the keys
	fmt.Print(string(s._prevKeyFromStdin))

	s._rtAction = 0
	s._inputDisplayBuffer = append(s._inputDisplayBuffer, s._prevKeyFromStdin)

	s._predictionDisplayed = false
	s._prevPredictionFDisplayLength = 0

	//s._lastInputLength = len(s._lastInput)

}

// TODO or @deprecated
func (s *shellHandler) checkForArrow() bool {

	if s._arrowAction == 2 || s._arrowAction == 3 {
		return true
	}
	return false
}

func (s *shellHandler) newLine() {
	//fmt.Print(thetabprefix)
	if s._rtAction == 1 {
		s._rtAction = 0
		fmt.Print("\n")
		s.printPrefix()
	}
}

/*******************************************************************************************************************

DRAWING




*/

/*
*
Print the prefix for the custom Shell environment
*/
func (s *shellHandler) printPrefix() {

	fmt.Print("\r")
	fmt.Print(s._prefixColor)
	fmt.Print("\r")
	fmt.Print(s._preFix)
	fmt.Print(COLOR_RESET)
}

/*
*
clear the current line within the active shell
*/
func (s *shellHandler) clearCurrentLine() {

	inputlength := s._inputDisplayBufferLength

	s._debugHandler.printVerbose(CLI_VERBOSE_SHELL, "\n-->shell: currentInputLength_ClearCurrLine: ")
	s._debugHandler.printVerbose(CLI_VERBOSE_SHELL, inputlength+s._preFixLength+s._prevPredictionFDisplayLength)

	fmt.Print("\033[2K")

	s._prevPredictionFDisplayLength = 0

	fmt.Print("\033[0G")
	//fmt.Print("\u001b[{n}")

}

/*
*
Reprints the current line
*/
func (s *shellHandler) reprintCurrentLine() {
	s.clearCurrentLine()

	s._currentInputBuffer = []Key{}

	s.printPrefix()

	fmt.Print(string(s._inputDisplayBuffer))
	s._inputDisplayBufferLength = len(s._inputDisplayBuffer)
}

/*
*
Removes a current Prediction and reprints the current displayed line with the most recent entered buffer,
if arrows were used, the buffer is emnpty
*/
func (s *shellHandler) removePrediction() {
	if !s._searchPredictionsState || s._attribs&PREDICTIONS == 0 {
		return
	}

	s.clearKeys(s._prevPredictionFDisplayLength + s._inputDisplayBufferLength)

	s.reprintCurrentLine()
	s._predictionDisplayed = false
}

/*
*
Clears n-Keys
*/
func (s *shellHandler) clearKeys(n int) {

	for i := 0; i < n; i++ {
		fmt.Print("\b")
	}
}

/*
*
Clears the current Terminal Window
*/
func (s *shellHandler) clearTerminal() {

	fmt.Print("\033[2J \033[H")
	fmt.Print("\033[2J \033[H")

}

/**********************************************************************************************************************************************************

ARROW Keys




*/

func (s *shellHandler) checkForArrowInput(keyInput Key) bool {

	//remove the arrow input

	l := len(s._inputDisplayBuffer)
	if l < 2 {
		return false
	}

	if (keyInput == ARROW_UP[2] ||
		keyInput == ARROW_DOWN[2] ||
		keyInput == ARROW_LEFT[2] ||
		keyInput == ARROW_RIGHT[2]) && (s._inputDisplayBuffer[l-1] == ARROW_UP[1] && s._inputDisplayBuffer[l-2] == ARROW_UP[0]) {

		//action 0 is button left
		//action 1 is button right
		//action 2 is button down
		//action 3 is button up

		s._arrowAction = 68 - int(keyInput)

		return true
	}

	return false
}

/*
*
Removes a triplet of bytes from the input buffer when possible
*/
func (s *shellHandler) removeArrowKeyStrokeFromDisplayBuffer() {
	l := len(s._inputDisplayBuffer)
	if l > 2 {
		s._inputDisplayBuffer = s._inputDisplayBuffer[0 : l-3]

	}
}

/**********************************************************************************************************************************************************

PREDICTIONS




*/

/*
*
Returns the last valid subword before a given Spacebar from the current inputbuffer that was not entered yet
*/
func (s *shellHandler) latestFullInput() (string, int32) {
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

/*
*
Checks on every Key enter a possible current prediction based on the current output
If more match with a given substring displays only the first hit
WIll also be the same on TAB complete
*/
func (s *shellHandler) checkForCurrentPrediction() {
	if s._attribs&PREDICTIONS == 0 {
		return
	}
	if s._searchPredictionsState {

		code := CLICODE(-1)

		latestWord, layer := s.latestFullInput()
		s._newestPrediction, code = s.cmdline.checkPredictions(strings.Split(string(s._currentInputBuffer), " "), latestWord, layer)

		s._currentPredictionAvailable = true

		if code&CLI_SUCCESS > 0 {
			s.displayPrediction()

		} else {
			s._currentPredictionAvailable = false
		}
	}

	if s._searchPredictionsState && s._prevKeyFromStdin == KEY_TAB && s._currentPredictionAvailable {

		s.handleTabCompletion()
	}
}

/*
*
Bundle Function to reprint the current line from a given inpput and add the previously recieved prediction from the parse tree
*/
func (s *shellHandler) handleTabCompletion() {
	s.reprintCurrentLine()
	s.addPreviousPrediction()
}

/*
*
Displays a possible completion prediction
*/
func (s *shellHandler) displayPrediction() {

	l := len(s._newestPrediction)
	q := len(s._latestFullWord)
	k := l - q - 1

	if k <= 0 || s._attribs&PREDICTIONS == 0 {
		s._prevPredictionFDisplayLength = 0
		return
	}

	s._prevPredictionFDisplayLength = k

	s._predictionDisplayed = true

	fmt.Print(COLOR_GRAY_D)

	fmt.Print(s._newestPrediction[q:])

	fmt.Print(COLOR_RESET)
}

/*
*
Adds the previously displayewd Tab completion prediction to the current input string buffer
*/
func (s *shellHandler) addPreviousPrediction() {

	l := len(s._newestPrediction)
	q := len(s._latestFullWord)
	k := l - q - 1

	if k <= 0 || s._attribs&PREDICTIONS == 0 {
		s._prevPredictionFDisplayLength = 0
		return
	}

	s._inputDisplayBuffer = append(s._inputDisplayBuffer, []Key(s._newestPrediction[q:])...)

	s._predictionDisplayed = false
	s._prevPredictionFDisplayLength = 0
	s._newestPrediction = ""

	s.reprintCurrentLine()
}

/*
*
Helps displaying all possible suggestions for a given VALID parse tree layer
*/
func (s *shellHandler) handleSuggestions() {

	if s._requestSuggestions {

		s._debugHandler.debugSuggestions(s)

		s._requestSuggestions = false
		if s.yesNoConfirm() {

			s._consumed = true

			fmt.Print("\nPrinting" + strconv.Itoa(s.cmdline.numberOfSuggestions(strings.Split(string(s._currentInputBuffer), " "), s._parseDepth)) + "Options")

		} else {
			fmt.Print("\naborting...")
		}

	}
}

func (s *shellHandler) debugInputBuffer() {
	s._debugHandler.printVerboseBuffer(s._inputDisplayBuffer)
}
